/*
The code below will call out to OPA to filter claims
and pass them in the response sent to the client
*/

var opaUrl = "https://{{OPA_URL}}/v1/data/authz/response";
var method = context.getVariable("request.verb");
var path = getRequestPath();
var originalResponseObj = JSON.parse(context.getVariable("response.content"));
var token = context.getVariable("request.queryparam.token.1");

//DEBUG: Print original response
print(context.getVariable("response.content"));

// Input to OPA
var inputObj = {
    'method': method,
    'path': path,
    'object': originalResponseObj,
    'token': token
};

var body = JSON.stringify({input: inputObj});

var req = new Request(opaUrl, 'POST', {'Content-Type':'application/json'}, body);

httpClient.send(req, onComplete);

function onComplete(response, error) {
    if (response) {

        // get the response object
        var opaResponseObj = response.content.asJSON;

        if (opaResponseObj.error) {
            throw new Error(opaResponseObj.error_description);
        }

        // get the HTTP status code from the response
        if (response.status >= 300) {
            //Set the http response code
            print("OPA denied request");
            // This condition will raise a fault
            context.setVariable("triggerError", "true");
            context.setVariable("error.status.code", response.status);
        }

        // process filtered response
        if (opaResponseObj.result !== undefined) {
            // set the enrolleeClaimSummaryList
            originalResponseObj['enrolleeClaimSummaryList'] =  opaResponseObj.result.enrolleeClaimSummaryList;

            // update response returned to client
            context.setVariable("response.content", JSON.stringify(originalResponseObj));

            //DEBUG: Print updated response
            print(context.getVariable("response.content"));
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
