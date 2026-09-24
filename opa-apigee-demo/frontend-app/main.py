from flask import Flask, render_template, g, redirect, url_for, request
from flask_oidc import OpenIDConnect
from okta import UsersClient
import requests

app = Flask(__name__)
app.config["OIDC_CLIENT_SECRETS"] = "client_secrets.json"
app.config["OIDC_COOKIE_SECURE"] = False
app.config["OIDC_CALLBACK_ROUTE"] = "/oidc/callback"
app.config["OIDC_SCOPES"] = ["openid", "email", "profile"]
app.config["SECRET_KEY"] = "secret"
app.config["OIDC_ID_TOKEN_COOKIE_NAME"] = "oidc_token"
oidc = OpenIDConnect(app)
okta_client = UsersClient("https://{{OKTA_DOMAIN}}", "{{OKTA_AUTH_TOKEN}}")

# user-id mapping
user_id = {"alice": "1001", "bob": "1002", "john": "1003", "jane": "1004"}

# mapping of the users the logged-in user is allowed to see
# this mapping is used to list the buttons on the dashboard
views = {"Alice Opa": ["alice", "bob", "john"], "Bob Rego": ["bob"], "John Opa": ["john"], "Jane Opa": ["jane"]}


@app.before_request
def before_request():
    if oidc.user_loggedin:
        g.user = okta_client.get_user(oidc.user_getfield("sub"))
    else:
        g.user = None


@app.route("/dev")
@app.route("/")
def index():
    return render_template("index.html")


@app.route("/dev/dashboard")
@app.route("/dashboard")
@oidc.require_login
def dashboard():
    info = oidc.user_getinfo(["name"])
    users = views[info['name']]
    return render_template("dashboard.html", users=users)


@app.route("/dev/login")
@app.route("/login")
@oidc.require_login
def login():
    return redirect(url_for(".dashboard"))


@app.route("/dev/logout")
@app.route("/logout")
def logout():
    oidc.logout()
    return redirect(url_for(".index"))


@app.route("/dev/profile")
@app.route("/profile")
@oidc.require_login
def profile():
    info = oidc.user_getinfo(["name", "email", "age", "enrolleeType"])
    return render_template("profile.html", profile=info)


@app.route('/dev/v1/claims/enrollee/<enrolleeName>', methods=['GET'])
@app.route('/v1/claims/enrollee/<enrolleeName>', methods=['GET'])
@oidc.require_login
def get_enrollee_claims(enrolleeName):

    # Apigee URL
    url = "https://{{APIGEE_URL}}/v1/claims/enrollee/{}".format(user_id[enrolleeName])

    # get user's ID token
    token = request.cookies.get('oidc_token')
    params = {'token': token}

    # dashboard users
    info = oidc.user_getinfo(["name"])
    users = views[info['name']]

    r = requests.get(url, params=params)

    if r.status_code == 200:
        claims = r.json().get("enrolleeClaimSummaryList")
        if claims is None or len(claims) == 0:
            return render_template("dashboard.html", users=users, status=403, error={"Result": "OPA denied request", "Status": "403 Forbidden"})
        else:
            return render_template("dashboard.html", users=users, status=r.status_code, claims=r.json()["enrolleeClaimSummaryList"])
    else:
        return render_template("dashboard.html", users=users, status=r.status_code, error=r.json())


if __name__ == "__main__":
    app.run()
