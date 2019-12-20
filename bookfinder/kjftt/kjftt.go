package kjftt

import (
	"net/http"
	"time"
)

//KJFTT finder for KJFTT
type KJFTT struct {
	client *http.Client
}

//NewKJFTT creates new KJFTT finder
func NewKJFTT(timeout time.Duration) *KJFTT {
	return &KJFTT{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}
