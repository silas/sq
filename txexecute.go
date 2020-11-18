// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package sq

import (
	"context"
)

func txExecute(ctx context.Context, tx *Transaction, fn func(Tx) error) (err error) {
	defer func() {
		if err == nil {
			// Ignore commit errors. The tx has already been committed by RELEASE.
			_ = tx.tx.Commit(ctx)
		} else {
			// We always need to execute a Rollback() so sql.DB releases the
			// connection.
			_ = tx.tx.Rollback(ctx)
		}
	}()
	// Specify that we intend to retry this txn in case of CockroachDB retryable
	// errors.
	if _, err = tx.tx.Exec(ctx, "SAVEPOINT sq"); err != nil {
		return err
	}

	for {
		err = fn(tx)
		if err == nil {
			// RELEASE acts like COMMIT in CockroachDB. We use it since it gives us an
			// opportunity to react to retryable errors, whereas tx.Commit() doesn't.
			if _, err = tx.tx.Exec(ctx, "RELEASE SAVEPOINT sq"); err == nil {
				return nil
			}
		}
		// We got an error; let's see if it's a retryable one and, if so, restart.
		if !IsError(err, CodeSerializationFailure) {
			return err
		}

		if _, retryErr := tx.tx.Exec(ctx, "ROLLBACK TO SAVEPOINT sq"); retryErr != nil {
			return err
		}
	}
}
