package sq

import "fmt"

type wherePart part

func newWherePart(pred interface{}, args ...interface{}) QueryBuilder {
	return &wherePart{pred: pred, args: args}
}

func (p wherePart) ToSQL() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case QueryBuilder:
		return pred.ToSQL()
	case map[string]interface{}:
		return Eq(pred).ToSQL()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string-keyed map or string, not %T", pred)
	}
	return
}

// WhereBuilder builds SQL where statements.
type WhereBuilder struct {
	StatementBuilderType

	whereParts []QueryBuilder
}

// NewWhereBuilder creates new instance of UpdateBuilder
func NewWhereBuilder(b StatementBuilderType) *WhereBuilder {
	return &WhereBuilder{StatementBuilderType: b}
}

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b *WhereBuilder) PlaceholderFormat(f PlaceholderFormat) *WhereBuilder {
	b.placeholderFormat = f
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b *WhereBuilder) Where(pred interface{}, args ...interface{}) *WhereBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

// Select returns a SelectBuilder for this WhereBuilder.
func (b *WhereBuilder) Select(columns ...string) *SelectBuilder {
	nb := NewSelectBuilder(b.StatementBuilderType).Columns(columns...)
	nb.whereParts = b.whereParts
	return nb
}

// Update returns a UpdateBuilder for this WhereBuilder.
func (b *WhereBuilder) Update(table string) *UpdateBuilder {
	nb := NewUpdateBuilder(b.StatementBuilderType).Table(table)
	nb.whereParts = b.whereParts
	return nb
}

// Delete returns a DeleteBuilder for this WhereBuilder.
func (b *WhereBuilder) Delete(what ...string) *DeleteBuilder {
	nb := NewDeleteBuilder(b.StatementBuilderType).What(what...)
	nb.whereParts = b.whereParts
	return nb
}
