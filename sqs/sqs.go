package sqs

import (
	"encoding/xml"
	"fmt"
	"github.com/badgerodon/aws"
	"io/ioutil"
	"net/http"
)

type (
	Parameters map[string]string
	Client     struct {
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
		params["Version"] = "2009-02-01"
	}
}

func (c *Client) Get(url string, params Parameters, dst interface{}) error {
	c.fillDefaults(params)

	if url == "" {
		endpoint := "sqs." + c.RegionName + ".amazonaws.com"
		url = "https://" + endpoint + "/"
	}
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
	aws.NewV4Signer(c.Auth, "sqs", c.RegionName).Sign(req)
	resp, err := http.DefaultClient.Do(req)
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

func (c *Client) AddPermission(url string, label string, userToAction map[string]string) (AddPermissionResponse, error) {
	params := map[string]string{
		"Action": "AddPermission",
		"Label":  label,
	}
	i := 1
	for u, a := range userToAction {
		params[fmt.Sprint("AWSAccountId.", i)] = u
		params[fmt.Sprint("ActionName.", i)] = a
		i++
	}
	var res AddPermissionResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) ChangeMessageVisibility(url string, receiptHandle string, visibilityTimeout int) (ChangeMessageVisibilityResponse, error) {
	params := map[string]string{
		"Action":            "ChangeMessageVisibility",
		"ReceiptHandle":     receiptHandle,
		"VisibilityTimeout": fmt.Sprint(visibilityTimeout),
	}
	var res ChangeMessageVisibilityResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) ChangeMessageVisibilityBatch(url string, requests []ChangeMessageVisibilityBatchRequest) (ChangeMessageVisiblityBatchResponse, error) {
	params := map[string]string{
		"Action": "ChangeMessageVisibilityBatch",
	}
	i := 1
	for _, req := range requests {
		params[fmt.Sprint("ChangeMessageVisibilityBatchRequestEntry.", i, ".Id")] = req.Id
		params[fmt.Sprint("ChangeMessageVisibilityBatchRequestEntry.", i, ".ReceiptHandle")] = req.ReceiptHandle
		params[fmt.Sprint("ChangeMessageVisibilityBatchRequestEntry.", i, ".VisibilityTimeout")] = fmt.Sprint(req.VisibilityTimeout)
		i++
	}
	var res ChangeMessageVisiblityBatchResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) CreateQueue(name string, attributes map[string]string) (CreateQueueResponse, error) {
	params := map[string]string{
		"Action":    "CreateQueue",
		"QueueName": name,
	}
	i := 1
	for k, v := range attributes {
		params[fmt.Sprint("Attributes.", i, ".Name")] = k
		params[fmt.Sprint("Attributes.", i, ".Value")] = v
		i++
	}
	var res CreateQueueResponse
	return res, c.Get("", params, &res)
}

func (c *Client) SetQueueAttributes(url string, attributes map[string]string) error {
	params := map[string]string{
		"Action": "SetQueueAttributes",
	}
	i := 1
	for k, v := range attributes {
		params[fmt.Sprint("Attributes.", i, ".Name")] = k
		params[fmt.Sprint("Attributes.", i, ".Value")] = v
		i++
	}
	return c.Get(url, params, nil)
}
