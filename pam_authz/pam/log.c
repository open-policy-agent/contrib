#include <string.h>
#include <stdarg.h>
#include <syslog.h>

#include <stdio.h>

#include "log.h"

#define LOG_LEVEL_COUNT 4

// Application identity displayed when logging.
static const char SYSLOG_IDENTITY[] = "OPA-PAM";

// Log levels.
// These constants should be used inside the PAM module to specify logging behavior.
static const char LOG_LEVEL_NONE[]  = "none";  // Log nothing.
static const char LOG_LEVEL_ERROR[] = "error"; // Log only errors.
static const char LOG_LEVEL_INFO[]  = "info";  // Log general info.
static const char LOG_LEVEL_DEBUG[] = "debug"; // Also log to stderr, and log verbosely.

// LOG_PRIORITY is a way of assigning a number (index) to each log level.
// Higher priority means more verbose logs.
static const char* LOG_PRIORITY[LOG_LEVEL_COUNT] = {
	LOG_LEVEL_NONE,
	LOG_LEVEL_ERROR,
	LOG_LEVEL_INFO,
	LOG_LEVEL_DEBUG,
};

// Explicit log priorities make for simpler and faster implementations.
// These must match the ordering in LOG_PRIORITY.
const int LOG_PRIORITY_NONE  = 0;
const int LOG_PRIORITY_ERROR = 1;
const int LOG_PRIORITY_INFO  = 2;
const int LOG_PRIORITY_DEBUG = 3;

static const int LOG_PRIORITY_INVALID = -1;

static int log_priority_for_log_level(const char *log_lvl) {
	int i;
	for (i = 0; i < LOG_LEVEL_COUNT; i++) {
		if (strcmp(LOG_PRIORITY[i], log_lvl) == 0) {
			return i;
		}
	}

	return LOG_PRIORITY_INVALID;
}

// session_log_priority is the log priority for the current session.
// It defaults to LOG_PRIORITY_INFO.
static int session_log_priority;

void initialize_default_log_session() {
	// Set session_log_priority default.
	session_log_priority = LOG_PRIORITY_INFO;
	pam_log(LOG_PRIORITY_INFO, "Defaulted to log level %s", LOG_LEVEL_INFO);
}

void initialize_log_session(const char *log_lvl) {
	int pri = log_priority_for_log_level(log_lvl);
	if (pri != LOG_PRIORITY_INVALID) {
		session_log_priority = pri;

		pam_log(LOG_PRIORITY_INFO, "Session log level is set to %s", log_lvl);
	} else {
		pam_log(LOG_PRIORITY_ERROR, "Invalid log level defined in PAM module: %s", log_lvl);

		initialize_default_log_session();
	}
}

void pam_log(int priority, const char *fmt, ...) {
	va_list args;
	va_start(args, fmt);

	vpam_log(priority, fmt, args);

	va_end(args);
}

void vpam_log(int priority, const char *fmt, va_list args) {
	if (session_log_priority >= LOG_PRIORITY_DEBUG) {
		// LOG_PERROR will additionally pipe logs to stderr.
		openlog(SYSLOG_IDENTITY, (LOG_PERROR | LOG_CONS | LOG_PID | LOG_NDELAY), LOG_AUTH);
	} else {
		openlog(SYSLOG_IDENTITY, (LOG_CONS | LOG_PID | LOG_NDELAY), LOG_AUTH);
	}

	// Should this log be logged?
	if (session_log_priority >= priority && session_log_priority > LOG_PRIORITY_NONE) {
		if (priority == LOG_PRIORITY_ERROR) {
			vsyslog(LOG_ERR, fmt, args);
		} else {
			// Log all non error logs as syslog LOG_INFO.
			vsyslog(LOG_INFO, fmt, args);
		}
	}

	// When openlog is used in a shared object library, it must be closed
	// before the library is unloaded. Closing it here is safest.
	closelog();
}