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

	pool, err := Connect(ctx, databaseURL())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	err = pool.Tx(ctx, func(tx Tx) error {
		table := fmt.Sprintf("test_%d", time.Now().Unix())
		column := "v"

		createTable := fmt.Sprintf("CREATE TABLE %s (%s text)", table, column)
		_, err = tx.Exec(ctx, Expr(createTable))
		require.NoError(t, err)

		input := "hello123"
		_, err = tx.Exec(ctx, Insert(table).
			SetMap(map[string]interface{}{column: input}))
		require.NoError(t, err)

		var output string
		err = tx.QueryRow(ctx, Select(column).From(table)).Scan(&output)
		require.NoError(t, err)

		require.Equal(t, input, output)

		return nil
	})
	require.NoError(t, err)
}
