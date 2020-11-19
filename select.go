package sq

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// SelectBuilder builds SQL SELECT statements.
type SelectBuilder interface {
	Prefix(sql string, args ...interface{}) SelectBuilder
	Distinct() SelectBuilder
	Columns(columns ...string) SelectBuilder
	Column(column interface{}, args ...interface{}) SelectBuilder
	From(from string) SelectBuilder
	JoinClause(join string) SelectBuilder
	Join(join string) SelectBuilder
	LeftJoin(join string) SelectBuilder
	RightJoin(join string) SelectBuilder
	Where(pred interface{}, args ...interface{}) SelectBuilder
	GroupBy(groupBys ...string) SelectBuilder
	Having(pred interface{}, rest ...interface{}) SelectBuilder
	OrderBy(orderBys ...string) SelectBuilder
	Limit(limit uint64) SelectBuilder
	Offset(offset uint64) SelectBuilder
	Suffix(sql string, args ...interface{}) SelectBuilder
	ToSQL() (sqlStr string, args []interface{}, err error)
}

type selectBuilder struct {
	prefixes    exprs
	distinct    bool
	columns     []QueryBuilder
	from        string
	joins       []string
	whereParts  []QueryBuilder
	groupBys    []string
	havingParts []QueryBuilder
	orderBys    []string

	limit       uint64
	limitValid  bool
	offset      uint64
	offsetValid bool

	suffixes exprs
}

// NewSelectBuilder creates new instance of SelectBuilder
func NewSelectBuilder() SelectBuilder {
	return &selectBuilder{}
}

// ToSQL builds the query into a SQL string and bound args.
func (b *selectBuilder) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(b.columns) == 0 {
		err = fmt.Errorf("select statements must have at least one result column")
		return
	}

	sql := &bytes.Buffer{}

	if len(b.prefixes) > 0 {
		args, _ = b.prefixes.AppendToSQL(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("SELECT ")

	if b.distinct {
		sql.WriteString("DISTINCT ")
	}

	if len(b.columns) > 0 {
		args, err = appendToSQL(b.columns, sql, ", ", args)
		if err != nil {
			return
		}
	}

	if len(b.from) > 0 {
		sql.WriteString(" FROM ")
		sql.WriteString(b.from)
	}

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

	if len(b.groupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(b.groupBys, ", "))
	}

	if len(b.havingParts) > 0 {
		sql.WriteString(" HAVING ")
		args, err = appendToSQL(b.havingParts, sql, " AND ", args)
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
func (b *selectBuilder) Prefix(sql string, args ...interface{}) SelectBuilder {
	b.prefixes = append(b.prefixes, Expr(sql, args...))
	return b
}

// Distinct adds a DISTINCT clause to the query.
func (b *selectBuilder) Distinct() SelectBuilder {
	b.distinct = true

	return b
}

// Columns adds result columns to the query.
func (b *selectBuilder) Columns(columns ...string) SelectBuilder {
	for _, str := range columns {
		b.columns = append(b.columns, newPart(str))
	}

	return b
}

// Column adds a result column to the query.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//   Column("IF(col IN ("+Placeholders(3)+"), 1, 0) as col", 1, 2, 3)
func (b *selectBuilder) Column(column interface{}, args ...interface{}) SelectBuilder {
	b.columns = append(b.columns, newPart(column, args...))

	return b
}

// From sets the FROM clause of the query.
func (b *selectBuilder) From(from string) SelectBuilder {
	b.from = from
	return b
}

// JoinClause adds a join clause to the query.
func (b *selectBuilder) JoinClause(join string) SelectBuilder {
	b.joins = append(b.joins, join)

	return b
}

// Join adds a JOIN clause to the query.
func (b *selectBuilder) Join(join string) SelectBuilder {
	return b.JoinClause("JOIN " + join)
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b *selectBuilder) LeftJoin(join string) SelectBuilder {
	return b.JoinClause("LEFT JOIN " + join)
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b *selectBuilder) RightJoin(join string) SelectBuilder {
	return b.JoinClause("RIGHT JOIN " + join)
}

// Where adds an expression to the WHERE clause of the query.
//
// Expressions are ANDed together in the generated SQL.
//
// Where accepts several types for its pred argument:
//
// nil OR "" - ignored.
//
// string - SQL expression.
// If the expression has SQL placeholders then a set of arguments must be passed
// as well, one for each placeholder.
//
// map[string]interface{} OR Eq - map of SQL expressions to values. Each key is
// transformed into an expression like "<key> = ?", with the corresponding value
// bound to the placeholder. If the value is nil, the expression will be "<key>
// IS NULL". If the value is an array or slice, the expression will be "<key> IN
// (?,?,...)", with one placeholder for each item in the value. These expressions
// are ANDed together.
//
// Where will panic if pred isn't any of the above types.
func (b *selectBuilder) Where(pred interface{}, args ...interface{}) SelectBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

// GroupBy adds GROUP BY expressions to the query.
func (b *selectBuilder) GroupBy(groupBys ...string) SelectBuilder {
	b.groupBys = append(b.groupBys, groupBys...)
	return b
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b *selectBuilder) Having(pred interface{}, rest ...interface{}) SelectBuilder {
	b.havingParts = append(b.havingParts, newWherePart(pred, rest...))
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b *selectBuilder) OrderBy(orderBys ...string) SelectBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

// Limit sets a LIMIT clause on the query.
func (b *selectBuilder) Limit(limit uint64) SelectBuilder {
	b.limit = limit
	b.limitValid = true
	return b
}

// Offset sets a OFFSET clause on the query.
func (b *selectBuilder) Offset(offset uint64) SelectBuilder {
	b.offset = offset
	b.offsetValid = true
	return b
}

// Suffix adds an expression to the end of the query
func (b *selectBuilder) Suffix(sql string, args ...interface{}) SelectBuilder {
	b.suffixes = append(b.suffixes, Expr(sql, args...))

	return b
}
