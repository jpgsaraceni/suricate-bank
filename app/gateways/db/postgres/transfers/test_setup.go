package transferspg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

var (
	testContext    = context.Background()
	testTime       = time.Now().Round(time.Hour)
	testMoney10, _ = money.NewMoney(10)
	testMoney15, _ = money.NewMoney(15)
)

func createTestTransferBatch(pool *pgxpool.Pool, ids []transfer.TransferId, originIds, destinationIds []account.AccountId, amount []int) error {
	const query = `
	INSERT INTO 
	transfers (
		id,
		account_origin_id,
		account_destination_id,
		amount,
		created_at
	)
VALUES
	($1, $2, $3, $4, $5);
`
	batch := &pgx.Batch{}

	for i := 0; i < len(ids); i++ {
		batch.Queue(query, ids[i], originIds[i], destinationIds[i], amount[i], testTime)
	}

	br := pool.SendBatch(testContext, batch)

	_, err := br.Exec()
	defer br.Close()
	if err != nil {
		return err
	}

	return nil
}
