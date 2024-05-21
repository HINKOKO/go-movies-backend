package models

type Genre struct {
	ID        int    `json:"id"`
	Genre     string `json:"genre"`
	Checked   bool   `json:"checked"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}
