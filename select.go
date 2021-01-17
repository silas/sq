package sq

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// SelectBuilder builds SQL SELECT statements.
type SelectBuilder interface {
	// Prefix adds an expression to the beginning of the query.
	Prefix(sql string, args ...interface{}) SelectBuilder

	// Distinct adds a DISTINCT clause to the query.
	Distinct() SelectBuilder

	// Columns adds result columns to the query.
	Columns(columns ...string) SelectBuilder

	// Column adds a result column to the query.
	// Unlike Columns, Column accepts args which will be bound to placeholders in
	// the columns string.
	//
	//   Column("IF(col IN ("+placeholders(3)+"), 1, 0) as col", 1, 2, 3)
	Column(column interface{}, args ...interface{}) SelectBuilder

	// From sets the FROM clause of the query.
	From(from string) SelectBuilder

	// JoinClause adds a join clause to the query.
	JoinClause(join string, args ...interface{}) SelectBuilder

	// Join adds a JOIN clause to the query.
	Join(join string, args ...interface{}) SelectBuilder

	// LeftJoin adds a LEFT JOIN clause to the query.
	LeftJoin(join string, args ...interface{}) SelectBuilder

	// RightJoin adds a RIGHT JOIN clause to the query.
	RightJoin(join string, args ...interface{}) SelectBuilder

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
	Where(pred interface{}, args ...interface{}) SelectBuilder

	// GroupBy adds GROUP BY expressions to the query.
	GroupBy(groupBys ...string) SelectBuilder

	// Having adds an expression to the HAVING clause of the query.
	//
	// See Where.
	Having(pred interface{}, rest ...interface{}) SelectBuilder

	// OrderBy adds ORDER BY expressions to the query.
	OrderBy(orderBys ...string) SelectBuilder

	// Limit sets a LIMIT clause on the query.
	Limit(limit uint64) SelectBuilder

	// Offset sets a OFFSET clause on the query.
	Offset(offset uint64) SelectBuilder

	// Suffix adds an expression to the end of the query.
	Suffix(sql string, args ...interface{}) SelectBuilder

	ToSQL() (sqlStr string, args []interface{}, err error)
}

type selectBuilder struct {
	prefixes    exprs
	distinct    bool
	columns     []StatementBuilder
	from        string
	joins       exprs
	whereParts  []StatementBuilder
	groupBys    []string
	havingParts []StatementBuilder
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
		args, err = b.joins.AppendToSQL(sql, " ", args)
		if err != nil {
			return
		}
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

	sqlStr = sql.String()
	return

}

func (b *selectBuilder) Prefix(sql string, args ...interface{}) SelectBuilder {
	b.prefixes = append(b.prefixes, expr{sql, args})
	return b
}

func (b *selectBuilder) Distinct() SelectBuilder {
	b.distinct = true

	return b
}

func (b *selectBuilder) Columns(columns ...string) SelectBuilder {
	for _, str := range columns {
		b.columns = append(b.columns, newPart(str))
	}

	return b
}

func (b *selectBuilder) Column(column interface{}, args ...interface{}) SelectBuilder {
	b.columns = append(b.columns, newPart(column, args...))

	return b
}

func (b *selectBuilder) From(from string) SelectBuilder {
	b.from = from
	return b
}

func (b *selectBuilder) JoinClause(join string, args ...interface{}) SelectBuilder {
	b.joins = append(b.joins, expr{join, args})

	return b
}

func (b *selectBuilder) Join(join string, args ...interface{}) SelectBuilder {
	return b.JoinClause("JOIN "+join, args...)
}

func (b *selectBuilder) LeftJoin(join string, args ...interface{}) SelectBuilder {
	return b.JoinClause("LEFT JOIN "+join, args...)
}

func (b *selectBuilder) RightJoin(join string, args ...interface{}) SelectBuilder {
	return b.JoinClause("RIGHT JOIN "+join, args...)
}

func (b *selectBuilder) Where(pred interface{}, args ...interface{}) SelectBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

func (b *selectBuilder) GroupBy(groupBys ...string) SelectBuilder {
	b.groupBys = append(b.groupBys, groupBys...)
	return b
}

func (b *selectBuilder) Having(pred interface{}, rest ...interface{}) SelectBuilder {
	b.havingParts = append(b.havingParts, newWherePart(pred, rest...))
	return b
}

func (b *selectBuilder) OrderBy(orderBys ...string) SelectBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

func (b *selectBuilder) Limit(limit uint64) SelectBuilder {
	b.limit = limit
	b.limitValid = true
	return b
}

func (b *selectBuilder) Offset(offset uint64) SelectBuilder {
	b.offset = offset
	b.offsetValid = true
	return b
}

func (b *selectBuilder) Suffix(sql string, args ...interface{}) SelectBuilder {
	b.suffixes = append(b.suffixes, expr{sql, args})

	return b
}
