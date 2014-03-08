package aws

import (
	"net/http"
	"net/url"
	"testing"
)

func TestSignV2(t *testing.T) {
	var req http.Request
	req.Header = make(http.Header)
	req.URL, _ = url.Parse("https://elasticmapreduce.amazonaws.com?Action=DescribeJobFlows&Version=2009-03-31&AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&SignatureVersion=2&SignatureMethod=HmacSHA256&Timestamp=2011-10-03T15%3A19%3A30")
	NewV2Signer(
		Auth{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"},
		"elasticmapreduce",
	).Sign(&req)
	t.Errorf("%v", req)
}
