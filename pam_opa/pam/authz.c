#include <jansson.h>

#include "display.h"
#include "authz.h"
#include "http.h"
#include "json.h"
#include "log.h"
#include "pull.h"
#include "sysinfo.h"

int engine_authz(const char *endpoint, struct DisplayResponses display_responses,
	struct PullResponses pull_responses, struct Sysinfo sysinfo) {
	// Create JSON objects.
	json_t *req_body_j          = json_object();
	json_t *input_j             = json_object();
	json_t *display_responses_j = json_object();
	json_t *pull_responses_j    = json_object();
	json_t *files_j             = json_object();
	json_t *env_vars_j          = json_object();
	json_t *sysinfo_j           = json_object();

	int i;

	// Construct display_responses object.
	for (i = 0; i < display_responses.count; i++) {
		// Try to add user's response values by specified key to the object.
		if (json_object_set_new(display_responses_j, display_responses.responses[i].key,
			json_string(display_responses.responses[i].input))) {

			pam_log(LOG_PRIORITY_ERROR, "Could not set display_responses '%s' to '%s'",
				display_responses.responses[i].key, display_responses.responses[i].input);
		}
	}

	// Add display_responses object to input object.
	if (json_object_set_new(input_j, "display_responses", display_responses_j)) {
		pam_log(LOG_PRIORITY_ERROR, "Could not set input 'display_responses'");
	}

	// Construct files object.
	for (i = 0; i < pull_responses.file_count; i++) {
		// A jansson stealing function must not be used	used here because
		// the contents are part of an object that has its own freeing process.
		if (json_object_set(files_j, pull_responses.files[i].name,
			pull_responses.files[i].contents)) {

			pam_log(LOG_PRIORITY_ERROR, "Could not set files '%s' to '%s'",
				pull_responses.files[i].name, pull_responses.files[i].contents);
		}
	}

	// Add files object to pull_responses object.
	if (json_object_set_new(pull_responses_j, "files", files_j)) {
		pam_log(LOG_PRIORITY_ERROR, "Could not set pull_responses 'files'");
	}

	// Construct env_vars object.
	for (i = 0; i < pull_responses.env_var_count; i++) {
		if (json_object_set_new(env_vars_j, pull_responses.env_vars[i].name,
			json_string(pull_responses.env_vars[i].value))) {

			pam_log(LOG_PRIORITY_ERROR, "Could not set env_vars '%s' to '%s'",
				pull_responses.env_vars[i].name, pull_responses.env_vars[i].value);
		}
	}

	// Add env_vars object to pull_responses object.
	if (json_object_set_new(pull_responses_j, "env_vars", env_vars_j)) {
		pam_log(LOG_PRIORITY_ERROR, "Could not set pull_responses 'env_vars'");
	}

	// Add pull_responses object to input object.
	if (json_object_set_new(input_j, "pull_responses", pull_responses_j)) {
		pam_log(LOG_PRIORITY_ERROR, "Could not set input 'pull_responses'");
	}

	// Construct sysinfo object.
	for (i = 0; i < sysinfo.count; i++) {
		if (json_object_set_new(sysinfo_j, sysinfo.items[i].id,
			json_string(sysinfo.items[i].value))) {

			pam_log(LOG_PRIORITY_ERROR, "Could not set sysinfo '%s' to '%s'",
				sysinfo.items[i].id, sysinfo.items[i].value);
		}
	}

	// Add sysinfo object to input object.
	if (json_object_set_new(input_j, "sysinfo", sysinfo_j)) {
		pam_log(LOG_PRIORITY_ERROR, "Could not set input 'sysinfo'");
	}

	// Add input object to top-level request object.
	if (json_object_set_new(req_body_j, "input", input_j)) {
		pam_log(LOG_PRIORITY_ERROR, "Could not set req_body 'input'");
	}

	json_t *result_j;
	http_request(HTTP_METHOD_POST, endpoint, json_dumps(req_body_j, JSON_COMPACT), &result_j);

	// Only the top level JSON object needs to be cleaned up, because all other
	// objects should be added to it via stealing functions.
	json_decref(req_body_j);

	if (result_j == NULL) {
		// Errors occurred during HTTP action will have ben logged there.
		// Nothing to do.
		return PAM_AUTH_ERR;
	}

	json_t *allow_j = json_object_get(result_j, "allow");
	if (!json_is_boolean(allow_j)) {
		return json_error_ret_int(result_j, PAM_AUTH_ERR,
			"Value of field 'allow' does not have type boolean in JSON response");	
	}

	int decision = PAM_AUTH_ERR;
	if (json_boolean_value(allow_j)) {
		decision = PAM_SUCCESS;
	} 

	json_t *errors_j = json_object_get(result_j, "errors");
	if (!json_is_array(errors_j)) {
		return json_error_ret_int(result_j, decision,
			"Value of field 'errors' does not have type array in JSON response");	
	}

	for (i = 0; i < json_array_size(errors_j); i++) {
		json_t *error_msg_j = json_array_get(errors_j, i);
		if (!json_is_string(error_msg_j)) {
			return json_error_ret_int(result_j, decision,
			"Value of %dth element of 'errors' does not have type string in JSON response", i);
		}

		pam_log(LOG_PRIORITY_ERROR, "Received authz error log from OPA: %s", json_string_value(error_msg_j));
	}

	json_decref(result_j);
	return decision;
}