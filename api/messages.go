package api

//ErrDetail ...
type ErrDetail struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

//ErrMessage is return message default
type ErrMessage struct {
	Message string      `json:"message"`
	Errors  []ErrDetail `json:"errors"`
}

//SuccessMessage is return Zen message
type SuccessMessage struct {
	Message string `json:"message"`
}
