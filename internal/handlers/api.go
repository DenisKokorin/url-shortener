package api

import "net/http"

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

type LongURLReponse struct {
	Status int    `json:"status"`
	URL    string `json:"url"`
}

func ErrorReponse(status int, msg string) Response {
	return Response{
		Status: status,
		Error:  msg,
	}
}

func ResponseOK(alias string) Response {
	return Response{
		Status: http.StatusCreated,
		Alias:  alias,
	}
}

func URLResponse(url string) LongURLReponse {
	return LongURLReponse{
		Status: http.StatusOK,
		URL:    url,
	}
}

type Service interface {
	GetShortURL(url string) (string, error)
	GetLongURL(alias string) (string, error)
}
