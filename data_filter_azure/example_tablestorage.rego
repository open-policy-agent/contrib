package example

default allow = false

allow {
    input.data[input.type][input.resourceName].actions[_] == input.action
}

allow {
    input.data[input.type][input.resourceName].actions[_] == "*"
}

allow {
    input.data[input.type]["*"].actions[_] == input.action
}

allow {
    input.data[input.type]["*"].actions[_] == "*"
}
