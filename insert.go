package sq

import (
	"bytes"
	"fmt"
	"strings"
)

// InsertBuilder builds SQL INSERT statements.
type InsertBuilder interface {
	// Prefix adds an expression to the beginning of the query.
	Prefix(sql string, args ...interface{}) InsertBuilder

	// Options adds keyword options before the INTO clause of the query.
	Options(options ...string) InsertBuilder

	// Into sets the INTO clause of the query.
	Into(into string) InsertBuilder

	// Columns adds insert columns to the query.
	Columns(columns ...string) InsertBuilder

	// Values adds a single row's values to the query.
	Values(values ...interface{}) InsertBuilder

	// Suffix adds an expression to the end of the query.
	Suffix(sql string, args ...interface{}) InsertBuilder

	// SetMap set columns and values for insert builder from a map of column name and value
	// note that it will reset all previous columns and values was set if any.
	SetMap(clauses map[string]interface{}) InsertBuilder

	ToSQL() (sqlStr string, args []interface{}, err error)
}

type insertBuilder struct {
	prefixes exprs
	options  []string
	into     string
	columns  []string
	values   [][]interface{}
	suffixes exprs
}

// NewInsertBuilder creates new instance of InsertBuilder
func NewInsertBuilder() InsertBuilder {
	return &insertBuilder{}
}

func (b *insertBuilder) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(b.into) == 0 {
		err = fmt.Errorf("insert statements must specify a table")
		return
	}
	if len(b.values) == 0 {
		err = fmt.Errorf("insert statements must have at least one set of values")
		return
	}

	sql := &bytes.Buffer{}

	if len(b.prefixes) > 0 {
		args, _ = b.prefixes.AppendToSQL(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("INSERT ")

	if len(b.options) > 0 {
		sql.WriteString(strings.Join(b.options, " "))
		sql.WriteString(" ")
	}

	sql.WriteString("INTO ")
	sql.WriteString(b.into)
	sql.WriteString(" ")

	if len(b.columns) > 0 {
		sql.WriteString("(")
		sql.WriteString(strings.Join(b.columns, ","))
		sql.WriteString(") ")
	}

	sql.WriteString("VALUES ")

	valuesStrings := make([]string, len(b.values))
	for r, row := range b.values {

		valueStrings := make([]string, len(row))

		for v, val := range row {
			switch typedVal := val.(type) {
			case StatementBuilder:
				var valSQL string
				var valArgs []interface{}

				valSQL, valArgs, err = typedVal.ToSQL()
				if err != nil {
					return
				}

				valueStrings[v] = valSQL
				args = append(args, valArgs...)
			default:
				valueStrings[v] = "?"
				args = append(args, val)
			}
		}

		valuesStrings[r] = fmt.Sprintf("(%s)", strings.Join(valueStrings, ","))
	}
	sql.WriteString(strings.Join(valuesStrings, ","))

	if len(b.suffixes) > 0 {
		sql.WriteString(" ")
		args, _ = b.suffixes.AppendToSQL(sql, " ", args)
	}

	sqlStr = sql.String()
	return
}

func (b *insertBuilder) Prefix(sql string, args ...interface{}) InsertBuilder {
	b.prefixes = append(b.prefixes, expr{sql, args})
	return b
}

func (b *insertBuilder) Options(options ...string) InsertBuilder {
	b.options = append(b.options, options...)
	return b
}

func (b *insertBuilder) Into(into string) InsertBuilder {
	b.into = into
	return b
}

func (b *insertBuilder) Columns(columns ...string) InsertBuilder {
	b.columns = append(b.columns, columns...)
	return b
}

func (b *insertBuilder) Values(values ...interface{}) InsertBuilder {
	b.values = append(b.values, values)
	return b
}

func (b *insertBuilder) Suffix(sql string, args ...interface{}) InsertBuilder {
	b.suffixes = append(b.suffixes, expr{sql, args})
	return b
}

func (b *insertBuilder) SetMap(clauses map[string]interface{}) InsertBuilder {
	// TODO: replace resetting previous values with extending existing ones?
	cols := make([]string, 0, len(clauses))
	vals := make([]interface{}, 0, len(clauses))

	for col, val := range clauses {
		cols = append(cols, col)
		vals = append(vals, val)
	}

	b.columns = cols
	b.values = [][]interface{}{vals}

	return b
}
