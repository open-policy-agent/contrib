package example1

is_admin { 
    input_group := input.groups[_]
    input_group == "admin"
}