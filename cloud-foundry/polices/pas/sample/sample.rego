package main

deny["1"] {
  input.path == "blah"
}

deny["we are only looking for nope"] {
  input.path == "nope"
}
