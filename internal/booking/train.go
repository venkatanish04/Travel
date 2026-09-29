package booking

type Train struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	Seats []Seat `json:"seats"`
}
