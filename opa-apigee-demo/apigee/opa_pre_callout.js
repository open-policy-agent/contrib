/*
The code below will call out to OPA for an authorization decision
*/

var opaUrl = "https://{{OPA_URL}}/v1/data/authz/allow";
var method = context.getVariable("request.verb");
var path = getRequestPath();
var token = context.getVariable("request.queryparam.token.1");

// Input to OPA
var inputObj = {
    'method': method,
    'path': path,
    'token': token
};

var body = JSON.stringify({input: inputObj});

var req = new Request(opaUrl, 'POST', {'Content-Type':'application/json'}, body);

httpClient.send(req, onComplete);

function onComplete(response, error) {
    if (response) {

        // get the response object from the exchange
        var responseObj = response.content.asJSON;

        if (responseObj.error) {
            throw new Error(responseObj.error_description);
        }

        // get the HTTP status code from the response
        if (response.status >= 300) {
            //Set the http response code
            // TODO: Raise a fault
            context.setVariable("error.status.code", response.status);
        }

        if(responseObj.result) {
            print("OPA allowed request");
        } else {
            print("OPA denied request");
            // This condition will raise a fault
            context.setVariable("triggerError", "true");
        }
    } else {
        throw new Error(error);
    }
}

function getRequestPath() {
    var path = context.getVariable("request.path");

    if (path != "/") {
        path = path.replace(/^\/|\/$/g, '');
        path = path.split("/");
    } else {
        path = ["/"];
    }
    return path;
}
