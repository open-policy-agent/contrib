package example2.v2

admin_groups_set := { g | g := data.admin_groups[_] }

is_admin {
  input_group := input.groups[_]
  admin_groups_set[ input_group ]
}