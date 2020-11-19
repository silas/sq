package sq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectBuilderToSQL(t *testing.T) {
	subQ := Select("aa", "bb").From("dd")
	b := Select("a", "b").
		Prefix("WITH prefix AS ?", 0).
		Distinct().
		Columns("c").
		Column("IF(d IN ("+Placeholders(3)+"), 1, 0) as stat_column", 1, 2, 3).
		Column(Expr("a > ?", 100)).
		Column(Alias(Eq{"b": []int{101, 102, 103}}, "b_alias")).
		Column(Alias(subQ, "subq")).
		From("e").
		JoinClause("CROSS JOIN j1").
		Join("j2").
		LeftJoin("j3").
		RightJoin("j4").
		Where("f = ?", 4).
		Where(Eq{"g": 5}).
		Where(map[string]interface{}{"h": 6}).
		Where(Eq{"i": []int{7, 8, 9}}).
		Where(Or{Expr("j = ?", 10), And{Eq{"k": 11}, Expr("true")}}).
		GroupBy("l").
		Having("m = n").
		OrderBy("o ASC", "p DESC").
		Limit(12).
		Offset(13).
		Suffix("FETCH FIRST ? ROWS ONLY", 14)

	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL :=
		"WITH prefix AS $1 " +
			"SELECT DISTINCT a, b, c, IF(d IN ($2,$3,$4), 1, 0) as stat_column, a > $5, " +
			"(b IN ($6,$7,$8)) AS b_alias, " +
			"(SELECT aa, bb FROM dd) AS subq " +
			"FROM e " +
			"CROSS JOIN j1 JOIN j2 LEFT JOIN j3 RIGHT JOIN j4 " +
			"WHERE f = $9 AND g = $10 AND h = $11 AND i IN ($12,$13,$14) AND (j = $15 OR (k = $16 AND true)) " +
			"GROUP BY l HAVING m = n ORDER BY o ASC, p DESC LIMIT 12 OFFSET 13 " +
			"FETCH FIRST $17 ROWS ONLY"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{0, 1, 2, 3, 100, 101, 102, 103, 4, 5, 6, 7, 8, 9, 10, 11, 14}
	assert.Equal(t, expectedArgs, args)
}

func BenchmarkSelectBuilderToSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {

		qb := Select("a", "b").
			Prefix("WITH prefix AS ?", 0).
			Distinct().
			Columns("c").
			Column("IF(d IN ("+Placeholders(3)+"), 1, 0) as stat_column", 1, 2, 3).
			Column(Expr("a > ?", 100)).
			Column(Eq{"b": []int{101, 102, 103}}).
			From("e").
			JoinClause("CROSS JOIN j1").
			Join("j2").
			LeftJoin("j3").
			RightJoin("j4").
			Where("f = ?", 4).
			Where(Eq{"g": 5}).
			Where(map[string]interface{}{"h": 6}).
			Where(Eq{"i": []int{7, 8, 9}}).
			Where(Or{Expr("j = ?", 10), And{Eq{"k": 11}, Expr("true")}}).
			GroupBy("l").
			Having("m = n").
			OrderBy("o ASC", "p DESC").
			Limit(12).
			Offset(13).
			Suffix("FETCH FIRST ? ROWS ONLY", 14)

		_, _, _ = qb.ToSQL()
	}
}

func TestSelectBuilderZeroOffsetLimit(t *testing.T) {
	qb := Select("a").
		From("b").
		Limit(0).
		Offset(0)

	sql, _, err := qb.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "SELECT a FROM b LIMIT 0 OFFSET 0"
	assert.Equal(t, expectedSQL, sql)
}

func TestSelectBuilderToSQLErr(t *testing.T) {
	_, _, err := Select().From("x").ToSQL()
	assert.Error(t, err)
}