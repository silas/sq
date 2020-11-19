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
type WhereBuilder interface {
	Where(pred interface{}, args ...interface{}) WhereBuilder
	Select(columns ...string) SelectBuilder
	Update(table string) UpdateBuilder
	Delete(what ...string) DeleteBuilder
}

type whereBuilder struct {
	whereParts []QueryBuilder
}

// NewWhereBuilder creates new instance of UpdateBuilder
func NewWhereBuilder() WhereBuilder {
	return &whereBuilder{}
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b *whereBuilder) Where(pred interface{}, args ...interface{}) WhereBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

// Select returns a SelectBuilder for this WhereBuilder.
func (b *whereBuilder) Select(columns ...string) SelectBuilder {
	nb := NewSelectBuilder().Columns(columns...)
	nb.(*selectBuilder).whereParts = b.whereParts
	return nb
}

// Update returns a UpdateBuilder for this WhereBuilder.
func (b *whereBuilder) Update(table string) UpdateBuilder {
	nb := NewUpdateBuilder().Table(table)
	nb.(*updateBuilder).whereParts = b.whereParts
	return nb
}

// Delete returns a DeleteBuilder for this WhereBuilder.
func (b *whereBuilder) Delete(what ...string) DeleteBuilder {
	nb := NewDeleteBuilder().What(what...)
	nb.(*deleteBuilder).whereParts = b.whereParts
	return nb
}
