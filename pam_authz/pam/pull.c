#include <stdlib.h>
#include <string.h>

#include <security/pam_appl.h>
#include <jansson.h>

#include "pull.h"
#include "http.h"
#include "json.h"
#include "log.h"

static void load_files(struct PullResponses *pull_responses_ptr, json_t *files_j) {
	int i;
	for (i = 0; i < json_array_size(files_j); i++) {
		json_t *filename_j = json_array_get(files_j, i);
		if (!json_is_string(filename_j)) {
			return json_error_ret_void(files_j,
				"Value of %dth element in 'files' does not have type string");
		}

		// Extend the files array.
		pull_responses_ptr->files = (struct File *)realloc(
			pull_responses_ptr->files,
			((pull_responses_ptr->file_count)+1) * sizeof(struct File));

		if (pull_responses_ptr->files == NULL) {
			pam_log(LOG_PRIORITY_ERROR,
				"Unable to allocate memory to store pull responses.");
		}

		pull_responses_ptr->files[pull_responses_ptr->file_count].name = strdup(
			json_string_value(filename_j));

		if (pull_responses_ptr->files[pull_responses_ptr->file_count].name == NULL) {
			pam_log(LOG_PRIORITY_ERROR,
				"Unable to allocate memory to store pull responses.");
		}

		json_error_t error;
		json_t *file_contents_j = json_load_file(json_string_value(filename_j), 0, &error);
		if (file_contents_j == NULL) {
			pam_log(LOG_PRIORITY_ERROR, "Error loading JSON file: %s", error.text);
		}

		char *file_json = json_dumps(file_contents_j, JSON_INDENT(2));
		pam_log(LOG_PRIORITY_DEBUG, "Loaded JSON from file %s:\n%s",
			json_string_value(filename_j), file_json);
		free(file_json);

		pull_responses_ptr->files[pull_responses_ptr->file_count].contents = file_contents_j;
		pull_responses_ptr->file_count++;
	}
}

static void load_env_vars(struct PullResponses *pull_responses_ptr, json_t *env_vars_j) {
	int i;
	for (i = 0; i < json_array_size(env_vars_j); i++) {
		json_t *env_var_name_j = json_array_get(env_vars_j, i);
		if (!json_is_string(env_var_name_j)) {
			return json_error_ret_void(env_vars_j,
				"Value of %dth element in 'env_vars' does not have type string");
		}

		// Extend the env_vars array.
		pull_responses_ptr->env_vars = (struct EnvVar *)realloc(
			pull_responses_ptr->env_vars,
			((pull_responses_ptr->env_var_count)+1) * sizeof(struct EnvVar));

		if (pull_responses_ptr->env_vars == NULL) {
			pam_log(LOG_PRIORITY_ERROR,
				"Unable to allocate memory to store pull responses.");
		}

		pull_responses_ptr->env_vars[pull_responses_ptr->env_var_count].name = strdup(
			json_string_value(env_var_name_j));

		if (pull_responses_ptr->env_vars[pull_responses_ptr->env_var_count].name == NULL) {
			pam_log(LOG_PRIORITY_ERROR,
				"Unable to allocate memory to store pull responses.");
		}

		char *env_var_value = getenv(json_string_value(env_var_name_j));

		pam_log(LOG_PRIORITY_DEBUG, "Loaded environment variable %s: %s",
			json_string_value(env_var_name_j), env_var_value);

		if (env_var_value != NULL) {
			pull_responses_ptr->env_vars[pull_responses_ptr->env_var_count].value = strdup(
			env_var_value);		
		} else {
			// Make this safe to assign and free.
			pull_responses_ptr->env_vars[pull_responses_ptr->env_var_count].value = strdup("");
		}
		
		pull_responses_ptr->env_var_count++;
	}
}

void engine_pull(const char *endpoint, struct PullResponses *pull_responses_ptr) {
	// Initialize empty responses, then fill it up as the user responses come in.
	pull_responses_ptr->file_count    = 0;
	pull_responses_ptr->env_var_count = 0;
	// An empty malloc here allows calling free() later without having to check anything.
	pull_responses_ptr->files    = (struct File *)malloc(0);
	pull_responses_ptr->env_vars = (struct EnvVar *)malloc(0);

	if (strcmp(endpoint, "") == 0) {
		pam_log(LOG_PRIORITY_INFO, "Pull endpoint is empty; not proceeding.");
		return;
	}

	// Get specification of what to pull from OPA.
	json_t *result_j;
	http_request(HTTP_METHOD_GET, endpoint, NULL, &result_j);	

	if (result_j == NULL) {
		// Errors occurred during HTTP action will have ben logged there.
		// Nothing to do.
		return;
	}

	json_t *files_j = json_object_get(result_j, "files");
	if (!json_is_array(files_j)) {
		return json_error_ret_void(result_j,
			"Value of field 'files' does not have type array in JSON response");
	}

	// Load requested JSON files.
	load_files(pull_responses_ptr, files_j);	

	json_t *env_vars_j = json_object_get(result_j, "env_vars");
	if (!json_is_array(env_vars_j)) {
		return json_error_ret_void(result_j,
			"Value of field 'env_vars' does not have type object in JSON response");
	}

	// Load requested environment variables.
	load_env_vars(pull_responses_ptr, env_vars_j);

	json_decref(result_j); // Clean up.
}

static void free_files(struct PullResponses *pull_responses_ptr) {
	int i;
	for (i = 0; i < pull_responses_ptr->file_count; i++) {
		free(pull_responses_ptr->files[i].name);
		json_decref(pull_responses_ptr->files[i].contents);
	}

	free(pull_responses_ptr->files);
}

static void free_env_vars(struct PullResponses *pull_responses_ptr) {
	int i;
	for (i = 0; i < pull_responses_ptr->env_var_count; i++) {
		free(pull_responses_ptr->env_vars[i].name);
		free(pull_responses_ptr->env_vars[i].value);
	}

	free(pull_responses_ptr->env_vars);
}

void free_pull_responses(struct PullResponses *pull_responses_ptr) {
	free_files(pull_responses_ptr);
	free_env_vars(pull_responses_ptr);
}