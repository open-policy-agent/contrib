```mermaid
sequenceDiagram

Actor->+Kong: Request resource

```
note right of Actor
Example Request:
GET https://kong:8000/myService
Autorization: Bearer eyJhbGciOiJ
end note
Kong->*kong-plugin-opa: Execute plugin on access phase
kong-plugin-opa->kong-plugin-opa: Get request path, method and parse bearer token from authorization header
kong-plugin-opa->+Open Policy Agent: Get Document with Input
note right of kong-plugin-opa
Example request:
POST /v1/data/opa/examples/allow_request
Content-Type: application/json
{   "input": {
.        "token": {
.            "payload": {
.                "sub": "1234567890",
.                "role": "dev"
.            },
.        },
.        "method": "GET",
.        "path": "/myService",
.        "headers" : {
.           "Accept" : "application/json",
.           "Content-Type" : "application/json"
.        }
.    }
}
end note
Open Policy Agent->Open Policy Agent: Evaluate rule
Open Policy Agent->-kong-plugin-opa: Return response
note left of Open Policy Agent
Example response:
Content-Type: application/json
{    "result":  true    }
end note
alt no result or result = false
kong-plugin-opa->Kong: kong.reponse.exit 403
Kong->Actor: 403 Access Forbidden
else result = true
kong-plugin-opa-->Kong: nothing returned
destroy kong-plugin-opa
Kong->+Resource: Forward request
Resource->-Kong: Resource response
Kong->-Actor: Resource response
end
