package models

import (
	"time"
)

//easyjson:json
type MoneyCategory struct {
	ID       uint    `json:"id" db:"category_id"`
	User     uint    `json:"-" db:"user_id"`
	Name     string  `json:"name,omitempty"`
	IsIncome bool    `json:"isIncome,omitempty" db:"is_income"` // FALSE – consumption
	Pic      *string `json:"picture,omitempty" db:"pic"`

	Sum int `json:"value,omitempty"`
}

//easyjson:json
type MoneyCategoryList []MoneyCategory

//easyjson:json
type MoneyOp struct {
	UUID  string  `json:"id,omitempty" db:"action_uuid"`
	User  uint    `json:"-" db:"user_id"`
	Delta int     `json:"value"`
	From  int     `json:"from" db:"from_category"`
	To    int     `json:"to" db:"to_category"`
	Photo *string `json:"photo,omitempty"`

	Added time.Time `json:"added"`
}

//easyjson:json
type MoneyOpHistory []MoneyOp
