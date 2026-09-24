package test 
import future.keywords
import data.main

input1 := {
  "role": "customer",
  "reputation": 10
}

input2 := {
  "role": "customer",
  "reputation": -10
}

input3 := {
  "role": "guest",
  "reputation": 100
}

input4 := {
  "role": "guest",
  "reputation": 100
}

input5 := {
  "role": "admin"
}

test_allow_medium_reputation_customer_to_review if {
  main.allow_review == true with input as input1
}

test_disallow_negative_reputation_customer_to_review if {
  not main.allow_review with input as input2
}

test_disallow_non_customer_to_review if {
  not main.allow_review with input as input3
}

test_allow_high_reputation_user_to_delete if {
  main.allow_delete == true with input as input4
}

test_allow_admin_to_delete if {
  main.allow_delete == true with input as input5
}

test_disallow_medium_reputation_customer_to_delete if {
  not main.allow_delete with input as input1
}

