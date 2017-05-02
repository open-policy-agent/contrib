package httpapi.authz

# Datasource for management
#import data.manager
# HTTP API request
manager = {"alice": "bob", "charlie": "bob"}

import input as http_api
# http_api = {
#   "user": "dave",
#   "path": ["finance", "salary", "alice"],
#   "method": "GET"
# }

default allow = false
allow {
  http_api.method = "GET"
  http_api.path = ["finance", "salary", username]
  username = http_api.user
}

allow {
  http_api.method = "GET"
  http_api.path = ["finance", "salary", username]
  manager[username] = http_api.user
}
