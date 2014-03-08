package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
)

var b64 = base64.StdEncoding

type V2Signer struct {
	auth        Auth
	serviceName string
}

func NewV2Signer(auth Auth, serviceName string) *V2Signer {
	return &V2Signer{auth: auth, serviceName: serviceName}
}

func (s *V2Signer) Sign(req *http.Request) {
	query := req.URL.Query()
	query.Set("AWSAccessKeyId", s.auth.AccessKeyID)
	query.Set("SignatureVersion", "2")
	query.Set("SignatureMethod", "HmacSHA256")

	// AWS specifies that the parameters in a signed request must
	// be provided in the natural order of the keys. This is distinct
	// from the natural order of the encoded value of key=value.
	// Percent and gocheck.Equals affect the sorting order.
	ks := []string{}
	for k, vs := range req.URL.Query() {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	qs := []string{}
	for _, k := range ks {
		for _, v := range query[k] {
			qs = append(qs, Encode(k)+"="+Encode(v))
		}
	}
	joined := strings.Join(qs, "&")
	payload := req.Method + "\n" +
		req.URL.Host + "\n" +
		req.URL.Path + "\n" +
		joined
	hash := hmac.New(sha256.New, []byte(s.auth.SecretAccessKey))
	hash.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	query.Set("Signature", signature)
}
