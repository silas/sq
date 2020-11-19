package sq

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqToSQL(t *testing.T) {
	b := Eq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id = ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestEqInToSQL(t *testing.T) {
	b := Eq{"id": []int{1, 2, 3}}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id IN (?,?,?)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1, 2, 3}
	assert.Equal(t, expectedArgs, args)
}

func TestNotEqToSQL(t *testing.T) {
	b := NotEq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id <> ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestEqNotInToSQL(t *testing.T) {
	b := NotEq{"id": []int{1, 2, 3}}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id NOT IN (?,?,?)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1, 2, 3}
	assert.Equal(t, expectedArgs, args)
}

func TestExprNilToSQL(t *testing.T) {
	var b StatementBuilder
	b = NotEq{"name": nil}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)

	expectedSQL := "name IS NOT NULL"
	assert.Equal(t, expectedSQL, sql)

	b = Eq{"name": nil}
	sql, args, err = b.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)

	expectedSQL = "name IS NULL"
	assert.Equal(t, expectedSQL, sql)
}

func TestLtToSQL(t *testing.T) {
	b := Lt{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id < ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestLtExprToSQL(t *testing.T) {
	b := Lt{"time": Expr("current_timestamp + interval '? minutes'", 10)}
	sql, args, err := b.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, args, []interface{}{10})
	assert.Equal(t, "time < current_timestamp + interval '? minutes'", sql)
}

func TestLtOrEqToSQL(t *testing.T) {
	b := LtOrEq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id <= ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestGtToSQL(t *testing.T) {
	b := Gt{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id > ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestGtOrEqToSQL(t *testing.T) {
	b := GtOrEq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id >= ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestNullTypeString(t *testing.T) {
	var b StatementBuilder
	var name sql.NullString

	b = Eq{"name": name}
	sql, args, err := b.ToSQL()

	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "name IS NULL", sql)

	assert.NoError(t, name.Scan("Name"))
	b = Eq{"name": name}
	sql, args, err = b.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"Name"}, args)
	assert.Equal(t, "name = ?", sql)
}

func TestNullTypeInt64(t *testing.T) {
	var userID sql.NullInt64
	assert.NoError(t, userID.Scan(nil))
	b := Eq{"user_id": userID}
	sql, args, err := b.ToSQL()

	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "user_id IS NULL", sql)

	assert.NoError(t, userID.Scan(10))
	b = Eq{"user_id": userID}
	sql, args, err = b.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, []interface{}{int64(10)}, args)
	assert.Equal(t, "user_id = ?", sql)
}

type dummyQueryBuilder int

func (d dummyQueryBuilder) ToSQL() (string, []interface{}, error) {
	return "DUMMY(?, ?)", []interface{}{int(d), int(d)}, nil
}

func TestExprQueryBuilder(t *testing.T) {
	b := Expr("EXISTS(?)", dummyQueryBuilder(42))
	sql, args, err := b.ToSQL()

	if assert.NoError(t, err) {
		assert.Equal(t, "EXISTS(DUMMY(?, ?))", sql)
		assert.Equal(t, []interface{}{42, 42}, args)
	}
}
