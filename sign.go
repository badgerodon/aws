package aws

import (
	"net/http"
)

const (
	ISO8601BasicFormat      = "20060102T150405Z"
	ISO8601BasicFormatShort = "20060102"
)

type Signer interface {
	Sign(req *http.Request)
}
