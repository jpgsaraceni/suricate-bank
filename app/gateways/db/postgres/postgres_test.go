package postgres

//TODO: separate test functions and testcases

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=suricate",
			"POSTGRES_DB=suricate",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://suricate:secret@%s/suricate?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dockerPool.MaxWait = 10 * time.Second
	if err = dockerPool.Retry(func() error {
		dbPool, err = ConnectPool(databaseUrl)

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	migration, err := os.ReadFile("./migrations/000001_accounts.up.sql")

	if err != nil {
		log.Fatalf("Could not read migration file: %s", err)
	}

	if _, err := dbPool.Exec(context.Background(), string(migration)); err != nil {
		log.Fatalf("Could not run migration: %s", err)
	}

	//Run tests
	code := m.Run()

	dbPool.Close()

	if err := dockerPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

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

	testCases := []testCase{
		{
			name: "successfully create account",
			account: account.Account{
				Id:        testId,
				Cpf:       testCpf,
				Name:      "Nice name",
				Secret:    testHash,
				CreatedAt: time.Now(),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			fmt.Println(balance)

			if err := repo.DebitAccount(testContext, account.Id, testMoney10); !errors.Is(err, tt.err) {
				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}

			balance, err = repo.GetBalance(testContext, account.Id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}

			fmt.Println(balance)
		})
	}
}
