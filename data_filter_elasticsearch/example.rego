package example

allow {
    p := data.policies[_]
    p.action = input.method
    p.resource = input.path
    r := data.elastic[p.resource[0]][_]
    all_conditions_true(p, r)
}

all_conditions_true(p, r) {
    not any_condition_false(p, r)
}

any_condition_false(p, r) {
    c := p.condition[_]
    not eval(r, c.operator, c.field, c.value)
}

# allow "admin" to see all posts
allow = true {
    input.user = "admin"
}

allow = true {
    input.method = "GET"
    input.path = ["posts"]
    allowed[x]
}

allow = true {
    input.method = "GET"
    input.path = ["posts", post_id]
    allowed[x]
    x.id = post_id
}

# equality operator
eval(r, "equal", r_field, i_field) {
    r[r_field] = input[i_field]
}

eval(r, "equal", r_field, i_field) {
    r[r_field] = i_field
}

# return posts authored by input.user
allowed[x] {
    x := data.elastic.posts[_]
    #data.elastic.posts[x]
    x.author == input.user
}

# return posts with clearance level greater than 0 and less than equal to 5
# but no posts from "it"
allowed[x] {
    x := data.elastic.posts[_]
    x.clearance <= 5
    x.clearance > 0
    x.department != "it"
}

# return posts containing the term "OPA" in their message
allowed[x] {
    x := data.elastic.posts[_]
    contains(x.message, "OPA")
}

# return posts who email address matches the ".org" domain
allowed[x] {
    x := data.elastic.posts[_]
    re_match("[a-zA-Z]+@[a-zA-Z]+.org", x.email)
}

# return posts liked by input.user
allowed[x] {
    x := data.elastic.posts[_]
    y := x.likes[_]
    y.name = input.user
}

# return posts followed by input.user
allowed[x] {
    x := data.elastic.posts[_]
    y := x.followers[_]
    y.info.first = input.user
}

# return posts by authors from CA
allowed[x] {
    x := data.elastic.posts[_]
    y := x.stats[_]
    y.authorstat.authorbio.state = "CA"
}
