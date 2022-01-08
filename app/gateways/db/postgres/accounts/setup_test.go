package accountspg

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
)

var testContext = context.Background()
var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dbPool = postgrestest.GetTestPool()

	m.Run()
}
