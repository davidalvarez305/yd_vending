package models

import "time"

type Transaction struct {
	TransactionID int       `json:"transaction_id" gorm:"primaryKey"`
	Amount        int       `gorm:"not null" json:"amount"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	Comments      string    `json:"comments"`
	IsFixed       bool      `gorm:"not null" json:"is_fixed"`
	IsExpense     bool      `gorm:"not null" json:"is_expense"`
	CategoryID    int       `json:"category_id" form:"category_id"`
	Category      *Category `gorm:"not null;column:category_id;foreignKey:CategoryID;constraint:OnDelete:NULL,OnUpdate:CASCADE" json:"category" form:"category"`
	UserID        int       `json:"user_id" form:"user_id"`
	User          *User     `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:DELETE,OnUpdate:CASCADE" json:"user" form:"user"`
}

type Category struct {
	CategoryID int    `json:"category_id" gorm:"primaryKey"`
	Name       string `gorm:"not null" json:"name"`
	UserID     int    `json:"user_id" form:"user_id"`
	User       *User  `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:DELETE,OnUpdate:CASCADE" json:"user" form:"user"`
}

type User struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null" json:"username" form:"username"`
	Password  string    `gorm:"not null" json:"password" form:"password"`
	Email     string    `gorm:"unique;not null" json:"email" form:"email"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}
