package api

import (
	"context"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type LongURLReponse struct {
	URL string `json:"url"`
}

func ErrorReponse(msg string) Response {
	return Response{
		Error: msg,
	}
}

func ResponseOK(alias string) Response {
	return Response{
		Alias: alias,
	}
}

func URLResponse(url string) LongURLReponse {
	return LongURLReponse{
		URL: url,
	}
}

type Service interface {
	GetShortURL(ctx context.Context, url string) (string, error)
	GetLongURL(ctx context.Context, alias string) (string, error)
}
