package sns

import (
	"encoding/xml"
	"fmt"
	"github.com/badgerodon/aws"
	"io/ioutil"
	"net/http"
)

type (
	Client struct {
		Auth       aws.Auth
		RegionName string
	}
)

func New(auth aws.Auth, regionName string) *Client {
	return &Client{auth, regionName}
}

func (c *Client) fillDefaults(params Parameters) {
	if _, ok := params["AWSAccessKeyID"]; !ok {
		params["AWSAccessKeyID"] = c.Auth.AccessKeyID
	}
	if _, ok := params["Version"]; !ok {
		params["Version"] = "2010-03-31"
	}
}

func (c *Client) Get(params Parameters, dst interface{}) error {
	c.fillDefaults(params)

	url := "https://sns." + c.RegionName + ".amazonaws.com"
	i := 0
	for k, v := range params {
		if i == 0 {
			url += "?"
		} else {
			url += "&"
		}
		url += aws.Encode(k) + "=" + aws.Encode(v)
		i++
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	aws.NewV4Signer(c.Auth, "sns", c.RegionName).Sign(req)
	resp, err := aws.DefaultHTTPClient.Do(req)
	if err != nil {
		return err
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	//fmt.Println(string(bs))

	if resp.StatusCode != 200 {
		var res ErrorResponse
		if xml.Unmarshal(bs, &res) == nil {
			return fmt.Errorf("Code:%v,Type:%v,Message:%v", res.Code, res.Type, res.Message)
		} else {
			return fmt.Errorf("%v, %v", resp.StatusCode, string(bs))
		}
	}

	if dst != nil {
		return xml.Unmarshal(bs, dst)
	}
	return nil
}
