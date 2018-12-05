package southxchange

type Balance struct {
	Currency    string
	Deposited   float64
	Available   float64
	Unconfirmed float64
}
