package models

type Movie struct {
	ID       int64     `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (m Movie) IsEmpty() bool {
	return m.ID == 0
}
