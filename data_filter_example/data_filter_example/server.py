#!/usr/bin/env python

import sqlite3
import base64
import json

import flask
from flask_bootstrap import Bootstrap

from data_filter_example import opa

app = flask.Flask(__name__, static_url_path='/static')
Bootstrap(app)

def get_post(post_id):
    decision = query_opa("GET", ["posts", post_id])
    if not decision.defined:
        raise flask.abort(403)

    sql = opa.splice(SELECT='posts.*', FROM='posts', WHERE='posts.id=?', decision=decision)

    result = query_db(sql, args=(post_id,), one=True)
    if result is None:
        raise flask.abort(404)

    return result


def list_posts():
    decision = query_opa("GET", ["posts"])
    if not decision.defined:
        raise flask.abort(403)

    sql = opa.splice(SELECT='posts.*', FROM='posts', decision=decision)

    return query_db(sql)


def create_post(post):
    decision = query_opa("POST", ["posts"], post=post)
    if not decision.defined:
        raise flask.abort(403)

    elif decision.sql is not None:
        # POST API does not support conditional results.
        raise flask.abort(500)

    db = get_db()
    c = db.cursor()
    insert_post(c, post)
    db.commit()


def query_opa(method, path, **kwargs):
    input = {
        'method': method,
        'path': path,
        'subject': make_subject(),
    }
    for key in kwargs:
        input[key] = kwargs[key]
    return opa.compile(q='data.example.allow==true',
                       input=input,
                       unknowns=['posts'])


@app.route('/api/posts/<post_id>', methods=["GET"])
def api_get_post(post_id):
    return flask.jsonify(get_post(post_id))


@app.route('/api/posts', methods=["GET"])
def api_list_posts():
    return flask.jsonify(list_posts())


@app.route('/api/posts', methods=["POST"])
def api_create_post():
    post = flask.request.get_json(force=True)
    return flask.jsonify(create_post(post))


@app.route('/')
def index():
    kwargs = {
        'username': flask.request.cookies.get('user', ''),
        'posts': list_posts(),
    }
    if kwargs['username'] in USERS:
        kwargs['user'] = USERS[kwargs['username']]
    return flask.render_template('index.html', **kwargs)


@app.route('/login', methods=['POST'])
def login():
    user = flask.request.values.get('username')
    response = app.make_response(flask.redirect(flask.request.referrer))
    response.set_cookie('user', user)
    if user in USERS:
        for c in COOKIES:
            if c in USERS[user]:
                response.set_cookie(c, base64.b64encode(json.dumps(USERS[user][c])))
    return response


@app.route('/logout', methods=['GET'])
def logout():
    response = app.make_response(flask.redirect(flask.request.referrer))
    response.set_cookie('user', '', expires=0)
    for c in COOKIES:
        response.set_cookie(c, '', expires=0)
    return response


def make_subject():
    subject = {}
    user = flask.request.cookies.get('user', '')
    if user:
        subject['user'] = user
    for c in COOKIES:
        v = flask.request.cookies.get(c, '')
        if v:
            subject[c] = json.loads(base64.b64decode(v))
    return subject


def get_db():
    db = getattr(flask.g, '_database', None)
    if db is None:
        db = flask.g._database = sqlite3.connect('posts.db')
    db.row_factory = make_dicts
    return db


@app.teardown_appcontext
def close_connection(e):
    db = getattr(flask.g, '_database', None)
    if db is not None:
        db.close()


def init_schema():
    db = get_db()
    c = db.cursor()
    for table in TABLES:
        c.execute('DROP TABLE IF EXISTS ' + table['name'])
        c.execute(table['schema'])
    db.commit()


def pump_db():
    db = get_db()
    c = db.cursor()
    for table in TABLES:
        for row in table['data']:
            insert_object(table['name'], c, row)
    db.commit()


def init_db():
    with app.app_context():
        init_schema()
        pump_db()


def query_db(query, args=(), one=False):
    cur = get_db().execute(query, args)
    rv = cur.fetchall()
    cur.close()
    return (rv[0] if rv else None) if one else rv


def insert_post(cursor, post):
    insert_object('posts', cursor, post)


def insert_object(table, cursor, obj):
    row_keys = sorted(obj.keys())
    keys = '(' + ','.join(row_keys) + ')'
    values = '(' + ','.join(['?'] * len(row_keys)) + ')'
    stmt = 'INSERT INTO {} {} VALUES {}'.format(table, keys, values)
    args = [str(obj[k]) for k in row_keys]
    print(str(stmt), args)
    cursor.execute(stmt, args)


def make_dicts(cursor, row):
    return dict((cursor.description[idx][0], value) for idx, value in enumerate(row))


POSTS = [{
    'id': 'post1',
    'name': 'Personalization Updated to v2.1',
    'author': 'bob',
    'department': 'dev',
    'clearance_level': 3,
    'content': 'Leverage agile frameworks to provide a robust synopsis for high level overviews. Iterative approaches to corporate strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment.',
}, {
    'id': 'post2',
    'name': 'Critical Vulnerability in Y2K patch (CVE-2018-DEADBEEF)',
    'author': 'bob',
    'department': 'sec',
    'clearance_level': 2,
    'content': 'Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring.',
}, {
    'id': 'post3',
    'name': 'Blockchain Service Mesh deployed',
    'author': 'alice',
    'department': 'sec',
    'clearance_level': 5,
    'content': 'Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line.'
}, {
    'id': 'post4',
    'name': 'Quantum Gigaflux encountering errors',
    'author': 'alice',
    'department': 'sec',
    'clearance_level': 10,
    'content': 'Collaboratively administrate turnkey channels whereas virtual e-tailers. Objectively seize scalable metrics whereas proactive e-services. Seamlessly empower fully researched growth strategies and interoperable internal or "organic" sources.',
}, {
    'id': 'post5',
    'name': 'Missing printer',
    'author': 'charlie',
    'department': 'company',
    'clearance_level': 1,
    'content': 'Objectively innovate empowered manufactured products whereas parallel platforms. Holisticly predominate extensible testing procedures for reliable supply chains. Dramatically engage top-line web services vis-a-vis cutting-edge deliverables.',
}, {
    'id': 'post6',
    'name': 'Loud keyboards considered harmful',
    'author': 'charlie',
    'department': 'hr',
    'clearance_level': 10,
    'content': 'Podcasting operational change management inside of workflows to establish a framework. Taking seamless key performance indicators offline to maximise the long tail. Keeping your eye on the ball while performing a deep dive on the start-up mentality to derive convergence on cross-platform integration.',
}]

TABLES = [
    {
        'name': 'posts',
        'schema': """CREATE TABLE posts (
                        id TEXT PRIMARY KEY
                        , name TEXT
                        , author TEXT
                        , content TEXT
                        , clearance_level INT
                        , department TEXT)""",
        'data': POSTS,
    },
]

COOKIES = [
    'departments',
    'clearance_level',
]

USERS = {
    "bob": {
        "departments": ["sec", "dev"],
        "clearance_level": 5,
    },
    "alice": {
        "departments": ["dev", "sec"],
        "clearance_level": 10,
    },
    "charlie": {
        "departments": ["hr"],
        "clearance_level": 10,
    }
}

if __name__ == '__main__':
    init_db()
    app.jinja_env.auto_reload = True
    app.config['TEMPLATES_AUTO_RELOAD'] = True
    app.run(debug=True)
