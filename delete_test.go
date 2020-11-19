package sq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteBuilderToSQL(t *testing.T) {
	b := Delete("").
		Prefix("WITH prefix AS ?", 0).
		From("a").
		Where("b = ?", 1).
		OrderBy("c").
		Limit(2).
		Offset(3).
		Suffix("RETURNING ?", 4)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL :=
		"WITH prefix AS $1 " +
			"DELETE FROM a WHERE b = $2 ORDER BY c LIMIT 2 OFFSET 3 " +
			"RETURNING $3"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{0, 1, 4}
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteWithoutFrom(t *testing.T) {
	b := Delete("a").
		Where("b = ?", 1)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "DELETE FROM a WHERE b = $1"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteBuilderZeroOffsetLimit(t *testing.T) {
	qb := Delete("").
		From("b").
		Limit(0).
		Offset(0)

	sql, _, err := qb.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "DELETE FROM b LIMIT 0 OFFSET 0"
	assert.Equal(t, expectedSQL, sql)
}

func TestDeleteBuilderToSQLErr(t *testing.T) {
	_, _, err := Delete("").ToSQL()
	assert.Error(t, err)
}