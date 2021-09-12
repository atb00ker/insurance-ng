package controllers

type ResponseMessage struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	ResponseSuccess string = "success"
	ResponseError   string = "error"
)
const (
	IsUserLoggedInErrorMessage string = "Key not found, is user logged in?"
)
