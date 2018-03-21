#ifndef PAM_PULL_H
#define PAM_PULL_H

#include <security/pam_appl.h>

#include <jansson.h>

// struct File represents a JSON file.
struct File {
	char   *name;
	json_t *contents;
};

// struct EnvVar represents an environment variable.
struct EnvVar {
	char *name;
	char *value;
};

// struct PullResponses holds data collected from the sytem.
struct PullResponses {
	int           file_count;
	struct File   *files;

	int           env_var_count;
	struct EnvVar *env_vars;
};

// engine_pull communicates with the policy engine to determine what to pull
// from the sytem. It performs the required operations, and stores data
// collected from the user in *pull_reponses_ptr.
//
// free_pull_responses should be called when this data is no longer required.
extern void
engine_pull(const char *url, struct PullResponses *pull_responses_ptr);

// free_pull_responses ensures that all allocated data within
// *pull_responses_ptr is freed. *pull_responses_ptr itself is not freed.
extern void
free_pull_responses(struct PullResponses *pull_responses_ptr);

#endif