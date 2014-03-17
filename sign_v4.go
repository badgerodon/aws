package aws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

type closableBuffer struct {
	*bytes.Buffer
}

type (
	V4Signer struct {
		auth        Auth
		serviceName string
		regionName  string
	}
)

func NewV4Signer(auth Auth, serviceName string, regionName string) *V4Signer {
	return &V4Signer{auth: auth, serviceName: serviceName, regionName: regionName}
}

func (s *V4Signer) Sign(req *http.Request) {
	req.Header.Set("host", req.Host)
	t := s.requestTime(req)
	sts := s.stringToSign(t, req)
	signature := s.signature(t, sts)
	auth := s.authorization(req.Header, t, signature)
	req.Header.Set("Authorization", auth)
	return
}

func (s *V4Signer) requestTime(req *http.Request) time.Time {
	// Get "x-amz-date" header
	date := req.Header.Get("x-amz-date")

	// Attempt to parse as ISO8601BasicFormat
	t, err := time.Parse(ISO8601BasicFormat, date)
	if err == nil {
		return t
	}

	// Attempt to parse as http.TimeFormat
	t, err = time.Parse(http.TimeFormat, date)
	if err == nil {
		req.Header.Set("x-amz-date", t.Format(ISO8601BasicFormat))
		return t
	}

	// Get "date" header
	date = req.Header.Get("date")

	// Attempt to parse as http.TimeFormat
	t, err = time.Parse(http.TimeFormat, date)
	if err == nil {
		return t
	}

	// Create a current time header to be used
	t = time.Now().UTC()
	req.Header.Set("x-amz-date", t.Format(ISO8601BasicFormat))
	return t
}

func (s *V4Signer) canonicalRequest(req *http.Request) string {
	c := new(bytes.Buffer)
	fmt.Fprintf(c, "%s\n", req.Method)
	fmt.Fprintf(c, "%s\n", s.canonicalURI(req.URL))
	fmt.Fprintf(c, "%s\n", s.canonicalQueryString(req.URL))
	fmt.Fprintf(c, "%s\n\n", s.canonicalHeaders(req.Header))
	fmt.Fprintf(c, "%s\n", s.signedHeaders(req.Header))
	if req.Header.Get("x-amz-content-sha256") == "STREAMING-AWS4-HMAC-SHA256-PAYLOAD" {
		fmt.Fprintf(c, "STREAMING-AWS4-HMAC-SHA256-PAYLOAD")
	} else {
		fmt.Fprintf(c, "%s", s.payloadHash(req))
	}
	str := c.String()
	return str
}

func (s *V4Signer) canonicalURI(u *url.URL) string {
	canonicalPath := u.RequestURI()
	if u.RawQuery != "" {
		canonicalPath = canonicalPath[:len(canonicalPath)-len(u.RawQuery)-1]
	}
	slash := strings.HasSuffix(canonicalPath, "/")
	canonicalPath = path.Clean(canonicalPath)
	if canonicalPath != "/" && slash {
		canonicalPath += "/"
	}
	return canonicalPath
}

func (s *V4Signer) canonicalQueryString(u *url.URL) string {
	var a []string
	for k, vs := range u.Query() {
		k = Encode(k)
		for _, v := range vs {
			if v == "" {
				a = append(a, k)
			} else {
				v = Encode(v)
				a = append(a, k+"="+v)
			}
		}
	}
	sort.Strings(a)
	return strings.Join(a, "&")
}

func (s *V4Signer) canonicalHeaders(h http.Header) string {
	i, a := 0, make([]string, len(h))
	for k, v := range h {
		for j, w := range v {
			v[j] = strings.Trim(w, " ")
		}
		sort.Strings(v)
		a[i] = strings.ToLower(k) + ":" + strings.Join(v, ",")
		i++
	}
	sort.Strings(a)
	return strings.Join(a, "\n")
}

func (s *V4Signer) signedHeaders(h http.Header) string {
	i, a := 0, make([]string, len(h))
	for k, _ := range h {
		a[i] = strings.ToLower(k)
		i++
	}
	sort.Strings(a)
	return strings.Join(a, ";")
}

func (s *V4Signer) payloadHash(req *http.Request) string {
	var b []byte
	if req.Body == nil {
		b = []byte("")
	} else {
		var err error
		b, err = ioutil.ReadAll(req.Body)
		if err != nil {
			// TODO: I REALLY DON'T LIKE THIS PANIC!!!!
			panic(err)
		}
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return s.hash(string(b))
}

func (s *V4Signer) stringToSign(t time.Time, req *http.Request) string {
	w := new(bytes.Buffer)
	fmt.Fprint(w, "AWS4-HMAC-SHA256\n")
	fmt.Fprintf(w, "%s\n", t.Format(ISO8601BasicFormat))
	fmt.Fprintf(w, "%s\n", s.credentialScope(t))
	fmt.Fprintf(w, "%s", s.hash(s.canonicalRequest(req)))
	return w.String()
}

func (s *V4Signer) credentialScope(t time.Time) string {
	return fmt.Sprintf("%s/%s/%s/aws4_request", t.Format(ISO8601BasicFormatShort), s.regionName, s.serviceName)
}

func (s *V4Signer) signature(t time.Time, sts string) string {
	h := s.hmac(s.derivedKey(t), []byte(sts))
	return fmt.Sprintf("%x", h)
}

func (s *V4Signer) derivedKey(t time.Time) []byte {
	h := s.hmac([]byte("AWS4"+s.auth.SecretAccessKey), []byte(t.Format(ISO8601BasicFormatShort)))
	h = s.hmac(h, []byte(s.regionName))
	h = s.hmac(h, []byte(s.serviceName))
	h = s.hmac(h, []byte("aws4_request"))
	return h
}
func (s *V4Signer) authorization(header http.Header, t time.Time, signature string) string {
	w := new(bytes.Buffer)
	fmt.Fprint(w, "AWS4-HMAC-SHA256 ")
	fmt.Fprintf(w, "Credential=%s/%s, ", s.auth.AccessKeyID, s.credentialScope(t))
	fmt.Fprintf(w, "SignedHeaders=%s, ", s.signedHeaders(header))
	fmt.Fprintf(w, "Signature=%s", signature)
	return w.String()
}

func (s *V4Signer) hash(in string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s", in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *V4Signer) hmac(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
