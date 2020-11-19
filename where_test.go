package sq

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWherePartsAppendToSQL(t *testing.T) {
	parts := []StatementBuilder{
		newWherePart("x = ?", 1),
		newWherePart(nil),
		newWherePart(Eq{"y": 2}),
	}
	sql := &bytes.Buffer{}
	args, _ := appendToSQL(parts, sql, " AND ", []interface{}{})
	assert.Equal(t, "x = ? AND y = ?", sql.String())
	assert.Equal(t, []interface{}{1, 2}, args)
}

func TestWherePartsAppendToSQLErr(t *testing.T) {
	parts := []StatementBuilder{newWherePart(1)}
	_, err := appendToSQL(parts, &bytes.Buffer{}, "", []interface{}{})
	assert.Error(t, err)
}

func TestWherePartNil(t *testing.T) {
	sql, _, _ := newWherePart(nil).ToSQL()
	assert.Equal(t, "", sql)
}

func TestWherePartErr(t *testing.T) {
	_, _, err := newWherePart(1).ToSQL()
	assert.Error(t, err)
}

func TestWherePartString(t *testing.T) {
	sql, args, _ := newWherePart("x = ?", 1).ToSQL()
	assert.Equal(t, "x = ?", sql)
	assert.Equal(t, []interface{}{1}, args)
}

func TestWherePartMap(t *testing.T) {
	test := func(pred interface{}) {
		sql, _, _ := newWherePart(pred).ToSQL()
		expect := []string{"x = ? AND y = ?", "y = ? AND x = ?"}
		if sql != expect[0] && sql != expect[1] {
			t.Errorf("expected one of %#v, got %#v", expect, sql)
		}
	}
	m := map[string]interface{}{"x": 1, "y": 2}
	test(m)
	test(Eq(m))
}

func TestWherePartNoArgs(t *testing.T) {
	_, _, err := newWherePart(Eq{"test": []string{}}).ToSQL()
	assert.Equal(t, err, fmt.Errorf("equality condition must contain at least one paramater"))
}

func TestWhereBuilderToSelectSQL(t *testing.T) {
	sql, args, err := Where(Eq{"hello": "world"}).
		Select("name").
		From("test").
		ToSQL()
	require.NoError(t, err)

	expectedSQL := "SELECT name FROM test WHERE hello = $1"
	require.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"world"}
	require.Equal(t, expectedArgs, args)
}

func TestWhereBuilderToDeleteSQL(t *testing.T) {
	sql, args, err := Where(Eq{"hello": "world"}).
		Delete("test").
		ToSQL()
	require.NoError(t, err)

	expectedSQL := "DELETE FROM test WHERE hello = $1"
	require.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"world"}
	require.Equal(t, expectedArgs, args)
}

func TestWhereBuilderToUpdateSQL(t *testing.T) {
	sql, args, err := Where(Eq{"hello": "world"}).
		Update("test").
		Set("hello", "universe").
		ToSQL()
	require.NoError(t, err)

	expectedSQL := "UPDATE test SET hello = $1 WHERE hello = $2"
	require.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"universe", "world"}
	require.Equal(t, expectedArgs, args)
}
