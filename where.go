package sq

import "fmt"

// WhereBuilder builds SQL where statements.
type WhereBuilder interface {
	// Where adds WHERE expressions to the query.
	//
	// See SelectBuilder.Where for more information.
	Where(pred interface{}, args ...interface{}) WhereBuilder

	// Select returns a SelectBuilder for this WhereBuilder.
	Select(columns ...string) SelectBuilder

	// Update returns a UpdateBuilder for this WhereBuilder.
	Update(table string) UpdateBuilder

	// Delete returns a DeleteBuilder for this WhereBuilder.
	Delete(table string) DeleteBuilder
}

type whereBuilder struct {
	whereParts []StatementBuilder
}

// NewWhereBuilder creates new instance of UpdateBuilder.
func NewWhereBuilder() WhereBuilder {
	return &whereBuilder{}
}

func (b *whereBuilder) Where(pred interface{}, args ...interface{}) WhereBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

func (b *whereBuilder) Select(columns ...string) SelectBuilder {
	nb := NewSelectBuilder().Columns(columns...)
	nb.(*selectBuilder).whereParts = b.whereParts
	return nb
}

func (b *whereBuilder) Update(table string) UpdateBuilder {
	nb := NewUpdateBuilder().Table(table)
	nb.(*updateBuilder).whereParts = b.whereParts
	return nb
}

func (b *whereBuilder) Delete(table string) DeleteBuilder {
	nb := NewDeleteBuilder().From(table)
	nb.(*deleteBuilder).whereParts = b.whereParts
	return nb
}

type wherePart part

func newWherePart(pred interface{}, args ...interface{}) StatementBuilder {
	return &wherePart{pred: pred, args: args}
}

func (p wherePart) ToSQL() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case StatementBuilder:
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