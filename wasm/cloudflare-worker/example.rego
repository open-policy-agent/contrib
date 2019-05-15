package example

allow {
    input.method = "GET"
}

allow {
    input.method = "POST"
    input.path = "/api"
}
