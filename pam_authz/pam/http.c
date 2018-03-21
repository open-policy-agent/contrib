#include <stdlib.h>
#include <string.h>

#include <curl/curl.h>
#include <jansson.h>

#include "http.h"
#include "json.h"
#include "log.h"
#include "flag.h"

char HTTP_METHOD_GET[]  = "GET";
char HTTP_METHOD_POST[] = "POST";

// join_url combines base and endpoint, and stores the result in dest.
// It returns dest.
static const char *join_url(char* dest, const char *base, const char *endpoint) {
	return strcat(strcpy(dest, base), endpoint);
}

struct curl_data {
	char *payload;
	size_t size;
};

static const int CURL_CALLBACK_ERR = -1;

// Custom callback to store response data in struct curl_data.
static size_t curl_callback (void *contents, size_t size, size_t nmemb, void *userp) {
	size_t real_size = size * nmemb;                  // Calculate buffer size.
	struct curl_data *p = (struct curl_data *) userp; // Cast pointer to custom type.

	// Expand buffer.
	p->payload = (char *)realloc(p->payload, p->size + real_size + 1);
	if (p->payload == NULL) {
		pam_log(LOG_PRIORITY_ERROR, "Failed to expand buffer in curl_callback.");
		free(p->payload);

		return CURL_CALLBACK_ERR;
	}

	// Populate buffer with contents.
	memcpy(&(p->payload[p->size]), contents, real_size);

	// Set new buffer size.
	p->size += real_size;

	// Ensure null termination.
	p->payload[p->size] = 0;

	// Callback expects size of buffer to be returned.
	return real_size;
}

int http_request(const char * method, const char* endpoint, char *req_body, json_t **result_j) {
	pam_log(LOG_PRIORITY_DEBUG, "Initializing HTTP request %s %s", method, endpoint);
	pam_log(LOG_PRIORITY_DEBUG, "HTTP request body: %s", req_body);

	if (result_j != NULL) {
		*result_j = NULL; // Enables callers to check if errors occured.
	}

	CURL* curl_handle = curl_easy_init();
	if (!curl_handle) {
		pam_log(LOG_PRIORITY_ERROR, "Unable to initialize cURL");
		return 0;
	}

	char url[2 * MAX_FLAG_SIZE];
	curl_easy_setopt(curl_handle, CURLOPT_URL, join_url(url, flag_opa_url, endpoint));
	curl_easy_setopt(curl_handle, CURLOPT_CUSTOMREQUEST, method);
	curl_easy_setopt(curl_handle, CURLOPT_NOPROGRESS, 1);
	curl_easy_setopt(curl_handle, CURLOPT_FAILONERROR, 1);
	curl_easy_setopt(curl_handle, CURLOPT_TIMEOUT, 5); // Set a 5 second timeout.

	// Set up the data object which will be populated by the callback.
	struct curl_data resp_data;
	resp_data.payload = (char *) malloc(1); // This will be realloced by libcurl.
	resp_data.size = 0;                     // Start with an empty payload.

	// Set headers.
	struct curl_slist *headers = NULL;
	// The request body is JSON.
	headers = curl_slist_append(headers, "Content-Type: application/json");
	// The response body can be JSON.
	headers = curl_slist_append(headers, "Accept: application/json");

	// Set the request body JSON.
	// This has the side effect of setting request headers to default, undesired values.
	// Ensure that proper headers are set afterwards.
	if (req_body != NULL) {
		curl_easy_setopt(curl_handle, CURLOPT_POSTFIELDS, req_body);
	}

	// Specify that the data should be written to our object.
	curl_easy_setopt(curl_handle, CURLOPT_WRITEDATA, (void *)&resp_data);

	// Specify that our callback should be used to write the data.
	curl_easy_setopt(curl_handle, CURLOPT_WRITEFUNCTION, curl_callback);

	// Perform a synchronous request.
	CURLcode resp_code = curl_easy_perform(curl_handle);

	// Clean up request objects.
	if (req_body != NULL) { // Caller expects req_body to be freed.
		free(req_body);
	}
	curl_easy_cleanup(curl_handle); // Clean up curl objects.
	curl_slist_free_all(headers);	// Clean up headers.

	if (resp_code != CURLE_OK) {
		pam_log(LOG_PRIORITY_ERROR,
			"HTTP request failed with error: %s", curl_easy_strerror(resp_code));

		free(resp_data.payload);
		return resp_code;
	}

	pam_log(LOG_PRIORITY_DEBUG, "HTTP request complete, libcURL returned with %d.", resp_code);
	pam_log(LOG_PRIORITY_DEBUG, "HTTP response body: %s", resp_data.payload);

	// Only try to parse JSON if the caller has provided an object for it.
	if (result_j != NULL) {
		// Define object to store JSON errors in.
		json_error_t error;
		json_t *resp_body_j = json_loads(resp_data.payload, 0, &error);

		// The response data is not needed anymore.
		free(resp_data.payload);

		if (!resp_body_j) {
			pam_log(LOG_PRIORITY_ERROR,
				"Error parsing JSON on line %d: %s", error.line, error.text);
			return resp_code;
		}

		if (!json_is_object(resp_body_j)) {
			return json_error_ret_int(resp_body_j, resp_code,
				"top level value of JSON response recieved is not type object");
		}

		json_t *result_j_actual = json_object_get(resp_body_j, "result");
		if (!json_is_object(result_j_actual)) {
			return json_info_ret_int(resp_body_j, resp_code,
				"Value of field 'result' does not have type object in JSON response. "
				"Please ensure that your endpoint flag '%s' matches your package path.", endpoint);
		}

		// Copy the result into a new object and destroy everything else.
		*result_j = json_deep_copy(result_j_actual);
		json_decref(resp_body_j);
	}

	return resp_code;
}