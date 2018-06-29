from data_filter_example import opa
import pytest
import requests


def put_policy(str):
    """Inserts a policy into OPA."""
    resp = requests.put('http://localhost:8181/v1/policies/test', data=str)
    resp.raise_for_status()


def clear_policies():
    """Deletes existing policies in OPA."""
    resp = requests.get("http://localhost:8181/v1/policies")
    resp.raise_for_status()
    body = resp.json()
    for policy in body["result"]:
        resp = requests.delete("http://localhost:8181/v1/policies/" + policy["id"])
        resp.raise_for_status()


one_table_assert_cases = [
    ('trivial', {
        "a": {
            "b": "foo"
        }
    }, '''package test

        p {
            data.q[x]
            x.b = input.a.b
        }''', True, '(("foo" = q.b))'),
    ('anonymous', {
        "a": {
            "c": "bar"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                    }''', False, None),
    [
        'inline', {
            "a": {
                "b": "foo"
            }
        }, '''package test
                    p {
                        data.q[_].b = input.a.b
                    }''', True, '(("foo" = q.b))'
    ],
    ('inline named var', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[i].b = input.a.b
                    }''', True, '(("foo" = q.b))'),
    ('assigned', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[_] = x
                        x.b = input.a.b
                    }''', True, '(("foo" = q.b))'),
    ('double eq', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[_] = x
                        x.b == input.a.b
                    }''', True, '(("foo" = q.b))'),
    ('conjunction', {
        "a": {
            "b": "foo",
            "c": "bar"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.c = input.a.c
                    }''', True, '(("foo" = q.b AND "bar" = q.c))'),
    ('disjunction data', {
        "a": {
            "b": "foo",
            "c": ["bar", "IT"]
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.c = input.a.c[_]
                    }''', True, '(("foo" = q.b AND "bar" = q.c) OR ("foo" = q.b AND "IT" = q.c))'),
    ('disjunction rules', {
        "a": {
            "b": "foo",
            "c": "bar"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                    }
                    p {
                        data.q[x]
                        x.c = input.a.c
                    }''', True, '(("foo" = q.b) OR ("bar" = q.c))'),
    ('undefined context', {
        "a": {
            "b": "foo",
            "c": "bar"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                    }
                    p {
                        data.r[x]  # data.r is undefined so this rule will not contribute to the result.
                        x.b = input.a.b
                    }
                    ''', True, '(("foo" = q.b))'),
    ('neq', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.exclude != true
                    }
                    ''', True, '(("foo" = q.b AND q.exclude != true))'),
    ('lt', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.n < 1
                    }
                    ''', True, '(("foo" = q.b AND q.n < 1))'),
    ('lte', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.n <= 1
                    }
                    ''', True, '(("foo" = q.b AND q.n <= 1))'),
    ('gt', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.n > 1
                    }
                    ''', True, '(("foo" = q.b AND q.n > 1))'),
    ('gte', {
        "a": {
            "b": "foo"
        }
    }, '''package test
                    p {
                        data.q[x]
                        x.b = input.a.b
                        x.n >= 1
                    }
                    ''', True, '(("foo" = q.b AND q.n >= 1))'),
    (
        'nested',
        {
            "a": 1
        },
        '''package test
                    p {
                        data.q[x]
                        abs(x.a) > input.a
                    }''',
        True,
        '((abs(q.a) > 1))',
    ),
    (
        'nested conjunction',
        {
            "a": 1
        },
        '''package test
                    p {
                        data.q[x]
                        x.b = 1
                        abs(x.a) > input.a
                    }''',
        True,
        '((q.b = 1 AND abs(q.a) > 1))',
    ),
    (
        'nested conjunction inline',
        {
            "a": 1
        },
        '''package test
                    p {
                        data.q[i].b = 1
                        abs(data.q[i].a) > input.a
                    }''',
        True,
        '((q.b = 1 AND abs(q.a) > 1))',
    ),
    (
        'intermediate vars',
        {},
        '''package test
        p {
            data.q = x
            x[i] = y
            y = z
            z.a = 1
            y.b = 2
        }''',
        True,
        '((q.a = 1 AND q.b = 2))',
    ),
    ('set based', {}, '''package test
        p {
            p1[x]
        }

        p1[x] {
            data.q[x].a = 1
        }

        p1[y] {
            data.q[y].b = 2
        }''', True, '((q.a = 1) OR (q.b = 2))'),
    (
        'unsupported built-in function',
        {},
        '''package test
        p {
            count(data.q[_].a) > input
        }''',
        opa.TranslationError("operator not supported: count"),
        None,
    ),
    (
        'non-relation expression',
        {},
        '''package test
        p {
            plus(data.q[_].a, 10, 10)
        }''',
        opa.TranslationError('too many arguments'),
        None,
    ),
    (
        'invalid row identifier',
        {},
        '''package test
        p {
            data.q.foo.bar = 10
        }''',
        opa.TranslationError('row identifier type not supported'),
        None,
    ),
]

