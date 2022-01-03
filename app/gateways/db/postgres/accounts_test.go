package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestAccount(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		account account.Account
		err     error
	}

	testId := account.AccountId(uuid.New())
	testCpf, _ := cpf.NewCpf("22061446035")
	testHash, _ := hash.NewHash("nicesecret")
	testMoney10, _ := money.NewMoney(10)
	testMoney30, _ := money.NewMoney(30)

	tt := testCase{
		name: "test account",
		account: account.Account{
			Id:        testId,
			Cpf:       testCpf,
			Name:      "Nice name",
			Secret:    testHash,
			CreatedAt: time.Now(),
		},
	}

	testContext := context.Background()

	repo := accountspg.NewRepository(dbPool)
	if err := repo.Create(testContext, &tt.account); !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}
	accounts, err := repo.Fetch(testContext)

	if !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}

	account, err := repo.GetById(testContext, accounts[0].Id)

	if !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}

	if err := repo.CreditAccount(testContext, account.Id, testMoney30); !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}

	balance, err := repo.GetBalance(testContext, account.Id)

	if !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}

	if balance != 30 {
		t.Fatalf("expected balance: 30 got balance: %d", balance)
	}

	if err := repo.DebitAccount(testContext, account.Id, testMoney10); !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}

	balance, err = repo.GetBalance(testContext, account.Id)

	if !errors.Is(err, tt.err) {
		t.Fatalf("expected error: %s got error: %s", tt.err, err)
	}

	if balance != 20 {
		t.Fatalf("expected balance: 20 got balance: %d", balance)
	}
}
