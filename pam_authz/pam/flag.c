#include <stdlib.h>
#include <string.h>

#include "flag.h"
#include "log.h"

#define FLAG_COUNT 5

// Valid flags that can be used as arguments to this PAM module inside /etc/pam.d/ files.
static const char FLAG_STR_OPA_URL[]          = "url";
static const char FLAG_STR_PULL_ENDPOINT[]    = "pull_endpoint";
static const char FLAG_STR_DISPLAY_ENDPOINT[] = "display_endpoint";
static const char FLAG_STR_AUTHZ_ENDPOINT[]   = "authz_endpoint";
static const char FLAG_STR_LOG_LEVEL[]        = "log_level";

// Default values of exposed variables.
char flag_opa_url[MAX_FLAG_SIZE]          = "";
char flag_pull_endpoint[MAX_FLAG_SIZE]    = ""; 
char flag_display_endpoint[MAX_FLAG_SIZE] = "";
char flag_authz_endpoint[MAX_FLAG_SIZE]   = "";
char flag_log_level[MAX_FLAG_SIZE]        = "";

static const struct FlagStrToVar {
	const char *flag_str;
	const char *var;
} FLAG_STR_TO_VAR[FLAG_COUNT] = {
	{FLAG_STR_OPA_URL,          flag_opa_url},
	{FLAG_STR_PULL_ENDPOINT,    flag_pull_endpoint},
	{FLAG_STR_DISPLAY_ENDPOINT, flag_display_endpoint},
	{FLAG_STR_AUTHZ_ENDPOINT,   flag_authz_endpoint},
	{FLAG_STR_LOG_LEVEL,        flag_log_level},
};

static const char *var_for_flag_str(const char* flag_str) {
	int i;
	for (i = 0; i < FLAG_COUNT; i++) {
		if (strcmp(FLAG_STR_TO_VAR[i].flag_str, flag_str) == 0) {
			return FLAG_STR_TO_VAR[i].var;
		}
	}

	return NULL;
}

void initialize_flags(int argc, const char **argv) {
	int i;
	for (i = 0; i < argc; i++) {
		pam_log(LOG_PRIORITY_INFO, "Parsing arg: %s", argv[i]);

		char *arg, *tofree;
		arg = tofree = strdup(argv[i]);
		if (arg == NULL) {
			pam_log(LOG_PRIORITY_ERROR, "Error reading arg: insufficient memory.");
			continue;
		}

		char *flag  = strsep(&arg, "=");
		char *value = strsep(&arg, "=");

		if (arg != NULL) {
			pam_log(LOG_PRIORITY_ERROR, "Got arg with multiple '=': %s", argv[i]);
			free(tofree);
			continue;
		}

		if (flag == NULL || value == NULL) {
			pam_log(LOG_PRIORITY_ERROR, "Could not parse arg: %s", argv[i]);
			free(tofree);
			continue;
		}

		const char *flag_var = var_for_flag_str(flag);
		if (flag_var == NULL) {
			pam_log(LOG_PRIORITY_ERROR, "Got unknown flag: %s", flag);
			free(tofree);
			continue;
		}

		strcpy((char *)flag_var, value);
		free(tofree);
	}
}