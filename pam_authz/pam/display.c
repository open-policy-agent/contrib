#include <stdlib.h>
#include <string.h>

#include <security/pam_appl.h>
#include <jansson.h>

#include "display.h"
#include "http.h"
#include "json.h"
#include "log.h"

#define DISPLAY_STYLE_COUNT 4

static const int DISPLAY_STYLE_INVALID = -1;

// Diplay style constants used in OPA.
static const char DISPLAY_STYLE_PROMPT_ECHO_ON[]  = "prompt_echo_on";
static const char DISPLAY_STYLE_PROMPT_ECHO_OFF[] = "prompt_echo_off";
static const char DISPLAY_STYLE_TEXT_INFO[]       = "info";
static const char DISPLAY_STYLE_ERROR_MSG[]       = "error";


static const struct DisplayStyleToPamInt {
	const char *style;
	const int  pam_int;
} DISPLAY_STYLE_TO_PAM_INT[DISPLAY_STYLE_COUNT] = {
	{DISPLAY_STYLE_PROMPT_ECHO_ON,  PAM_PROMPT_ECHO_ON},
	{DISPLAY_STYLE_PROMPT_ECHO_OFF, PAM_PROMPT_ECHO_OFF},
	{DISPLAY_STYLE_ERROR_MSG,       PAM_ERROR_MSG},
	{DISPLAY_STYLE_TEXT_INFO,       PAM_TEXT_INFO},
};

static int pam_int_for_display_style(const char *style) {
	int i;
	for (i = 0; i < DISPLAY_STYLE_COUNT; i++) {
		if (strcmp(DISPLAY_STYLE_TO_PAM_INT[i].style, style) == 0) {
			return DISPLAY_STYLE_TO_PAM_INT[i].pam_int;
		}
	}

	return DISPLAY_STYLE_INVALID;
}

static struct pam_conv *get_pam_conv(pam_handle_t *pamh) {
	if (!pamh)
	return NULL;

	struct pam_conv *conv;

	pam_log(LOG_PRIORITY_DEBUG, "Retrieving struct pam_conv from the PAM handle.");

	if (pam_get_item(pamh, PAM_CONV, (const void**)&conv) != PAM_SUCCESS)
		return NULL;

	return conv;
}

static char *call_pam_conv(pam_handle_t *pamh, int msg_style, char *message) {
	// conv is always called with parameter num_msg == 1 as a way to
	// make it compatible with both Linux-PAM and Solaris' PAM,
	// see https://linux.die.net/man/3/pam_conv
	int num_msg = 1;

	// Create a struct pam_message array and populate it with a single object.
	struct pam_message *msg_array = (struct pam_message *)(malloc(sizeof(struct pam_message)));
	msg_array[0].msg_style = msg_style;
	msg_array[0].msg = message;

	// Create a struct pam_response array.
	struct pam_response *resp_array = NULL;

	struct pam_conv *obj = get_pam_conv(pamh);
	if (obj == NULL)
		return NULL;

	pam_log(LOG_PRIORITY_DEBUG,
		"Calling application-defined conversation function to display message '%s' to "
		"the user.", message);

	int conv_resp = obj->conv(num_msg,
		(const struct pam_message **)&msg_array, &resp_array, obj->appdata_ptr);
	if (conv_resp != PAM_SUCCESS) {
		pam_log(LOG_PRIORITY_ERROR,
			"Received error from application-defined conversation function: %s",
			pam_strerror(pamh, conv_resp));

		return NULL;
	}

	char *user_resp = NULL;
	if (resp_array[0].resp != NULL) {
		pam_log(LOG_PRIORITY_DEBUG, "Collected a prompt response from user.");

		user_resp = strdup(resp_array[0].resp);
		free((resp_array[0].resp));
	}

	free(resp_array);
	free(msg_array);

	return user_resp;
}

