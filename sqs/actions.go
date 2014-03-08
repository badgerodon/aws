package sqs

import (
	"github.com/badgerodon/aws"
	"net/http"
	"net/url"
)

type (
	Parameter struct {
		Key, Value string
	}
	Client struct {
		AccessKeyId, SecretAccessKey string
		Region                       string
	}
)

func (this *Client) Get(endpoint, params []Parameter) error {

}

func Get(endpoint, params []Parameter) error {
	urlStr := endpoint + "?Action=" + url.QueryEscape(action)
	for _, param := range params {
		urlStr += "&" + url.QueryEscape(param.Key) + "=" + url.QueryEscape(param.Value)
	}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}
	//v4.Sign(req)
	return nil
}
