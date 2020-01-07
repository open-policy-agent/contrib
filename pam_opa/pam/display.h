#ifndef PAM_DISPLAY_H
#define PAM_DISPLAY_H

#include <security/pam_appl.h>

// struct DisplayResponse holds the input entered by the user for a prompt,
// along with the key associated with that prompt.
struct DisplayResponse {
	const char *key;
	const char *input;
};

// struct DisplayResponses holds an array of struct DisplayResponse,
// along with the the size of said array.
struct DisplayResponses {
	int count;
	struct DisplayResponse *responses;
};

// engine_display communicates with the policy engine to determine what to display
// to the user. It performs the required display operations, including user prompts,
// and stores data collected from the user in *display_reponses_ptr.
//
// free_display_responses should be called when this data is no longer required.
extern void
engine_display(pam_handle_t *pamh, const char *url, struct DisplayResponses *display_responses_ptr);

// free_display_responses ensures that all allocated data within
// *display_responses_ptr is freed. *display_responses_ptr itself is not freed.
extern void
free_display_responses(struct DisplayResponses *display_responses_ptr);

#endif