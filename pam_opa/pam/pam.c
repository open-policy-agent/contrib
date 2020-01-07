#include <stdlib.h>
#include <string.h>

#include <security/pam_appl.h>
#include <security/pam_modules.h>

#include "authz.h"
#include "display.h"
#include "flag.h"
#include "log.h"
#include "pull.h"
#include "sysinfo.h"

#ifndef PAM_EXTERN
#define PAM_EXTERN
#endif

static void initialize(int argc, const char **argv) {
	// Log level desired by the PAM module is not yet known,
	// so initialize using defaults instead.
	initialize_default_log_session();

	// Load flags from the PAM module.
	initialize_flags(argc, argv);

	// Renew the log session using desired log level.
	initialize_log_session(flag_log_level);
}

PAM_EXTERN int
pam_sm_authenticate(pam_handle_t *pamh, int flags, int argc, const char **argv) {
	pam_log(LOG_PRIORITY_DEBUG, "Application invoked pam_sm_authenticate.");

	initialize(argc, argv);

	struct DisplayResponses display_responses;
	struct PullResponses pull_responses;
	struct Sysinfo sysinfo;

	pam_log(LOG_PRIORITY_DEBUG, "Commencing display cycle.");
	engine_display(pamh, flag_display_endpoint, &display_responses);

	pam_log(LOG_PRIORITY_DEBUG, "Commencing pull cycle.");
	engine_pull(flag_pull_endpoint, &pull_responses);

	pam_log(LOG_PRIORITY_DEBUG, "Collecting system information.");
	load_sysinfo(pamh, &sysinfo);

	pam_log(LOG_PRIORITY_DEBUG, "Commencing authz cycle.");
	int authz = engine_authz(flag_authz_endpoint,
		display_responses, pull_responses, sysinfo);

	pam_log(LOG_PRIORITY_DEBUG, "Freeing allocated data.");
	free_display_responses(&display_responses);
	free_pull_responses(&pull_responses);
	free_sysinfo(&sysinfo);

	return authz;
}

PAM_EXTERN int
pam_sm_acct_mgmt(pam_handle_t *pamh, int flags, int argc, const char **argv) {
	pam_log(LOG_PRIORITY_DEBUG, "Application invoked pam_sm_acct_mgmt.");

	// Do exactly the same thing for both auth and account invocations.
	pam_sm_authenticate(pamh, flags, argc, argv);

	return PAM_SUCCESS;
}

PAM_EXTERN int
pam_sm_setcred( pam_handle_t *pamh, int flags, int argc, const char **argv) {
	pam_log(LOG_PRIORITY_DEBUG, "Application invoked pam_sm_setcred.");
	return PAM_SUCCESS;
}

