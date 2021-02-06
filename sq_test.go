package sq

import (
	"context"
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
	column := "v"
	input := "hello123"

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
		createTable := fmt.Sprintf("CREATE TABLE %s (%s text)", table, column)
		_, err = tx.Exec(ctx, Expr(createTable))
		require.NoError(t, err)

		_, err = tx.Exec(ctx, Insert(table).
			SetMap(map[string]interface{}{column: input}))
		require.NoError(t, err)

		var output string
		rows, err := tx.Query(ctx, Select(column).From(table))
		require.NoError(t, err)
		require.True(t, rows.Next())
		err = rows.Scan(&output)
		require.Equal(t, input, output)
		require.False(t, rows.Next())

		var output2 string
		err = tx.QueryRow(ctx, Select(column).From(table)).Scan(&output2)
		require.NoError(t, err)
		require.Equal(t, input, output2)

		return nil
	})
	require.NoError(t, err)

	var output string
	err = pool.QueryRow(ctx, Select(column).From(table)).Scan(&output)
	require.NoError(t, err)
	require.Equal(t, input, output)

	_, err = pool.Exec(ctx, Delete(table))
	require.NoError(t, err)

	rows, err := pool.Query(ctx, Select(column).From(table))
	require.NoError(t, err)
	require.False(t, rows.Next())
}
