#ifndef PAM_LOG_H
#define PAM_LOG_H

#include <stdarg.h>

extern const int LOG_PRIORITY_NONE;
extern const int LOG_PRIORITY_ERROR;
extern const int LOG_PRIORITY_INFO;
extern const int LOG_PRIORITY_DEBUG;

// initialize_default_log_session starts a logging session with default
// parameters. A log session must be initialized before any invocation of
// pam_log or vpam_log.
extern void
initialize_default_log_session();

// initialize_log_session starts a logging session with the given log
// priority level. A log session must be initialized before any invocation
// of pam_log or vpam_log.
extern void
initialize_log_session(const char *log_lvl);

// pam_log takes a log priority, and compares it with the session's log
// priority to determine if/how to log the message.
// fmt need not be suffixed with newline characters.
extern void
pam_log(int priority, const char *fmt, ...);

// vpam_log does the exact same thing as pam_log, but it is useful
// for calls from variadic functions.
extern void
vpam_log(int priority, const char *fmt, va_list args);

#endif