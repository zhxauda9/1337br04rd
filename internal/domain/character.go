package domain

type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Gender  string `json:"gender"`
	Type    string `json:"type"`
}
