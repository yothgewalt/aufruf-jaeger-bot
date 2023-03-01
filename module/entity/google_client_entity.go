package entity

import (
	"net/http"
)

type GoogleClient interface {
	GetClient(jsonKey []byte, scope ...string) *http.Client
}
