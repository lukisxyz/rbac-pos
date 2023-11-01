package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxKey struct{}

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(
		ctx, TxKey{}, tx,
	)
}

func ExtractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(TxKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
