package dto

// Param : parameter data structure
type Param struct {
	FromAccountNmbr int
	ToAccountNmbr   int
	Amount          float64
}

type ReqTransferParam struct {
	ToAccountNmbr int     `json:"to_account_number"`
	Amount        float64 `json:"amount"`
}
