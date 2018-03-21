#ifndef PAM_SYSINFO_H
#define PAM_SYSINFO_H

#include <security/pam_appl.h>

// struct Sysinfo holds key-value system data.
struct SysinfoItem {
	char *id;
	char *value;
};

// struct Sysinfo holds all the system information collected.
struct Sysinfo {
	int                count;
	struct SysinfoItem *items;
};

extern void
load_sysinfo(pam_handle_t *pamh, struct Sysinfo *sysinfo_ptr);

extern void
free_sysinfo(struct Sysinfo *sysinfo_ptr);

#endif