void
engine_display(pam_handle_t *pamh, const char *endpoint,
	struct DisplayResponses *display_responses_ptr) {
	// Initialize empty responses, then fill it up as the user responses come in.
	display_responses_ptr->count = 0;
	// An empty malloc here allows calling free() later without having to check anything.
	display_responses_ptr->responses = (struct DisplayResponse *)malloc(0);

	if (strcmp(endpoint, "") == 0) {
		pam_log(LOG_PRIORITY_INFO, "Display endpoint is empty; not proceeding.");
		return;
	}

	json_t *result_j;
	http_request(HTTP_METHOD_GET, endpoint, NULL, &result_j);

	if (result_j == NULL) {
		// Errors occurred during HTTP action will have ben logged there.
		// Nothing to do.
		return;
	}

	json_t *display_spec_j = json_object_get(result_j, "display_spec");
	if (!json_is_array(display_spec_j)) {
		return json_error_ret_void(result_j,
			"Value of field 'display_spec' does not have type array in JSON response");
	}

	int i;
	for (i = 0; i < json_array_size(display_spec_j); i++) {
		json_t *display_spec_elem_j, *message_j, *style_j, *key_j;

		display_spec_elem_j = json_array_get(display_spec_j, i);
		if (!json_is_object(display_spec_elem_j)) {
			return json_error_ret_void(result_j,
				"Value of %dth element in 'display_spec' does not have type object in "
				"JSON response", i);
		}

		message_j = json_object_get(display_spec_elem_j, "message");
		if (!json_is_string(message_j)) {
			return json_error_ret_void(result_j,
				"Value of 'message' in %dth element of 'display_spec' does not have type "
				"string in JSON response", i);
		}

		style_j = json_object_get(display_spec_elem_j, "style");
		if (!json_is_string(style_j)) {
			return json_error_ret_void(result_j,
				"Value of 'style' in %dth element of 'display_spec' does not have type "
				"string in JSON response", i);
		}

		int pam_style = pam_int_for_display_style(json_string_value(style_j));
		if (pam_style == DISPLAY_STYLE_INVALID) {
			pam_log(LOG_PRIORITY_ERROR,
				"Received invalid display style: %s", json_string_value(style_j));
		} else {
			// Perform conversation if display style is valid.

			char* user_resp = call_pam_conv(pamh, pam_style, (char *)json_string_value(message_j));

			if (pam_style == PAM_PROMPT_ECHO_ON || pam_style == PAM_PROMPT_ECHO_OFF) {
				key_j = json_object_get(display_spec_elem_j, "key");
				if (!json_is_string(key_j)) {
					return json_error_ret_void(result_j,
						"Value of 'key' in %dth element of 'display_spec' does not have type "
						"string in JSON response", i);
				}

				// Extend the responses array.
				display_responses_ptr->responses = (struct DisplayResponse *)realloc(
					display_responses_ptr->responses,
					((display_responses_ptr->count)+1) * sizeof(struct DisplayResponse));

				if (display_responses_ptr->responses == NULL) {
					pam_log(LOG_PRIORITY_ERROR,
						"Unable to allocate memory to store display responses.");
				}

				display_responses_ptr->responses[display_responses_ptr->count].key = strdup(
					json_string_value(key_j));

				if (display_responses_ptr->responses[display_responses_ptr->count].key == NULL) {
					pam_log(LOG_PRIORITY_ERROR,
						"Unable to allocate memory to store display responses.");
				}

				display_responses_ptr->responses[display_responses_ptr->count].input = user_resp;
				display_responses_ptr->count++;
			}
		}
	}

	json_decref(result_j); // Clean up.
}

void free_display_responses(struct DisplayResponses *display_responses_ptr) {
	int i;
	for (i = 0; i < display_responses_ptr->count; i++) {
		free((char *)display_responses_ptr->responses[i].key);
		free((char *)display_responses_ptr->responses[i].input);
	}

	free(display_responses_ptr->responses);
}