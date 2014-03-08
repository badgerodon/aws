package v4

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"net/url"
	"testing"
)

func TestGetSignatureKey(t *testing.T) {
	key := "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"
	dateStamp := "20120215"
	regionName := "us-east-1"
	serviceName := "iam"

	result := hex.EncodeToString(getSignatureKey(key, dateStamp, regionName, serviceName))
	expected := "f4780e2d9f65fa895f9c67b32ce1baf0b0d8a43505a000a1a9e090d414db404d"

	if expected != result {
		t.Errorf("invalid signature. expected %v, got %v", expected, result)
	}
}

func TestCanonicalHeaders(t *testing.T) {
	var req http.Request
	req.Header = make(http.Header)
	req.Header.Set("host", "iam.amazonaws.com")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("My-header1", "    a   b   c ")
	req.Header.Set("x-amz-date", "20120228T024136Z")
	req.Header.Set("My-Header2", `    "a   b   c"`)

	result := canonicalHeaders(&req)
	expected := "content-type:application/x-www-form-urlencoded; charset=utf-8\nhost:iam.amazonaws.com\nmy-header1:a b c\nmy-header2:\"a   b   c\"\nx-amz-date:20120228T024136Z\n"

	if result != expected {
		t.Errorf("invalid canonical headers. expected \n%s, got \n%s.", expected, result)
	}
}

func TestSignRequest(t *testing.T) {
	var req http.Request
	req.Header = make(http.Header)
	req.Method = "POST"
	req.URL, _ = url.Parse("http://iam.amazonaws.com")
	//Authorization: AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=ced6826de92d2bdeed8f846f0bf508e8559e98e4b0199114b84c54174deb456c
	req.Host = "iam.amazonaws.com"
	req.Header.Set("host", "iam.amazonaws.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("x-amz-date", "20110909T233600Z")
	var buf bytes.Buffer
	buf.WriteString("Action=ListUsers&Version=2010-05-08")
	req.Body = closableBuffer{&buf}

	err := Sign(&req, "AKIDEXAMPLE", "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY", "us-east-1", "iam")
	if err != nil {
		t.Errorf("expected successful signing, got %v", err)
	}
	result := req.Header.Get("Authorization")
	expected := "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=ced6826de92d2bdeed8f846f0bf508e8559e98e4b0199114b84c54174deb456c"
	if result != expected {
		t.Errorf("invalid signing, expected \n%s, got \n%s", result, expected)
	}

}
