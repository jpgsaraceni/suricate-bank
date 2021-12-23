package transferuc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository transfer.Repository
		want       []transfer.Transfer
		err        error
	}

	var testUUID1, _ = uuid.NewUUID()
	var testUUID2, _ = uuid.NewUUID()
	var testUUID3, _ = uuid.NewUUID()
	var testUUID4, _ = uuid.NewUUID()

	testCases := []testCase{
		{
			name: "successfully fetch 1 transfer",
			repository: transfer.MockRepository{
				OnFetch: func() ([]transfer.Transfer, error) {

					return []transfer.Transfer{
						{
							Id: transfer.TransferId(testUUID1),
						},
					}, nil
				},
			},
			want: []transfer.Transfer{
				{
					Id: transfer.TransferId(testUUID1),
				},
			},
		},
		{
			name: "successfully fetch 4 transfers",
			repository: transfer.MockRepository{
				OnFetch: func() ([]transfer.Transfer, error) {

					return []transfer.Transfer{
						{
							Id: transfer.TransferId(testUUID1),
						},
						{
							Id: transfer.TransferId(testUUID2),
						},
						{
							Id: transfer.TransferId(testUUID3),
						},
						{
							Id: transfer.TransferId(testUUID4),
						},
					}, nil
				},
			},
			want: []transfer.Transfer{
				{
					Id: transfer.TransferId(testUUID1),
				},
				{
					Id: transfer.TransferId(testUUID2),
				},
				{
					Id: transfer.TransferId(testUUID3),
				},
				{
					Id: transfer.TransferId(testUUID4),
				},
			},
		},
		{
			name: "no existent transfers error",
			repository: transfer.MockRepository{
				OnFetch: func() ([]transfer.Transfer, error) {

					return []transfer.Transfer{}, nil
				},
			},
			want: []transfer.Transfer{},
			err:  errNoTransfersToFetch,
		},
		{
			name: "repository throws error",
			repository: transfer.MockRepository{
				OnFetch: func() ([]transfer.Transfer, error) {

					return []transfer.Transfer{}, errFetchTransfers
				},
			},
			want: []transfer.Transfer{},
			err:  errFetchTransfers,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository, MockCrediter{}, MockDebiter{}}

			transfersList, err := uc.Fetch()

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(transfersList, tt.want) {
				t.Errorf("got %v expected %v", transfersList, tt.want)
			}
		})
	}
}
