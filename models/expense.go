package models

import "time"

type Expense struct {
	ExpenseAmount   int       `json:"expenseAmount" form:"expenseAmount"`
	ExpenseCategory int       `json:"expenseCategory" form:"expenseCategory"`
	ExpenseDate     time.Time `json:"expenseDate" form:"expenseDate"`
	ExpenseComments string    `json:"expenseComments" form:"expenseComments"`
}
