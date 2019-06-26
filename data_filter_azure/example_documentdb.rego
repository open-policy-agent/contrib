package documentdb.example

default allow = false

allow {
    allowed[_]
}

allowed[user] {
    data.permissions[user]
    user.registry == input.registry
    user.user == input.user
    user.map[_].type == input.type
    user.map[_].name == input.resourceName
    contains(user.map[_].actions, input.action)
}

allowed[user] {
    data.permissions[user]
    user.registry == input.registry
    user.user == input.user
    user.map[_].type == input.type
    user.map[_].name == "*"
    contains(user.map[_].actions, input.action)
}

allowed[user] {
    data.permissions[user]
    user.registry == input.registry
    user.user == input.user
    user.map[_].type == input.type
    user.map[_].name == input.resourceName
    contains(user.map[_].actions, "*")
}

allowed[user] {
    data.permissions[user]
    user.registry == input.registry
    user.user == input.user
    user.map[_].type == input.type
    user.map[_].name == "*"
    contains(user.map[_].actions, "*")
}