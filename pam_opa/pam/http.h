#ifndef PAM_HTTP_H
#define PAM_HTTP_H

#include <curl/curl.h>
#include <jansson.h>

// HTTP method constants for use in calls to http_request.
extern char HTTP_METHOD_GET[];
extern char HTTP_METHOD_POST[];

// http_request performs an HTTP request to url with method and body req_body.
// http_request will free req_body after using it.
// If resp_body is not NULL, the contents of the response are allocated to *resp_body.
extern int
http_request(const char * method, const char* url, char *req_body, json_t **result_j);

#endif