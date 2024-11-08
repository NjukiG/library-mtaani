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
	// Books []Book `gorm:"foreignKey:AuthorID"`
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
	Price       int
	Copies      int64
	Description string
	Trending    bool
	// AuthorID    uint
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

// Cart Model
type Cart struct {
	gorm.Model
	UserID    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;"`
}

// CartItem Model
type CartItem struct {
	gorm.Model
	CartID   uint `json:"cart_id"`
	BookID   uint `json:"book_id"`
	Book     Book `json:"book" gorm:"foreignKey:BookID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Quantity int  `json:"quantity" gorm:"not null;default:1"`
	UserID   uint `json:"user_id"`
}
