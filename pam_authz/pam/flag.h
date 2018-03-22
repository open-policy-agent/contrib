#ifndef PAM_FLAG_H
#define PAM_FLAG_H

#define MAX_FLAG_SIZE 256

// These variables are populated by initialize_flags() using arguments
// that this PAM module is called with.
extern char flag_opa_url[MAX_FLAG_SIZE];
extern char flag_pull_endpoint[MAX_FLAG_SIZE];
extern char flag_display_endpoint[MAX_FLAG_SIZE];
extern char flag_authz_endpoint[MAX_FLAG_SIZE];
extern char flag_log_level[MAX_FLAG_SIZE];

// initialize_flags populates the variables above using arguments that
// this PAM module is called with.
extern void
initialize_flags(int argc, const char **argv);

#endif