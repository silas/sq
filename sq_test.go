package sq

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func databaseURL() string {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgresql://postgres:postgres@127.0.0.1:5432/postgres"
	}
	return url
}

func TestPool(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()

	table := fmt.Sprintf("test_%d", time.Now().Unix())
	column1 := "name"
	column2 := "create_time"
	name1 := "jane"
	name2 := "mike"

	pool, err := Connect(ctx, databaseURL())
	require.NoError(t, err)

	t.Cleanup(func() {
		dropTable := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
		_, err = pool.Exec(ctx, Expr(dropTable))
		if err != nil {
			t.Logf("Failed to drop test table: %s", err.Error())
		}

		pool.Close()
	})

	err = pool.Tx(ctx, func(tx Tx) error {
		createTable := fmt.Sprintf("CREATE TABLE %s (%s text PRIMARY KEY, %s TIMESTAMPTZ)", table, column1, column2)
		_, err = tx.Exec(ctx, Expr(createTable))
		require.NoError(t, err)

		_, err = tx.Exec(ctx, Insert(table).
			SetMap(map[string]interface{}{column1: name1, column2: time.Now()}))
		require.NoError(t, err)

		_, err = tx.Exec(ctx, Insert(table).
			SetMap(map[string]interface{}{column1: name2}))
		require.NoError(t, err)

		var output1, output2 string
		rows, err := tx.Query(ctx, Select(column1).From(table))
		require.NoError(t, err)
		require.True(t, rows.Next())
		err = rows.Scan(&output1)
		require.NoError(t, err)
		require.Equal(t, name1, output1)
		require.True(t, rows.Next())
		err = rows.Scan(&output2)
		require.NoError(t, err)
		require.Equal(t, name2, output2)
		require.False(t, rows.Next())

		var output3 string
		err = tx.QueryRow(ctx, Select(column1).From(table)).Scan(&output3)
		require.NoError(t, err)
		require.Equal(t, name1, output3)

		type Person struct {
			Name       string
			CreateTime *time.Time
		}
		var output4 []*Person
		err = tx.All(ctx, Select(column1).From(table), &output4)
		require.NoError(t, err)
		require.Len(t, output4, 2)
		require.NotNil(t, output4[0])
		require.Equal(t, name1, output4[0].Name)
		require.Nil(t, output4[0].CreateTime)
		require.NotNil(t, output4[1])
		require.Equal(t, name2, output4[1].Name)
		require.Nil(t, output4[1].CreateTime)

		var output5 Person
		err = tx.One(ctx, Select(column1, column2).From(table).Limit(1), &output5)
		require.NoError(t, err)
		require.Equal(t, name1, output5.Name)
		require.False(t, output5.CreateTime.IsZero())

		err = tx.One(ctx, Select(column1, column2).From(table), &output5)
		require.EqualError(t, err, "scany: expected 1 row, got: 2")

		err = tx.One(ctx, Select(column1, column2).From(table).Offset(5), &output5)
		require.True(t, errors.Is(err, ErrNoRows))

		return nil
	})
	require.NoError(t, err)

	var output string
	err = pool.QueryRow(ctx, Select(column1).From(table)).Scan(&output)
	require.NoError(t, err)
	require.Equal(t, name1, output)

	_, err = pool.Exec(ctx, Delete(table))
	require.NoError(t, err)

	rows, err := pool.Query(ctx, Select(column1).From(table))
	require.NoError(t, err)
	require.False(t, rows.Next())
}
