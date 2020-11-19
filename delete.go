package sq

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Builder

// DeleteBuilder builds SQL DELETE statements.
type DeleteBuilder struct {
	prefixes   exprs
	what       []string
	from       string
	joins      []string
	whereParts []QueryBuilder
	orderBys   []string

	limit       uint64
	limitValid  bool
	offset      uint64
	offsetValid bool

	suffixes exprs
}

// NewDeleteBuilder creates new instance of DeleteBuilder
func NewDeleteBuilder() *DeleteBuilder {
	return &DeleteBuilder{}
}

// ToSQL builds the query into a SQL string and bound args.
func (b *DeleteBuilder) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(b.from) == 0 {
		err = fmt.Errorf("delete statements must specify a From table")
		return
	}

	sql := &bytes.Buffer{}

	if len(b.prefixes) > 0 {
		args, _ = b.prefixes.AppendToSQL(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("DELETE ")
	// following condition helps to avoid duplicate "from" value in DELETE query
	// e.g. "DELETE a FROM a ..." which is valid for MySQL but not for PostgreSQL
	if len(b.what) > 0 && (len(b.what) != 1 || b.what[0] != b.from) {
		sql.WriteString(strings.Join(b.what, ", "))
		sql.WriteString(" ")
	}

	sql.WriteString("FROM ")
	sql.WriteString(b.from)

	if len(b.joins) > 0 {
		sql.WriteString(" ")
		sql.WriteString(strings.Join(b.joins, " "))
	}

	if len(b.whereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSQL(b.whereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(b.orderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.orderBys, ", "))
	}

	// TODO: limit == 0 and offset == 0 are valid. Need to go dbr way and implement offsetValid and limitValid
	if b.limitValid {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.FormatUint(b.limit, 10))
	}

	if b.offsetValid {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.FormatUint(b.offset, 10))
	}

	if len(b.suffixes) > 0 {
		sql.WriteString(" ")
		args, _ = b.suffixes.AppendToSQL(sql, " ", args)
	}

	sqlStr, err = ReplacePlaceholders(sql.String())
	return
}

// Prefix adds an expression to the beginning of the query
func (b *DeleteBuilder) Prefix(sql string, args ...interface{}) *DeleteBuilder {
	b.prefixes = append(b.prefixes, Expr(sql, args...))
	return b
}

// From sets the FROM clause of the query.
func (b *DeleteBuilder) From(from string) *DeleteBuilder {
	b.from = from
	return b
}

// What sets names of tables to be used for deleting from
func (b *DeleteBuilder) What(what ...string) *DeleteBuilder {
	filteredWhat := make([]string, 0, len(what))
	for _, item := range what {
		if len(item) > 0 {
			filteredWhat = append(filteredWhat, item)
		}
	}

	b.what = filteredWhat
	if len(filteredWhat) == 1 {
		b.From(filteredWhat[0])
	}

	return b
}

// Where adds WHERE expressions to the query.
func (b *DeleteBuilder) Where(pred interface{}, args ...interface{}) *DeleteBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b *DeleteBuilder) OrderBy(orderBys ...string) *DeleteBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

// Limit sets a LIMIT clause on the query.
func (b *DeleteBuilder) Limit(limit uint64) *DeleteBuilder {
	b.limit = limit
	b.limitValid = true
	return b
}

// Offset sets a OFFSET clause on the query.
func (b *DeleteBuilder) Offset(offset uint64) *DeleteBuilder {
	b.offset = offset
	b.offsetValid = true

	return b
}

// Suffix adds an expression to the end of the query
func (b *DeleteBuilder) Suffix(sql string, args ...interface{}) *DeleteBuilder {
	b.suffixes = append(b.suffixes, Expr(sql, args...))

	return b
}

// JoinClause adds a join clause to the query.
func (b *DeleteBuilder) JoinClause(join string) *DeleteBuilder {
	b.joins = append(b.joins, join)

	return b
}

// Join adds a JOIN clause to the query.
func (b *DeleteBuilder) Join(join string) *DeleteBuilder {
	return b.JoinClause("JOIN " + join)
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b *DeleteBuilder) LeftJoin(join string) *DeleteBuilder {
	return b.JoinClause("LEFT JOIN " + join)
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b *DeleteBuilder) RightJoin(join string) *DeleteBuilder {
	return b.JoinClause("RIGHT JOIN " + join)
}
