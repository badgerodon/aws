package aws

import (
	"net/http"
	"net/url"
	"testing"
)

func TestV2SignerBuildsStringToSign(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://elasticmapreduce.amazonaws.com?Action=DescribeJobFlows&Version=2009-03-31&AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&SignatureVersion=2&SignatureMethod=HmacSHA256&Timestamp=2011-10-03T15%3A19%3A30", nil)
	result := NewV2Signer(
		Auth{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"},
	).stringToSign(req)
	expected := "GET\nelasticmapreduce.amazonaws.com\n/\nAWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&Action=DescribeJobFlows&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=2011-10-03T15%3A19%3A30&Version=2009-03-31"
	if expected != result {
		t.Errorf("error creating string to sign, expected: \n%s, got: \n%s", expected, result)
	}
}

func TestV2SignerSigns(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://elasticmapreduce.amazonaws.com?Action=DescribeJobFlows&Version=2009-03-31&AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&SignatureVersion=2&SignatureMethod=HmacSHA256&Timestamp=2011-10-03T15%3A19%3A30", nil)
	NewV2Signer(
		Auth{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"},
	).Sign(req)
	result := req.URL.Query().Get("Signature")
	expected, _ := url.QueryUnescape("i91nKc4PWAt0JJIdXwz9HxZCJDdiy6cf%2FMj6vPxyYIs%3D")
	if result != expected {
		t.Errorf("error signing: expected \n%s, got \n%s", expected, result)
	}
}
