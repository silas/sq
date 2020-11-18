package sq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertBuilderToSQL(t *testing.T) {
	b := Insert("").
		Prefix("WITH prefix AS ?", 0).
		Into("a").
		Options("DELAYED", "IGNORE").
		Columns("b", "c").
		Values(1, 2).
		Values(3, Expr("? + 1", 4)).
		Suffix("RETURNING ?", 5)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL :=
		"WITH prefix AS $1 " +
			"INSERT DELAYED IGNORE INTO a (b,c) VALUES ($2,$3),($4,$5 + 1) " +
			"RETURNING $6"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{0, 1, 2, 3, 4, 5}
	assert.Equal(t, expectedArgs, args)
}

func TestInsertBuilderToSQLErr(t *testing.T) {
	_, _, err := Insert("").Values(1).ToSQL()
	assert.Error(t, err)

	_, _, err = Insert("x").ToSQL()
	assert.Error(t, err)
}

func TestInsertBuilderPlaceholders(t *testing.T) {
	b := Insert("test").Values(1, 2)

	sql, _, _ := b.PlaceholderFormat(Question).ToSQL()
	assert.Equal(t, "INSERT INTO test VALUES (?,?)", sql)

	sql, _, _ = b.PlaceholderFormat(Dollar).ToSQL()
	assert.Equal(t, "INSERT INTO test VALUES ($1,$2)", sql)
}

func TestInsertBuilderSetMap(t *testing.T) {
	b := Insert("table").SetMap(Eq{"field1": 1})

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "INSERT INTO table (field1) VALUES ($1)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}
