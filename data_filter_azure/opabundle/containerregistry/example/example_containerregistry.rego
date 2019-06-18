package containerregistry.example

default allow = false

allow {
    data[input.registry][input.user][input.type][input.resourceName].actions[_] == input.action
}

allow {
    data[input.registry][input.user][input.type][input.resourceName].actions[_] == "*"
}

allow {
    data[input.registry][input.user][input.type]["*"].actions[_] == input.action
}

allow {
    data[input.registry][input.user][input.type]["*"].actions[_] == "*"
}