package sns

import (
	"fmt"
)

func (c *Client) AddPermission(req AddPermissionRequest) (AddPermissionResponse, error) {
	params := map[string]string{
		"Action":   "AddPermission",
		"Label":    req.Label,
		"TopicArn": req.TopicArn,
	}
	for i, accountId := range req.AWSAccountIds {
		params[fmt.Sprint("AWSAccountId.", i+1)] = accountId
	}
	for i, actionName := range req.ActionNames {
		params[fmt.Sprint("ActionName.", i+1)] = actionName
	}
	var res AddPermissionResponse
	return res, c.Get(params, &res)
}

func (c *Client) CreateTopic(req CreateTopicRequest) (CreateTopicResponse, error) {
	params := map[string]string{
		"Action": "CreateTopic",
		"Name":   req.Name,
	}
	var res CreateTopicResponse
	return res, c.Get(params, &res)
}
