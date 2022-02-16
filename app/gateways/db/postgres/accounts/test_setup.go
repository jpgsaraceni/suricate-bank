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

const (
	testAmount10 = 10
	testAmount20 = 20
	testAmount30 = 30
)

var (
	testContext    = context.Background()
	testHash, _    = hash.NewHash("nicesecret")
	testTime       = time.Now().Round(time.Hour)
	testMoney10, _ = money.NewMoney(testAmount10)
	testMoney20, _ = money.NewMoney(testAmount20)
	testMoney30, _ = money.NewMoney(testAmount30)
)

func createTestAccount(pool *pgxpool.Pool, id account.ID, cpf string, amount int) error {
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

func CreateTestAccountBatch(pool *pgxpool.Pool, ids []account.ID, cpfs []string, amount []int) error {
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

	if _, err := br.Exec(); err != nil {
		return err
	}
	br.Close()

	return nil
}
