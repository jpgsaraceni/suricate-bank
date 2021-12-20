package transfer

type Repository interface {
	Create(transfer *Transfer) error
}
