package httpapi.authz

import input as http_api
# http_api = {
#   "path": ["finance", "salary", "alice"],
#   "user": "alice",
#   "method": "GET",
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYWxpY2UiLCJhenAiOiJhbGljZSIsInN1Ym9yZGluYXRlcyI6W10sImhyIjpmYWxzZX0.rz3jTY033z-NrKfwrK89_dcLF7TN4gwCMj-fVBDyLoM"
# }

# io.jwt.decode takes one argument (the encoded token) and has three outputs:
# the decoded header, payload and signature, in that order. Our policy only
# cares about the payload, so we ignore the others.
token = {"payload": payload} { io.jwt.decode(http_api.token, _, payload, _) }

# Ensure that the token was issued to the user supplying it.
user_owns_token { http_api.user = token.payload.azp }

default allow = false

# Allow users to get their own salaries.
allow {
  http_api.method = "GET"
  http_api.path = ["finance", "salary", username]
  username = token.payload.user
  user_owns_token
}

# Allow managers to get their subordinate' salaries.
allow {
  http_api.method = "GET"
  http_api.path = ["finance", "salary", username]
  token.payload.subordinates[_] = username
  user_owns_token
}
