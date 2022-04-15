package models

type Category struct {
	ID       uint
	Category string   `json:"category_name"`
	Posts    *[]Posts `json:"-"`
}
