#include <jansson.h>

#include "json.h"
#include "log.h"

void json_error_ret_void(json_t *j, const char *fmt, ...) {
	va_list args;
	va_start(args, fmt);

	// Log error.
	vpam_log(LOG_PRIORITY_ERROR, fmt, args);

	va_end(args);

	// Free up JSON object memory.
	json_decref(j);
}

int json_error_ret_int(json_t *j, int ret, const char *fmt, ...) {
	va_list args;
	va_start(args, fmt);

	// Log error.
	vpam_log(LOG_PRIORITY_ERROR, fmt, args);

	va_end(args);

	// Free up JSON object memory.
	json_decref(j);

	return ret;
}

int json_info_ret_int(json_t *j, int ret, const char *fmt, ...) {
	va_list args;
	va_start(args, fmt);

	// Log info.
	vpam_log(LOG_PRIORITY_INFO, fmt, args);

	va_end(args);

	// Free up JSON object memory.
	json_decref(j);

	return ret;
}