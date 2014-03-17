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

func (c *Client) DeleteMessage(url string, receiptHandle string) (DeleteMessageResponse, error) {
	params := map[string]string{
		"Action":        "DeleteMessage",
		"ReceiptHandle": receiptHandle,
	}
	var res DeleteMessageResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) DeleteMessageBatch(url string, requests []DeleteMessageBatchRequest) (DeleteMessageBatchResponse, error) {
	params := map[string]string{
		"Action": "DeleteMessageBatch",
	}
	i := 1
	for _, req := range requests {
		params[fmt.Sprint("DeleteMessageBatchRequestEntry.", i, ".Id")] = req.Id
		params[fmt.Sprint("DeleteMessageBatchRequestEntry.", i, ".ReceiptHandle")] = req.ReceiptHandle
		i++
	}
	var res DeleteMessageBatchResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) DeleteQueue(url string) (DeleteQueueResponse, error) {
	params := map[string]string{
		"Action": "DeleteQueue",
	}
	var res DeleteQueueResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) GetQueueAttributes(url string, attributeNames []string) (GetQueueAttributesResponse, error) {
	params := map[string]string{
		"Action": "GetQueueAttributes",
	}
	i := 1
	for _, a := range attributeNames {
		params[fmt.Sprint("AttributeName.", i)] = a
		i++
	}
	var res GetQueueAttributesResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) GetQueueUrl(queueName string, queueOwnerAWSAccountId string) (GetQueueUrlResponse, error) {
	params := map[string]string{
		"Action":    "GetQueueUrl",
		"QueueName": queueName,
	}
	if queueOwnerAWSAccountId != "" {
		params["QueueOwnerAWSAccountId"] = queueOwnerAWSAccountId
	}
	var res GetQueueUrlResponse
	return res, c.Get("", params, &res)
}

func (c *Client) ListDeadLetterSourceQueues(url string) (ListDeadLetterSourceQueuesResponse, error) {
	params := map[string]string{
		"Action":   "ListDeadLetterSourceQueues",
		"QueueUrl": url,
	}
	var res ListDeadLetterSourceQueuesResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) ListQueues(queueNamePrefix string) (ListQueuesResponse, error) {
	params := map[string]string{
		"Action": "ListQueues",
	}
	if queueNamePrefix != "" {
		params["QueueNamePrefix"] = queueNamePrefix
	}
	var res ListQueuesResponse
	return res, c.Get("", params, &res)
}

func (c *Client) ReceiveMessage(url string, req ReceiveMessageRequest) (ReceiveMessageResponse, error) {
	params := map[string]string{
		"Action": "ReceiveMessage",
	}
	if req.MaxNumberOfMessages > 0 {
		params["MaxNumberOfMessages"] = fmt.Sprint(req.MaxNumberOfMessages)
	}
	if req.VisibilityTimeout > 0 {
		params["VisibilityTimeout"] = fmt.Sprint(req.VisibilityTimeout)
	}
	if req.WaitTimeSeconds > 0 {
		params["WaitTimeSeconds"] = fmt.Sprint(req.WaitTimeSeconds)
	}
	i := 1
	for _, a := range req.Attributes {
		params[fmt.Sprint("AttributeName.", i)] = a
		i++
	}
	var res ReceiveMessageResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) RemovePermission(url string, label string) (RemovePermissionResponse, error) {
	params := map[string]string{
		"Action": "RemovePermission",
		"Label":  label,
	}
	var res RemovePermissionResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) SendMessage(url string, body string, delaySeconds int) (SendMessageResponse, error) {
	params := map[string]string{
		"Action":      "SendMessage",
		"MessageBody": body,
	}
	if delaySeconds > 0 {
		params["DelaySeconds"] = fmt.Sprint(delaySeconds)
	}
	var res SendMessageResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) SendMessageBatch(url string, requests []SendMessageBatchRequest) (SendMessageBatchResponse, error) {
	params := map[string]string{
		"Action": "SendMessageBatch",
	}
	i := 1
	for _, req := range requests {
		params[fmt.Sprint("SendMessageBatchRequestEntry.", i, ".Id")] = req.Id
		params[fmt.Sprint("SendMessageBatchRequestEntry.", i, ".MessageBody")] = req.MessageBody
		params[fmt.Sprint("SendMessageBatchRequestEntry.", i, ".DelaySeconds")] = fmt.Sprint(req.DelaySeconds)
		i++
	}
	var res SendMessageBatchResponse
	return res, c.Get(url, params, &res)
}

func (c *Client) SetQueueAttributes(url string, attributes map[string]string) (SetQueueAttributesResponse, error) {
	params := map[string]string{
		"Action": "SetQueueAttributes",
	}
	i := 1
	for k, v := range attributes {
		params[fmt.Sprint("Attributes.", i, ".Name")] = k
		params[fmt.Sprint("Attributes.", i, ".Value")] = v
		i++
	}
	var res SetQueueAttributesResponse
	return res, c.Get(url, params, &res)
}
