package invoice

type Invoice struct {
	ID             string `json:"id" db:"id"`
	Amount         uint   `json:"amount" db:"amount"`
	Denomination   string `json:"denomination" db:"denomination"`
	PaymentAddress string `json:"address" db:"address"`
}
