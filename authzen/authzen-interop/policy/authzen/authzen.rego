package authzen

import rego.v1

default allow["decision"] := false

# Anyone can read users or todos
allow["decision"] if input.action.name == "can_read_user"

allow["decision"] if input.action.name == "can_read_todos"

allow["decision"] if {
	input.action.name == "can_create_todo"
	can_create
}

allow["decision"] if {
	input.action.name == "can_update_todo"
	can_update
}

allow["decision"] if {
	input.action.name == "can_delete_todo"
	can_delete
}

user_is_admin if "admin" in data.users[input.subject.id].roles

user_is_evil_genius if "evil_genius" in data.users[input.subject.id].roles

user_is_editor if "editor" in data.users[input.subject.id].roles

user_is_owner if input.resource.properties.ownerID == data.users[input.subject.id].email

can_create if user_is_admin

can_create if user_is_editor

can_update if user_is_evil_genius

can_update if {
	user_is_editor
	user_is_owner
}

can_delete if user_is_admin

can_delete if {
	user_is_editor
	user_is_owner
}
