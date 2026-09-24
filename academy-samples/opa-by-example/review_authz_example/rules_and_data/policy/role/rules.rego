package policy.role 
import future.keywords

is_moderator := true if {
  input.role == "admin"
}

is_moderator := true if {
  input.reputation >= 100
}

is_customer := true if {
 input.role == "customer"
}

