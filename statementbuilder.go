package sq

// StatementBuilder is the interface that wraps the ToSQL method.
//
// ToSQL returns a SQL representation of the StatementBuilder, along with a slice of args
// as passed to e.g. database/sql.Exec. It can also return an error.
type StatementBuilder interface {
	ToSQL() (string, []interface{}, error)
}
