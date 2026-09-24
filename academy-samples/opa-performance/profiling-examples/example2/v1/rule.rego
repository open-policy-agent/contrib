package example2.v1

is_admin { 
    input_group := input.groups[_]
    admin_group := data.admin_groups[_]
    input_group == admin_group
}