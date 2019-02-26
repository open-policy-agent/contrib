"""@todo."""

import json


class TranslationSettings(object):
    """@todo."""

    def __init__(self, quoteType='"'):
        """@todo."""
        self.quoteType = quoteType

    def __call__(self, quoteType='"'):
        """@todo."""
        self.quoteType = quoteType


class Union(object):
    """@todo."""

    def __init__(self, clauses, transSet: TranslationSettings):
        """@todo."""
        self.clauses = clauses


class InnerJoin(object):
    """@todo."""

    def __init__(self, tables, expr, transSet: TranslationSettings):
        """@todo."""
        self.tables = tables
        self.expr = expr
        self.settings = transSet

    def sql(self):
        """@todo."""
        return (
            " ".join(["INNER JOIN " + t for t in self.tables])
            + " ON "
            + self.expr.sql()
        )


class Where(object):
    """@todo."""

    def __init__(self, expr, transSet: TranslationSettings):
        """@todo."""
        self.expr = expr
        self.settings = transSet

    def sql(self):
        """@todo."""
        return "WHERE " + self.expr.sql()


class Disjunction(object):
    """@todo."""

    def __init__(self, conjunction, transSet: TranslationSettings):
        """@todo."""
        self.conjunction = conjunction
        self.settings = transSet

    def sql(self):
        """@todo."""
        return "(" + " OR ".join([c.sql() for c in self.conjunction]) + ")"


class Conjunction(object):
    """@todo."""

    def __init__(self, relation, transSet: TranslationSettings):
        """@todo."""
        self.relation = relation
        self.settings = transSet

    def sql(self):
        """@todo."""
        if len(self.relation) == 0:
            return "1"
        return "(" + " AND ".join([r.sql() for r in self.relation]) + ")"


class Relation(object):
    """@todo."""

    def __init__(self, operator, lhs, rhs, transSet: TranslationSettings):
        """@todo."""
        self.operator = operator
        self.lhs = lhs
        self.rhs = rhs
        self.settings = transSet

    def sql(self):
        """@todo."""
        return "%s %s %s" % (self.lhs.sql(), self.operator.sql(), self.rhs.sql())


class Column(object):
    """@todo."""

    def __init__(self, name, transSet: TranslationSettings, table=""):
        """@todo."""
        self.table = table
        self.name = name
        self.settings = transSet

    def sql(self):
        """@todo."""
        if self.table:
            return "%s.%s" % (self.table, self.name)
        return str(self.name)


class Call(object):
    """@todo."""

    def __init__(self, operator, operands, transSet: TranslationSettings):
        """@todo."""
        self.operator = operator
        self.operands = operands
        self.settings = transSet

    def sql(self):
        """@todo."""
        return self.operator + "(" + ", ".join(o.sql() for o in self.operands) + ")"


class Constant(object):
    """@todo."""

    def __init__(self, value, transSet: TranslationSettings):
        """@todo."""
        self.value = value
        self.settings = transSet

    def sql(self):
        """@todo."""
        tr = list(json.dumps(self.value))

        if tr[0] == '"' and tr[-1] == '"':
            tr[0] = self.settings.quoteType
            tr[-1] = self.settings.quoteType
        return str(tr)


class RelationOp(object):
    """@todo."""

    def __init__(self, value, transSet: TranslationSettings):
        """@todo."""
        self.value = value
        self.settings = transSet

    def sql(self):
        """@todo."""
        return self.value


def walk(node, vis, transSet: TranslationSettings):
    """@todo."""
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


def pretty_print(node, transSet: TranslationSettings):
    """@todo."""

    class printer(object):
        """@todo."""

        def __init__(self, indent):
            """@todo."""
            self.indent = indent

        def __call__(self, node):
            """@todo."""
            print(" " * self.indent, node.__class__.__name__)
            return printer(self.indent + 2)

    vis = printer(0)
    walk(node, vis, transSet)
