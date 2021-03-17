package sq

import (
	"errors"
	"strings"
)

// WithBuilder builds SQL WITH statements.
type WithBuilder interface {
	// With adds a new common table expression to the query.
	With(name string, field ...string) WithBuilder

	// Recursive adds RECURSIVE to the WITH clause.
	Recursive() WithBuilder

	// As sets the AS part of the common table expression.
	As(b StatementBuilder) WithBuilder

	// Union sets the UNION recursive subquery for the common table expression.
	Union(b StatementBuilder) WithBuilder

	// UnionAll sets the UNION ALL recursive subquery for the common table expression.
	UnionAll(b StatementBuilder) WithBuilder

	// Select returns the SelectBuilder for the parent query.
	Select(columns ...string) SelectBuilder
}

type withPart struct {
	name     string
	fields   []string
	as       StatementBuilder
	union    StatementBuilder
	unionAll bool
}

// NewWithBuilder creates new instance of WithBuilder.
func NewWithBuilder() WithBuilder {
	return &withBuilder{}
}

type withBuilder struct {
	withParts []*withPart
	recursive bool
	err       error
}

func (w *withBuilder) With(name string, field ...string) WithBuilder {
	w.withParts = append(w.withParts, &withPart{
		name:   name,
		fields: field,
	})
	return w
}

func (w *withBuilder) current() *withPart {
	if len(w.withParts) == 0 {
		w.err = errors.New("with statements must have WITH")
		return &withPart{}
	}
	return w.withParts[len(w.withParts)-1]
}

func (w *withBuilder) Recursive() WithBuilder {
	w.recursive = true
	return w
}

func (w *withBuilder) As(b StatementBuilder) WithBuilder {
	w.current().as = b
	return w
}

func (w *withBuilder) Union(b StatementBuilder) WithBuilder {
	c := w.current()
	c.union = b
	c.unionAll = false
	return w
}

func (w *withBuilder) UnionAll(b StatementBuilder) WithBuilder {
	c := w.current()
	c.union = b
	c.unionAll = true
	return w
}

func (w *withBuilder) Select(columns ...string) SelectBuilder {
	b := &selectBuilder{}
	b.Columns(columns...)

	if w.err != nil {
		b.prefixes = []expr{{err: w.err}}
		return b
	}

	if len(w.withParts) == 0 {
		return b
	}

	var sql strings.Builder
	var args []interface{}
	var sqlErr error

	sql.WriteString("WITH ")

	if w.recursive {
		sql.WriteString("RECURSIVE ")
	}

	for i, part := range w.withParts {
		if part.name == "" {
			sqlErr = errors.New("with statements must have a non-empty name")
			break
		}
		if part.as == nil {
			sqlErr = errors.New("with statements must have AS statement")
			break
		}

		if i > 0 {
			sql.WriteString(", ")
		}

		sql.WriteString(part.name)

		if len(part.fields) > 0 {
			sql.WriteString("(")
			for i, field := range part.fields {
				if i > 0 {
					sql.WriteString(", ")
				}
				sql.WriteString(field)
			}
			sql.WriteString(")")
		}

		sql.WriteString(" AS (")
		s, a, e := part.as.ToSQL()
		if e != nil {
			sqlErr = e
			break
		}
		sql.WriteString(s)
		args = append(args, a...)

		if part.union != nil {
			if part.unionAll {
				sql.WriteString(" UNION ALL ")
			} else {
				sql.WriteString(" UNION ")
			}

			s, a, e := part.union.ToSQL()
			if e != nil {
				sqlErr = e
				break
			}
			sql.WriteString(s)
			args = append(args, a...)
		}

		sql.WriteString(")")
	}

	if sqlErr != nil {
		b.prefixes = []expr{{err: sqlErr}}
	} else {
		b.prefixes = []expr{{sql: sql.String(), args: args}}
	}

	return b
}
