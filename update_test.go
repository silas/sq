package sq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBuilderToSQL(t *testing.T) {
	b := Update("").
		Prefix("WITH prefix AS ?", 0).
		Table("a").
		Set("b", Expr("? + 1", 1)).
		SetMap(Eq{"c": 2}).
		Where("d = ?", 3).
		OrderBy("e").
		Limit(4).
		Offset(5).
		Suffix("RETURNING ?", 6)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL :=
		"WITH prefix AS $1 " +
			"UPDATE a SET b = $2 + 1, c = $3 WHERE d = $4 " +
			"ORDER BY e LIMIT 4 OFFSET 5 " +
			"RETURNING $5"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{0, 1, 2, 3, 6}
	assert.Equal(t, expectedArgs, args)
}

func TestUpdateBuilderZeroOffsetLimit(t *testing.T) {
	qb := Update("a").
		Set("b", true).
		Limit(0).
		Offset(0)

	sql, args, err := qb.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "UPDATE a SET b = $1 LIMIT 0 OFFSET 0"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{true}
	assert.Equal(t, expectedArgs, args)
}

func TestUpdateBuilderToSQLErr(t *testing.T) {
	_, _, err := Update("").Set("x", 1).ToSQL()
	assert.Error(t, err)

	_, _, err = Update("x").ToSQL()
	assert.Error(t, err)
}