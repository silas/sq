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

func TestDeleteFromAndWhatDiffer(t *testing.T) {
	b := Delete("b").
		From("a").
		Where("b = ?", 1)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "DELETE b FROM a WHERE b = $1"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteFromAndWhatSame(t *testing.T) {
	b := Delete("a").
		From("a").
		Where("b = ?", 1)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "DELETE FROM a WHERE b = $1"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
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

func TestDeleteSQLMultipleTables(t *testing.T) {
	b := Delete("a1", "a2").
		From("z1 AS a1").
		JoinClause("INNER JOIN a2 ON a1.id = a2.ref_id").
		Join("a3").
		Where("b = ?", 1)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL :=
		"DELETE a1, a2 " +
			"FROM z1 AS a1 " +
			"INNER JOIN a2 ON a1.id = a2.ref_id " +
			"JOIN a3 " +
			"WHERE b = $1"

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

func TestIssue11(t *testing.T) {
	b := Delete("a").
		From("A a").
		Join("B b ON a.c = b.c").
		Where("b.d = ?", 1).
		Limit(2)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "DELETE a FROM A a " +
		"JOIN B b ON a.c = b.c " +
		"WHERE b.d = $1 " +
		"LIMIT 2"

	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}
