package policy

deny["foo"] {
    input.foo == "bar"
}

deny["bar"] {
    input.bar == "foo"
}
