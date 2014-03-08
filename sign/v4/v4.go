package v4

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

type closableBuffer struct {
	*bytes.Buffer
}

var (
	maybeStringRE = regexp.MustCompile(`^([^"]*)("[^"]*")?(.*)$`)
	spaceRE       = regexp.MustCompile(`\s+`)

	Algorithm = "AWS4-HMAC-SHA256"
)

func (this closableBuffer) Close() error {
	return nil
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func getSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kSecret := []byte("AWS4" + key)
	kDate := hmacSHA256(kSecret, dateStamp)
	kRegion := hmacSHA256(kDate, regionName)
	kService := hmacSHA256(kRegion, serviceName)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}

func hashedPayload(req *http.Request) string {
	var buf bytes.Buffer
	rdr := io.TeeReader(req.Body, &buf)
	defer req.Body.Close()

	h := sha256.New()
	// todo: handle error?
	io.Copy(h, rdr)
	req.Body = closableBuffer{&buf}
	return hex.EncodeToString(h.Sum(nil))
}

func signedHeaders(req *http.Request) string {
	hs := make([]string, 0, len(req.Header))
	for h, _ := range req.Header {
		hs = append(hs, strings.ToLower(h))
	}
	sort.Strings(hs)
	return strings.Join(hs, ";")
}

func trimNonQuoted(str string) string {
	return strings.TrimSpace(spaceRE.ReplaceAllLiteralString(str, " "))
}
func trimAll(str string) string {
	matches := maybeStringRE.FindStringSubmatch(str)
	if len(matches) < 4 {
		return trimNonQuoted(str)
	}
	return trimNonQuoted(matches[1]) + matches[2] + trimNonQuoted(matches[3])
}

func canonicalHeaders(req *http.Request) string {
	str := ""
	hs := make([]string, 0, len(req.Header))
	for h, _ := range req.Header {
		hs = append(hs, strings.ToLower(h))
	}
	sort.Strings(hs)

	for _, h := range hs {
		vs := req.Header[http.CanonicalHeaderKey(h)]
		for i, v := range vs {
			vs[i] = trimAll(v)
		}
		str += h + ":" + strings.Join(vs, ",") + "\n"
	}
	return str
}

func canonicalQueryString(req *http.Request) string {
	str := ""
	vs := req.URL.Query()
	keys := make([]string, 0, len(vs))
	for k, _ := range vs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			str += "&"
		}
		str += url.QueryEscape(k) + "=" + url.QueryEscape(vs.Get(k))
	}

	return str
}

func canonicalURI(req *http.Request) string {
	path := req.URL.Path
	if path == "" {
		path = "/"
	}
	return path
}

func canonicalRequest(req *http.Request) string {
	return req.Method + "\n" +
		canonicalURI(req) + "\n" +
		canonicalQueryString(req) + "\n" +
		canonicalHeaders(req) + "\n" +
		signedHeaders(req) + "\n" +
		hashedPayload(req)
}

func hashedCanonicalRequest(req *http.Request) string {
	h := sha256.New()
	h.Write([]byte(canonicalRequest(req)))
	return hex.EncodeToString(h.Sum(nil))
}

func Sign(req *http.Request, accessKeyId, secretAccessKey, regionName, serviceName string) error {
	date := req.Header.Get("x-amz-date")
	if date == "" {
		date = time.Now().Format(time.RFC3339)
		req.Header.Set("x-amz-date", date)
	}

	credentialScope := date[:8] + "/" + regionName + "/" + serviceName + "/aws4_request"

	stringToSign := Algorithm + "\n" +
		date + "\n" +
		credentialScope + "\n" +
		hashedCanonicalRequest(req)

	key := getSignatureKey(secretAccessKey, date[:8], regionName, serviceName)

	signature := hex.EncodeToString(hmacSHA256(key, stringToSign))

	req.Header.Set("Authorization",
		Algorithm+" "+
			"Credential="+accessKeyId+"/"+credentialScope+", "+
			"SignedHeaders="+signedHeaders(req)+", "+
			"Signature="+signature)

	return nil
}
