package resources

type Ticker struct {
	Base   string `json:"base"`
	Target string `json:"target"`
	Price  string `json:"price"`
}
type ResponceBTC struct {
	Ticker  Ticker `json:"ticker"`
	Success bool   `json:"success"`
	Err     string `json:"err"`
}
