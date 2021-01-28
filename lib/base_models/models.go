package models

const (
	// Success :
	Success = "201"
	// InsufBalance : Insufficient Balance
	InsufBalance = "502"
	// NotFound :
	NotFound = "404"
	// BadRequest :
	BadRequest = "400"
	// InternalServerError :
	InternalServerError = "500"
)

//Response data structure
type Response struct {
	StatusCode string      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

//Result data structure
type Result struct {
	Data  interface{}
	Error error
}

// ToResponseOK create response ok based on result struct
func (r *Result) ToResponseOK() *Response {
	return &Response{
		StatusCode: "200",
		Data:       r.Data,
		Message:    "Request Success",
	}
}

// ResponseBadReq : return bad request
func ResponseBadReq(msg string) *Response {
	return &Response{
		StatusCode: BadRequest,
		Message:    msg,
	}
}

// ResponseNotFound : return status not found
func ResponseNotFound(msg string) *Response {
	return &Response{
		StatusCode: NotFound,
		Message:    msg,
	}
}

func ResponseInsufBalance(msg string) *Response {
	return &Response{
		StatusCode: InsufBalance,
		Message:    msg,
	}
}
