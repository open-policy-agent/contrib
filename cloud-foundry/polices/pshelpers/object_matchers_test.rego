package pshelpers

test_path_value_match_with_correct_path_value  {
  expected := true
  inputObject = {"random": {"stuff": "moo" }}
  path := ["random", "stuff"]
  value := "moo" 
  actual := path_value_match(inputObject, path, value)
  actual == expected
}


test_path_value_match_incorrect_path {
  expected := false
  inputObject = {"random": {"stuff": "moo" }}
  path := ["incorrect", "path"]
  value := "moo" 
  actual := path_value_match(inputObject, path, value)
  actual == expected
}

test_path_value_match_incorrect_value {
  expected := false
  inputObject = {"random": {"stuff": "moo" }}
  path := ["random", "stuff"]
  value := "incorrect" 
  actual := path_value_match(inputObject, path, value)
  actual == expected
}
