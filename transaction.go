package southxchange

type Transaction struct {
	Date          Timestamp
	Type          string
	Amount        float64
	TotalBalance  float64
	Price         float64
	OtherAmount   float64
	OtherCurrency string
	OrderCode     string
	Status        string
	Address       string
	Hash          string
}
