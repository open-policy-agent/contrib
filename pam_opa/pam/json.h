#include <jansson.h>

// json_error_ret logs the given error message and performs
// JSON cleanup before returning.
extern void
json_error_ret_void(json_t *j, const char *fmt, ...);

// json_error_ret_int logs the given error message and performs
// JSON cleanup before returning with the given error code.
extern int
json_error_ret_int(json_t *j, int ret, const char *fmt, ...);

// json_info_ret_int logs the given info message and performs
// JSON cleanup before returning with the given error code.
extern int
json_info_ret_int(json_t *j, int ret, const char *fmt, ...);