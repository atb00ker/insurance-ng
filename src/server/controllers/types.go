package controllers

// ResponseMessage is the generic response type send on any error
type ResponseMessage struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// Response message status options
const (
	ResponseSuccess string = "success"
	ResponseError   string = "error"
)

// Response common error message options
const (
	IsUserLoggedInErrorMessage string = "key not found, is user logged in?"
)
