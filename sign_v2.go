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
	auth Auth
}

func NewV2Signer(auth Auth) *V2Signer {
	return &V2Signer{auth: auth}
}

func (s *V2Signer) stringToSign(req *http.Request) string {
	query := req.URL.Query()
	// AWS specifies that the parameters in a signed request must
	// be provided in the natural order of the keys. This is distinct
	// from the natural order of the encoded value of key=value.
	// Percent and gocheck.Equals affect the sorting order.
	ks := []string{}
	for k, _ := range query {
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
	path := req.URL.Path
	if path == "" {
		path = "/"
	}
	return req.Method + "\n" +
		req.URL.Host + "\n" +
		path + "\n" +
		joined
}

func (s *V2Signer) Sign(req *http.Request) {
	query := req.URL.Query()
	query.Set("AWSAccessKeyId", s.auth.AccessKeyID)
	query.Set("SignatureVersion", "2")
	query.Set("SignatureMethod", "HmacSHA256")
	req.URL.RawQuery = query.Encode()

	payload := s.stringToSign(req)
	hash := hmac.New(sha256.New, []byte(s.auth.SecretAccessKey))
	hash.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	query.Set("Signature", signature)
	req.URL.RawQuery = query.Encode()
}
