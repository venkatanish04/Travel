package booking

type Flight struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	Seats []Seat `json:"seats"`
}
