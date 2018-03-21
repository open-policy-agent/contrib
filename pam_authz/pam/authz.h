#ifndef PAM_AUTHZ_H
#define PAM_AUTHZ_H

#include "display.h"
#include "pull.h"
#include  "sysinfo.h"

extern int
engine_authz(const char *url, struct DisplayResponses display_responses,
	struct PullResponses pull_responses, struct Sysinfo sysinfo);

#endif