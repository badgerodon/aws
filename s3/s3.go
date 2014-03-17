package s3

import (
	"github.com/badgerodon/aws"
	"net/http"
)

type (
	Client struct {
		Auth       aws.Auth
		RegionName string
	}
	Request interface {
		Encode(*http.Request)
	}
	Response interface {
		Decode(*http.Response) error
	}
)

func New(auth aws.Auth, regionName string) *Client {
	return &Client{auth, regionName}
}

func (c *Client) Do(req Request, res Response) error {
	subdomain := "s3"
	if c.RegionName != "us-east-1" {
		subdomain = "s3-" + c.RegionName
	}
	url := "https://" + subdomain + ".amazonaws.com"
	hreq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Encode(hreq)
	hreq.Header.Set("x-amz-content-sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	aws.NewV4Signer(c.Auth, "s3", c.RegionName).Sign(hreq)

	resp, err := aws.DefaultHTTPClient.Do(hreq)
	if err != nil {
		return err
	}

	return res.Decode(resp)
}
