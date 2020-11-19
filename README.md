# sq

Fluent SQL generator for Go.

``` go
package main

import (
	"context"
	"log"

	"github.com/silas/sq"
)

func main() {
	ctx := context.Background()
	pool, err := sq.Connect(ctx, "postgresql://postgres:postgres@127.0.0.1:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}

	qb := sq.Select("user").
		From("user").
		Limit(1)

	var user string
	err = pool.Tx(ctx, func(tx sq.Tx) error {
		return tx.QueryRow(ctx, qb).Scan(&user)
	})
	if err != nil {
		log.Fatal(err)
	}

	println(user)
}
```

This is a fork of [sqrl][sqrl].

[sqrl]: https://github.com/elgris/sqrl
