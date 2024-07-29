package example

import rego.v1

policies := [
	{
		"subjects": ["group:users"],
		"resource": ["posts"],
		"action": "GET",
		"condition": [
			{
				"field": "author",
				"ref": "user",
				"operator": "equal",
			},
			{
				"field": "id",
				"value": "post1",
				"operator": "equal",
			},
		],
	},
	{
		"subjects": ["group:admin"],
		"resource": ["posts"],
		"action": "GET",
		"condition": [{
			"field": "clearance",
			"ref": "level",
			"value": 9,
			"operator": "greater_than",
		}],
	},
]

### Entry point to the policy.
### Example with conditions being interpreted.

allow if {
	some p in policies
	p.action = input.method
	p.resource = input.path
	some r in data.elastic[p.resource[0]]
	all_conditions_true(p, r)
}

all_conditions_true(p, r) if {
	not any_condition_false(p, r)
}

any_condition_false(p, r) if {
	some c in p.condition
	not eval(r, c.operator, c)
}

### Equality operator

# OPA Query: "bob" = data.elastic.posts[_].author; "post1" = data.elastic.posts[_].id
# ES  Query: {name:author value:bob boost:<nil> queryName:TermQuery}
eval(r, "equal", c) if {
	r[c.field] == input[c.ref] # regal ignore:external-reference
}

# OPA Query: "bob" = data.elastic.posts[_].author; "post1" = data.elastic.posts[_].id
# ES  Query: {name:id value:post1 boost:<nil> queryName:TermQuery}
eval(r, "equal", c) if {
	r[c.field] == c.value
}

### Greater Than operator

# OPA Query: gt(data.elastic.posts[_].clearance, 9)
# ES  Query: {name:clearance from:9 to:<nil> timeZone: includeLower:false
#             includeUpper:true boost:<nil> queryName: format: relation:}
eval(r, "greater_than", c) if {
	r[c.field] > c.value
}

### Sample Input to OPA.
# {
#     "method": "GET",
#     "path": ["posts"],
#     "user": "bob"
# }
### Sample Output from Elasticsearch.
# {
#   "result": [
#     {
#       "id": "post1",
#       "author": "bob",
#       "message": "My first post",
#       "department": "dev",
#       "email": "bob@abc.com",
#       "clearance": 2,
#       "action": "read",
#       "resource": "",
#       "conditions": [],
#       "likes": [],
#       "followers": [],
#       "stats": []
#     },
#     {
#       "id": "post5",
#       "author": "ben",
#       "message": "Hii from Ben",
#       "department": "ceo",
#       "email": "ben@opa.com",
#       "clearance": 10,
#       "action": "read",
#       "resource": "",
#       "conditions": [],
#       "likes": [],
#       "followers": [],
#       "stats": []
#     },
#     {
#       "id": "post8",
#       "author": "ben",
#       "message": "This is OPA's time",
#       "department": "ceo",
#       "email": "ben@opa.com",
#       "clearance": 10,
#       "action": "read",
#       "resource": "",
#       "conditions": [],
#       "likes": [],
#       "followers": [],
#       "stats": []
#     }
#   ]
# }
