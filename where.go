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
	whereParts []QueryBuilder
}

// NewWhereBuilder creates new instance of UpdateBuilder
func NewWhereBuilder() *WhereBuilder {
	return &WhereBuilder{}
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
	nb := NewSelectBuilder().Columns(columns...)
	nb.whereParts = b.whereParts
	return nb
}

// Update returns a UpdateBuilder for this WhereBuilder.
func (b *WhereBuilder) Update(table string) *UpdateBuilder {
	nb := NewUpdateBuilder().Table(table)
	nb.whereParts = b.whereParts
	return nb
}

// Delete returns a DeleteBuilder for this WhereBuilder.
func (b *WhereBuilder) Delete(what ...string) *DeleteBuilder {
	nb := NewDeleteBuilder().What(what...)
	nb.whereParts = b.whereParts
	return nb
}
