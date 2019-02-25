import json


class TranslateSettings:
    def __init__(self, quoteType='"'):
        self.quoteType = quoteType

    def __call__(self, quoteType='"'):
        self.quoteType = quoteType


class Union(object):
    def __init__(self, clauses, transSet: TranslateSettings):
        self.clauses = clauses


class InnerJoin(object):
    def __init__(self, tables, expr, transSet: TranslateSettings):
        self.tables = tables
        self.expr = expr
        self.settings = transSet

    def sql(self):
        return (
            " ".join(["INNER JOIN " + t for t in self.tables])
            + " ON "
            + self.expr.sql(transSet)
        )


class Where(object):
    def __init__(self, expr, transSet: TranslateSettings):
        self.expr = expr
        self.settings = transSet

    def sql(self):
        return "WHERE " + self.expr.sql(transSet)


class Disjunction(object):
    def __init__(self, conjunction, transSet: TranslateSettings):
        self.conjunction = conjunction
        self.settings = transSet

    def sql(self):
        return "(" + " OR ".join([c.sql(transSet) for c in self.conjunction]) + ")"


class Conjunction(object):
    def __init__(self, relation, transSet: TranslateSettings):
        self.relation = relation
        self.settings = transSet

    def sql(self):
        if len(self.relation) == 0:
            return "1"
        return "(" + " AND ".join([r.sql(transSet) for r in self.relation]) + ")"


class Relation(object):
    def __init__(self, operator, lhs, rhs, transSet: TranslateSettings):
        self.operator = operator
        self.lhs = lhs
        self.rhs = rhs
        self.settings = transSet

    def sql(self):
        return "%s %s %s" % (
            self.lhs.sql(transSet),
            self.operator.sql(transSet),
            self.rhs.sql(transSet),
        )


class Column(object):
    def __init__(self, name, transSet: TranslateSettings, table=""):
        self.table = table
        self.name = name
        self.settings = transSet

    def sql(self):
        if self.table:
            return "%s.%s" % (self.table, self.name)
        return str(self.name)


class Call(object):
    def __init__(self, operator, operands, transSet: TranslateSettings):
        self.operator = operator
        self.operands = operands
        self.settings = transSet

    def sql(self):
        return (
            self.operator
            + "("
            + ", ".join(o.sql(transSet) for o in self.operands)
            + ")"
        )


class Constant(object):
    def __init__(self, value, transSet: TranslateSettings):
        self.value = value
        self.settings = transSet

    def sql(self):
        tr = json.dumps(self.value)
        if tr[0] == '"' and tr[-1] == '"':
            tr[0] = transSet.quoteType
            tr[-1] = transSet.quoteType
        return tr


class RelationOp(object):
    def __init__(self, value, transSet: TranslateSettings):
        self.value = value
        self.settings = transSet

    def sql(self):
        return self.value


def walk(node, vis, transSet: TranslateSettings):
    next = vis(node)
    if next is None:
        return
    if isinstance(node, Union):
        for c in node.clauses:
            walk(c, next, transSet)
    elif isinstance(node, Where):
        walk(node.expr, next, transSet)
    elif isinstance(node, InnerJoin):
        walk(node.expr, next, transSet)
    elif isinstance(node, Disjunction):
        for child in node.conjunction:
            walk(child, next, transSet)
    elif isinstance(node, Conjunction):
        for child in node.relation:
            walk(child, next, transSet)
    elif isinstance(node, Relation):
        walk(node.operator, next, transSet)
        walk(node.lhs, next, transSet)
        walk(node.rhs, next, transSet)
    elif isinstance(node, Call):
        walk(node.operator, next)
        for o in node.operands:
            walk(o, next, transSet)


def pretty_print(node, transSet: TranslateSettings):
    class printer(object):
        def __init__(self, indent):
            self.indent = indent

        def __call__(self, node):
            print(" " * self.indent, node.__class__.__name__)
            return printer(self.indent + 2)

    vis = printer(0)
    walk(node, vis, transSet)
