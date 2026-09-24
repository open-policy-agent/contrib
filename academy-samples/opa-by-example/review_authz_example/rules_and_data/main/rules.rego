package main 
import future.keywords  # this is a special import to enable some newer keywords
import data.policy.role  # import the role package

default allow_review := false
default allow_delete := false

allow_review := true if {  # allow a customer with reputation >= 0 to review
  input.role == "customer"
  input.reputation >= 0
}

allow_delete := true if {  # allow a moderator to delete
  role.is_moderator == true
}

