package models

import "gorm.io/gorm"

// User Object/Model
type User struct {
	gorm.Model
	Name     string
	Role     string
	Email    string `gorm:"unique"`
	Password string `json:"-"`
}

type Author struct {
	gorm.Model
	Name  string
	Email string
	Books []Book `gorm:"foreignKey:AuthorID"`
}

type Category struct {
	gorm.Model
	Title string
	Books []Book `gorm:"foreignKey:CategoryID"`
}

type Book struct {
	gorm.Model
	Title       string
	ImageUrl    string
	Price       float64
	Copies      int64
	Description string
	Trending    bool
	AuthorID    uint
	CategoryID  uint
	// Author      Author `gorm:"foreignKey:AuthorID"`
	Borrows []Borrow
}

type Borrow struct {
	gorm.Model
	UserID  uint
	BookID  uint
	DueDate string
	Book    Book `gorm:"foreignKey:BookID"`
	User    User `gorm:"foreignKey:UserID"`
}
