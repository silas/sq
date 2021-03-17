package sq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithSingle(t *testing.T) {
	qb := With("hello").
		As(
			Select("name").
				From("test").
				Where(Gt{"age": 50}),
		).
		Select("name").
		From("hello").
		Where("name ILIKE ?", "m%")
	sql, args, err := qb.ToSQL()

	assert.NoError(t, err)

	expectedSQL := "WITH " +
		"hello AS (SELECT name FROM test WHERE age > ?) " +
		"SELECT name FROM hello WHERE name ILIKE ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{50, "m%"}
	assert.Equal(t, expectedArgs, args)
}

func TestWithMultiple(t *testing.T) {
	qb := With("one").As(Select("a1").From("a")).
		With("two").As(Select("b1").From("b")).
		Select("a1").
		From("one")
	sql, args, err := qb.ToSQL()

	assert.NoError(t, err)

	expectedSQL := "WITH " +
		"one AS (SELECT a1 FROM a), " +
		"two AS (SELECT b1 FROM b) " +
		"SELECT a1 FROM one"
	assert.Equal(t, expectedSQL, sql)

	assert.Empty(t, args)
}

func TestWithUpdate(t *testing.T) {
	qb := With("one").As(Select("a1").From("a")).
		With("two").As(Update("b").SetMap(map[string]interface{}{"b1": 1})).
		Select("a1").
		From("one")
	sql, args, err := qb.ToSQL()

	assert.NoError(t, err)

	expectedSQL := "WITH " +
		"one AS (SELECT a1 FROM a), " +
		"two AS (UPDATE b SET b1 = ?) " +
		"SELECT a1 FROM one"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestWithRecursive(t *testing.T) {
	qb := With("search_graph", "id", "link", "data", "depth").
		Recursive().
		As(
			Select("g.id", "g.link", "g.data", "1").
				From("graph g"),
		).
		UnionAll(
			Select("g.id", "g.link", "g.data", "sg.depth + 1").
				From("graph g, search_graph sg").
				Where("g.id = sg.link"),
		).
		Select("*").
		From("search_graph")
	sql, args, err := qb.ToSQL()

	assert.NoError(t, err)

	expectedSQL := "WITH RECURSIVE " +
		"search_graph(id, link, data, depth) AS (" +
		"SELECT g.id, g.link, g.data, 1 FROM graph g " +
		"UNION ALL " +
		"SELECT g.id, g.link, g.data, sg.depth + 1 FROM graph g, search_graph sg WHERE g.id = sg.link" +
		") SELECT * FROM search_graph"
	assert.Equal(t, expectedSQL, sql)

	assert.Empty(t, args)
}

func TestWithError(t *testing.T) {
	_, _, err := NewWithBuilder().
		As(Select("t").From("test")).
		Select("t").From("test").
		ToSQL()
	require.EqualError(t, err, "with statements must have WITH")

	_, _, err = With("").
		As(Select("t").From("test")).
		Select("t").From("hello").
		ToSQL()
	require.EqualError(t, err, "with statements must have a non-empty name")

	_, _, err = With("hello").
		Select("t").From("hello").
		ToSQL()
	require.EqualError(t, err, "with statements must have AS statement")
}
