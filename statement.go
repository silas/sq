package sq

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...string) SelectBuilder {
	return NewSelectBuilder().Columns(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(table string) InsertBuilder {
	return NewInsertBuilder().Into(table)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table string) UpdateBuilder {
	return NewUpdateBuilder().Table(table)
}

// Delete returns a new DeleteBuilder for given table names.
//
// See DeleteBuilder.From.
func Delete(table string) DeleteBuilder {
	return NewDeleteBuilder().From(table)
}

// Where returns a new WhereBuilder.
func Where(pred interface{}, args ...interface{}) WhereBuilder {
	return NewWhereBuilder().Where(pred, args...)
}

// Case returns a new CaseBuilder.
func Case(what ...interface{}) CaseBuilder {
	b := &caseBuilder{}

	switch len(what) {
	case 0:
		return b
	case 1:
		return b.what(what[0])
	default:
		return b.what(newPart(what[0], what[1:]...))
	}
}
