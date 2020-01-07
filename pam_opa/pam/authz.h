#ifndef PAM_AUTHZ_H
#define PAM_AUTHZ_H

#include "display.h"
#include "pull.h"
#include "sysinfo.h"

// engine_authz sends the data collected from previous cycles to the policy engine
// to make the authorization decision. It logs errors that OPA sends back, if any,
// and then returns either PAM_SUCCESS or PAM_AUTH_ERR.
extern int
engine_authz(const char *url, struct DisplayResponses display_responses,
	struct PullResponses pull_responses, struct Sysinfo sysinfo);

#endif