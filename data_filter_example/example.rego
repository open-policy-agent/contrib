package example

# data.example.allow == true
#
# {"input": {"method": "GET", "path": ...}, "unknowns": ["data.posts"]}

allow {
    input.method = "POST"
    input.path = ["posts"]
    input.subject.user  # authenticated users can create posts
}

allow {
    input.method = "GET"
    input.path = ["posts", post_id]
    allowed[post]
    post.id = post_id
}

allow {
    input.method = "GET"
    input.path = ["posts"]
    allowed[_]
}

allowed[post] {
    data.posts[post]
    post.author = input.subject.user
}

allowed[post] {
    data.posts[post]
    post.department = input.subject.departments[_]
    post.clearance_level <= input.subject.clearance_level
}

allowed[post] {
    data.posts[post]
    post.department = "company"
}
