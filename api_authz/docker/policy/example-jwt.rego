package httpapi.authz

default allow = false

# Allow users to get their own salaries.
allow {
  some username
  input.method == "GET"
  input.path = ["finance", "salary", username]
  token.payload.user == username
  user_owns_token
}

# Allow managers to get their subordinate' salaries.
allow {
  some username
  input.method == "GET"
  input.path = ["finance", "salary", username]
  token.payload.subordinates[_] == username
  user_owns_token
}

# Allow HR members to get anyone's salary.
allow {
  input.method == "GET"
  input.path = ["finance", "salary", _]
  token.payload.hr == true
  user_owns_token
}

# Ensure that the token was issued to the user supplying it.
user_owns_token { input.user == token.payload.azp }

# Helper to get the token payload.
token = {"payload": payload} {
  [_, payload, _] := io.jwt.decode(input.token)
}
