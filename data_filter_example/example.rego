package example

import rego.v1

# data.example.allow == true
#
# {"input": {"method": "GET", "path": ...}, "unknowns": ["data.posts"]}

allow if {
	input.method == "POST"
	input.path == ["posts"]
	input.subject.user # authenticated users can create posts
}

allow if {
	some post_id
	input.method == "GET"
	input.path = ["posts", post_id]

	some post in allowed
	post.id == post_id
}

allow if {
	input.method == "GET"
	input.path == ["posts"]
	count(allowed) > 0
}

allowed contains post if {
	some post in data.posts
	post.author == input.subject.user
}

allowed contains post if {
	some post in data.posts
	post.department in input.subject.departments
	post.clearance_level <= input.subject.clearance_level
}

allowed contains post if {
	some post in data.posts
	post.department == "company"
}
