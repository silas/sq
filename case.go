package sq

import (
	"bytes"
	"errors"
)

// CaseBuilder builds SQL CASE construct which could be used as parts of queries.
type CaseBuilder interface {
	When(when interface{}, then interface{}) CaseBuilder
	Else(expr interface{}) CaseBuilder
	ToSQL() (sqlStr string, args []interface{}, err error)
}

// queryBuilderBuffer is a helper that allows to write many QueryBuilders one by one
// without constant checks for errors that may come from QueryBuilder
type queryBuilderBuffer struct {
	bytes.Buffer
	args []interface{}
	err  error
}

// WriteSQL converts QueryBuilder to SQL strings and writes it to buffer
func (b *queryBuilderBuffer) WriteSQL(item QueryBuilder) {
	if b.err != nil {
		return
	}

	var str string
	var args []interface{}
	str, args, b.err = item.ToSQL()

	if b.err != nil {
		return
	}

	b.WriteString(str)
	b.WriteByte(' ')
	b.args = append(b.args, args...)
}

func (b *queryBuilderBuffer) ToSQL() (string, []interface{}, error) {
	return b.String(), b.args, b.err
}

// whenPart is a helper structure to describe SQLs "WHEN ... THEN ..." expression
type whenPart struct {
	when QueryBuilder
	then QueryBuilder
}

func newWhenPart(when interface{}, then interface{}) whenPart {
	return whenPart{newPart(when), newPart(then)}
}

type caseBuilder struct {
	whatPart  QueryBuilder
	whenParts []whenPart
	elsePart  QueryBuilder
}

// ToSQL implements QueryBuilder
func (b *caseBuilder) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(b.whenParts) == 0 {
		err = errors.New("case expression must contain at lease one WHEN clause")

		return
	}

	sql := queryBuilderBuffer{}

	sql.WriteString("CASE ")
	if b.whatPart != nil {
		sql.WriteSQL(b.whatPart)
	}

	for _, p := range b.whenParts {
		sql.WriteString("WHEN ")
		sql.WriteSQL(p.when)
		sql.WriteString("THEN ")
		sql.WriteSQL(p.then)
	}

	if b.elsePart != nil {
		sql.WriteString("ELSE ")
		sql.WriteSQL(b.elsePart)
	}

	sql.WriteString("END")

	return sql.ToSQL()
}

// what sets optional value for CASE construct "CASE [value] ..."
func (b *caseBuilder) what(expr interface{}) CaseBuilder {
	b.whatPart = newPart(expr)
	return b
}

// When adds "WHEN ... THEN ..." part to CASE construct
func (b *caseBuilder) When(when interface{}, then interface{}) CaseBuilder {
	// TODO: performance hint: replace slice of WhenPart with just slice of parts
	// where even indices of the slice belong to "when"s and odd indices belong to "then"s
	b.whenParts = append(b.whenParts, newWhenPart(when, then))
	return b
}

// Else sets optional "ELSE ..." part for CASE construct
func (b *caseBuilder) Else(expr interface{}) CaseBuilder {
	b.elsePart = newPart(expr)
	return b

}
