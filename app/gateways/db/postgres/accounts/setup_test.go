package accountspg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

var (
	testContext    = context.Background()
	testHash, _    = hash.NewHash("nicesecret")
	testTime       = time.Now().Round(time.Hour)
	testMoney10, _ = money.NewMoney(10)
	testMoney20, _ = money.NewMoney(20)
	testMoney30, _ = money.NewMoney(30)
)

func createTestAccount(pool *pgxpool.Pool, id account.AccountId, cpf string, amount int) error {
	const query = `
		INSERT INTO 
			accounts (
				id,
				name,
				cpf,
				secret,
				balance,
				created_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6);
	`

	_, err := pool.Exec(
		testContext,
		query,
		id,
		"nice name",
		cpf,
		testHash.Value(),
		amount,
		testTime,
	)

	if err != nil {
		return err
	}

	return nil
}

func createTestAccountBatch(pool *pgxpool.Pool, ids []account.AccountId, cpfs []string, amount []int) error {
	const query = `
	INSERT INTO 
		accounts (
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		)
	VALUES
		($1, $2, $3, $4, $5, $6);
`
	batch := &pgx.Batch{}

	for i := 0; i < len(ids); i++ {
		batch.Queue(query, ids[i], "nice name", cpfs[i], testHash.Value(), amount[i], testTime)
	}

	br := pool.SendBatch(testContext, batch)

	_, err := br.Exec()
	defer br.Close()
	if err != nil {
		return err
	}

	return nil
}
