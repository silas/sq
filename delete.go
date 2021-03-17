package sq

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// DeleteBuilder builds SQL DELETE statements.
type DeleteBuilder interface {
	// Prefix adds an expression to the beginning of the query.
	Prefix(sql string, args ...interface{}) DeleteBuilder

	// From sets the FROM clause of the query.
	From(from string) DeleteBuilder

	// Where adds WHERE expressions to the query.
	Where(pred interface{}, args ...interface{}) DeleteBuilder

	// OrderBy adds ORDER BY expressions to the query.
	OrderBy(orderBys ...string) DeleteBuilder

	// Limit sets a LIMIT clause on the query.
	Limit(limit uint64) DeleteBuilder

	// Offset sets a OFFSET clause on the query.
	Offset(offset uint64) DeleteBuilder

	// Suffix adds an expression to the end of the query
	Suffix(sql string, args ...interface{}) DeleteBuilder

	// JoinClause adds a join clause to the query.
	JoinClause(join string) DeleteBuilder

	// Join adds a JOIN clause to the query.
	Join(join string) DeleteBuilder

	// LeftJoin adds a LEFT JOIN clause to the query.
	LeftJoin(join string) DeleteBuilder

	// RightJoin adds a RIGHT JOIN clause to the query.
	RightJoin(join string) DeleteBuilder

	ToSQL() (sqlStr string, args []interface{}, err error)
}

type deleteBuilder struct {
	prefixes   exprs
	from       string
	joins      []string
	whereParts []StatementBuilder
	orderBys   []string

	limit       uint64
	limitValid  bool
	offset      uint64
	offsetValid bool

	suffixes exprs
}

// NewDeleteBuilder creates new instance of DeleteBuilder
func NewDeleteBuilder() DeleteBuilder {
	return &deleteBuilder{}
}

func (b *deleteBuilder) ToSQL() (sqlStr string, args []interface{}, err error) {
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

	sqlStr = sql.String()
	return
}

func (b *deleteBuilder) Prefix(sql string, args ...interface{}) DeleteBuilder {
	b.prefixes = append(b.prefixes, expr{sql: sql, args: args})
	return b
}

func (b *deleteBuilder) From(from string) DeleteBuilder {
	b.from = from
	return b
}

func (b *deleteBuilder) Where(pred interface{}, args ...interface{}) DeleteBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

func (b *deleteBuilder) OrderBy(orderBys ...string) DeleteBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

func (b *deleteBuilder) Limit(limit uint64) DeleteBuilder {
	b.limit = limit
	b.limitValid = true
	return b
}

func (b *deleteBuilder) Offset(offset uint64) DeleteBuilder {
	b.offset = offset
	b.offsetValid = true

	return b
}

func (b *deleteBuilder) Suffix(sql string, args ...interface{}) DeleteBuilder {
	b.suffixes = append(b.suffixes, expr{sql: sql, args: args})

	return b
}

func (b *deleteBuilder) JoinClause(join string) DeleteBuilder {
	b.joins = append(b.joins, join)

	return b
}

func (b *deleteBuilder) Join(join string) DeleteBuilder {
	return b.JoinClause("JOIN " + join)
}

func (b *deleteBuilder) LeftJoin(join string) DeleteBuilder {
	return b.JoinClause("LEFT JOIN " + join)
}

func (b *deleteBuilder) RightJoin(join string) DeleteBuilder {
	return b.JoinClause("RIGHT JOIN " + join)
}
