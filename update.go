package sq

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type setClause struct {
	column string
	value  interface{}
}

// UpdateBuilder builds SQL UPDATE statements.
type UpdateBuilder interface {
	// Prefix adds an expression to the beginning of the query.
	Prefix(sql string, args ...interface{}) UpdateBuilder

	// Table sets the table to be updated.
	Table(table string) UpdateBuilder

	// Set adds SET clauses to the query.
	Set(column string, value interface{}) UpdateBuilder

	// SetMap is a convenience method which calls .Set for each key/value pair in clauses.
	SetMap(clauses map[string]interface{}) UpdateBuilder

	// Where adds WHERE expressions to the query.
	//
	// See SelectBuilder.Where for more information.
	Where(pred interface{}, args ...interface{}) UpdateBuilder

	// OrderBy adds ORDER BY expressions to the query.
	OrderBy(orderBys ...string) UpdateBuilder

	// Limit sets a LIMIT clause on the query.
	Limit(limit uint64) UpdateBuilder

	// Offset sets a OFFSET clause on the query.
	Offset(offset uint64) UpdateBuilder

	// Suffix adds an expression to the end of the query.
	Suffix(sql string, args ...interface{}) UpdateBuilder

	ToSQL() (sqlStr string, args []interface{}, err error)
}

type updateBuilder struct {
	prefixes   exprs
	table      string
	setClauses []setClause
	whereParts []StatementBuilder
	orderBys   []string

	limit       uint64
	limitValid  bool
	offset      uint64
	offsetValid bool

	suffixes exprs
}

// NewUpdateBuilder creates new instance of UpdateBuilder.
func NewUpdateBuilder() UpdateBuilder {
	return &updateBuilder{}
}

func (b *updateBuilder) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(b.table) == 0 {
		err = fmt.Errorf("update statements must specify a table")
		return
	}
	if len(b.setClauses) == 0 {
		err = fmt.Errorf("update statements must have at least one Set clause")
		return
	}

	sql := &bytes.Buffer{}

	if len(b.prefixes) > 0 {
		args, _ = b.prefixes.AppendToSQL(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("UPDATE ")
	sql.WriteString(b.table)

	sql.WriteString(" SET ")
	setSQLs := make([]string, len(b.setClauses))
	for i, setClause := range b.setClauses {
		var valSQL string
		switch typedVal := setClause.value.(type) {
		case StatementBuilder:
			var valArgs []interface{}
			valSQL, valArgs, err = typedVal.ToSQL()
			if err != nil {
				return
			}
			args = append(args, valArgs...)
		default:
			valSQL = "?"
			args = append(args, typedVal)
		}
		setSQLs[i] = fmt.Sprintf("%s = %s", setClause.column, valSQL)
	}
	sql.WriteString(strings.Join(setSQLs, ", "))

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

	sqlStr, err = replacePlaceholders(sql.String())
	return
}

func (b *updateBuilder) Prefix(sql string, args ...interface{}) UpdateBuilder {
	b.prefixes = append(b.prefixes, expr{sql, args})
	return b
}

func (b *updateBuilder) Table(table string) UpdateBuilder {
	b.table = table
	return b
}

func (b *updateBuilder) Set(column string, value interface{}) UpdateBuilder {
	b.setClauses = append(b.setClauses, setClause{column: column, value: value})
	return b
}

func (b *updateBuilder) SetMap(clauses map[string]interface{}) UpdateBuilder {
	keys := make([]string, len(clauses))
	i := 0
	for key := range clauses {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		val := clauses[key]
		b.Set(key, val)
	}
	return b
}

func (b *updateBuilder) Where(pred interface{}, args ...interface{}) UpdateBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

func (b *updateBuilder) OrderBy(orderBys ...string) UpdateBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

func (b *updateBuilder) Limit(limit uint64) UpdateBuilder {
	b.limit = limit
	b.limitValid = true
	return b
}

func (b *updateBuilder) Offset(offset uint64) UpdateBuilder {
	b.offset = offset
	b.offsetValid = true
	return b
}

func (b *updateBuilder) Suffix(sql string, args ...interface{}) UpdateBuilder {
	b.suffixes = append(b.suffixes, expr{sql, args})

	return b
}
