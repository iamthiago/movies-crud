package models

type Movie struct {
	ID       int64  `json:"id"`
	Isbn     string `json:"isbn"`
	Title    string `json:"title"`
	Director string `json:"director"`
}

func (m Movie) IsEmpty() bool {
	return m.ID == 0
}
