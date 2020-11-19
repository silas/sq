package sq

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...string) *SelectBuilder {
	return NewSelectBuilder().Columns(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(into string) *InsertBuilder {
	return NewInsertBuilder().Into(into)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table string) *UpdateBuilder {
	return NewUpdateBuilder().Table(table)
}

// Delete returns a new DeleteBuilder for given table names.
//
// See DeleteBuilder.Table.
func Delete(what ...string) *DeleteBuilder {
	return NewDeleteBuilder().What(what...)
}

// Where returns a new WhereBuilder.
func Where(pred interface{}, args ...interface{}) *WhereBuilder {
	return NewWhereBuilder().Where(pred, args...)
}

// Case returns a new CaseBuilder
// "what" represents case value
func Case(what ...interface{}) *CaseBuilder {
	b := &CaseBuilder{}

	switch len(what) {
	case 0:
	case 1:
		b = b.what(what[0])
	default:
		b = b.what(newPart(what[0], what[1:]...))

	}
	return b
}
