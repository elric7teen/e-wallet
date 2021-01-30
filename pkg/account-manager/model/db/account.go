package model

// Customer :
type Customer struct {
	CustomerNumber   int    `json:"customer_number" gorm:"column:customer_number"`
	CustomerName     string `json:"customer_name" gorm:"column:customer_name"`
	CustomerPassword string `json:"customer_password" gorm:"column:customer_password"`
}

// CustomerAccount :
type CustomerAccount struct {
	AccountNumber  int     `json:"account_number" gorm:"column:account_number"`
	CustomerNumber int     `json:"customer_number" gorm:"column:customer_number"`
	Balance        float64 `json:"account_balance" gorm:"column:account_balance"`
}

// AccountInfo :
type AccountInfo struct {
	AccountNumber int     `json:"account_number" gorm:"column:account_number"`
	CustomerName  string  `json:"customer_name" gorm:"column:customer_name"`
	Balance       float64 `json:"account_balance" gorm:"column:account_balance"`
}

// NewCustomerAccount : Customer Account Builder
func (ca CustomerAccount) NewCustomerAccount(accNmbr, custNmbr int, balance float64) *CustomerAccount {
	ca.AccountNumber = accNmbr
	ca.CustomerNumber = custNmbr
	ca.Balance = balance
	return &ca
}

// NewCustomer : Customer Builder
func (c Customer) NewCustomer(custNmbr int, custName, custPass string) *Customer {
	c.CustomerNumber = custNmbr
	c.CustomerName = custName
	c.CustomerPassword = custPass
	return &c
}
