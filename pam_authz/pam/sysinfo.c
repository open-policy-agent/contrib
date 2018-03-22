#include <stdlib.h>
#include <string.h>

#include <security/pam_appl.h>

#include "sysinfo.h"
#include "log.h"

#define MAX_SYSINFO_ID_LENGTH 256
#define SYSINFO_COUNT 4

static const char SYSINFO_ID_PAM_USER   [MAX_SYSINFO_ID_LENGTH] = "pam_username";
static const char SYSINFO_ID_PAM_SERVICE[MAX_SYSINFO_ID_LENGTH] = "pam_service";
static const char SYSINFO_ID_PAM_RUSER  [MAX_SYSINFO_ID_LENGTH] = "pam_req_username";
static const char SYSINFO_ID_PAM_RHOST  [MAX_SYSINFO_ID_LENGTH] = "pam_req_hostname";

static const struct SysinfoIdToPamInt {
	const char *id;
	int        pam_int;
} SYSINFO_ID_TO_PAM_INT[SYSINFO_COUNT] = {
	{SYSINFO_ID_PAM_USER,    PAM_USER},
	{SYSINFO_ID_PAM_SERVICE, PAM_SERVICE},
	{SYSINFO_ID_PAM_RUSER,   PAM_RUSER},
	{SYSINFO_ID_PAM_RHOST,   PAM_RHOST},
};

static const int SYSINFO_INVALID = -1;

static int get_pam_int_for_sysinfo_id(const char *sysinfo_id) {
	int i;
	for (i = 0; i < SYSINFO_COUNT; i++) {
		if (strcmp(SYSINFO_ID_TO_PAM_INT[i].id, sysinfo_id) == 0) {
			return SYSINFO_ID_TO_PAM_INT[i].pam_int;
		}
	}

	return SYSINFO_INVALID;	
}

static char *pam_get_item_string(pam_handle_t *pamh, int item_type) {
	char *str;
	if (pam_get_item(pamh, item_type, (const void**)&str) != PAM_SUCCESS)
		return NULL;

	if (str == NULL)
		return NULL;

	return strdup(str);
}

void load_sysinfo(pam_handle_t *pamh, struct Sysinfo *sysinfo_ptr) {
	int i;

	sysinfo_ptr->items = calloc(SYSINFO_COUNT, sizeof(struct Sysinfo));
	sysinfo_ptr->count = SYSINFO_COUNT;

	for (i = 0; i < SYSINFO_COUNT; i++) {
		if (sysinfo_ptr->items == NULL) {
			pam_log(LOG_PRIORITY_ERROR,
				"Unable to allocate memory to store sys responses.");
		}

		sysinfo_ptr->items[i].id = strdup(
			SYSINFO_ID_TO_PAM_INT[i].id);

		if (sysinfo_ptr->items[i].id == NULL) {
			pam_log(LOG_PRIORITY_ERROR,
				"Unable to allocate memory to store sys responses.");
		}

		char *sysinfo_value = pam_get_item_string(pamh, SYSINFO_ID_TO_PAM_INT[i].pam_int);
		if (sysinfo_value == NULL) {
			sysinfo_value = strdup("");
		}
		pam_log(LOG_PRIORITY_DEBUG, "Loaded sysinfo %s: %s",
			SYSINFO_ID_TO_PAM_INT[i].id, sysinfo_value);

		sysinfo_ptr->items[i].value = sysinfo_value;
	}
}

void free_sysinfo(struct Sysinfo *sysinfo_ptr) {
	int i;
	for (i = 0; i < sysinfo_ptr->count; i++) {
		free(sysinfo_ptr->items[i].id);
		free(sysinfo_ptr->items[i].value);
	}

	free(sysinfo_ptr->items);
}