multi_table_assert_cases = [
    (
        'simple join',
        {},
        '''package test
        p {
            data.q[x].a = data.r[y].b
        }''',
        True,
        [[['r'], '(q.a = r.b)']],
    ),
    (
        'three-way join',
        {},
        '''package test
        p {
            data.q[x].a = data.r[y].b
            data.q[x].c = data.s[z].c
        }''',
        True,
        [[['s', 'r'], '(q.a = r.b AND q.c = s.c)']],
    ),
    (
        'mixed',
        {},
        '''package test
        p {
            data.q[x].a = 10
        }
        p {
            data.q[y].a = data.r[z].b
        }''',
        True,
        ['((q.a = 10))', [['r'], '(q.a = r.b)']],
    ),
    (
        'self-join',
        {},
        '''package test
        p {
            data.q[_].a = 10
            data.q[_].b = 20
        }''',
        opa.TranslationError('self-joins not supported'),
        [],
    ),
]


@pytest.mark.parametrize(
    'note,input,policy,exp_defined,exp_sql',
    one_table_assert_cases,
)
def test_compile_one_table(note, input, policy, exp_defined, exp_sql):
    crunch('data.test.p = true', input, ['q'], 'q', policy, exp_defined, ['WHERE ' + exp_sql]
           if exp_sql is not None else None)


@pytest.mark.parametrize('note,input,policy,exp_defined,exp_sql', one_table_assert_cases)
def test_compile_one_table_double_eq(note, input, policy, exp_defined, exp_sql):
    crunch('data.test.p == true', input, ['q'], 'q', policy, exp_defined, ['WHERE ' + exp_sql]
           if exp_sql is not None else None)


@pytest.mark.parametrize('note,input,policy,exp_defined,exp_sql', multi_table_assert_cases)
def test_compile_multi_table(note, input, policy, exp_defined, exp_sql):
    clauses = []
    for clause in exp_sql:
        if isinstance(clause, str):
            clauses.append('WHERE ' + clause)
        else:
            joins = ' '.join('INNER JOIN ' + t for t in clause[0])
            clauses.append(joins + ' ON ' + clause[1])
    crunch(
        'data.test.p = true',
        input,
        ['q', 'r', 's'],
        'q',
        policy,
        exp_defined,
        clauses,
    )


def crunch(query, input, unknowns, from_table, policy, exp_defined, exp_sql):
    clear_policies()
    put_policy(policy)
    try:
        result = opa.compile(query, input, unknowns, from_table)
    except opa.TranslationError as e:
        if not isinstance(exp_defined, opa.TranslationError):
            raise
        assert str(exp_defined) in str(e)
    else:
        assert result.defined == exp_defined
        if result.defined:
            if exp_sql is None:
                assert result.sql is None
            else:
                assert [c.sql() for c in result.sql.clauses] == exp_sql
        else:
            assert result.sql is None
