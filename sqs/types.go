package sqs

import (
	"encoding/xml"
	"fmt"
)

type (
	BatchResultErrorEntry struct {
		Code        string
		Id          string
		Message     string
		SenderFault bool
	}
	Message struct {
		Attributes    map[string]string
		Body          string
		MD5OfBody     string
		MessageId     string
		ReceiptHandle string
	}
	ResponseMetadata struct {
		RequestId string
		BoxUsage  float64
	}
	SendMessageBatchResultEntry struct {
		Id               string
		MessageId        string
		MD5OfMessageBody string
	}
	ErrorResponse struct {
		Type      string `xml:"Error>Type"`
		Code      string `xml:"Error>Code"`
		Message   string `xml:"Error>Message"`
		RequestId string
	}
)

// Requests
type (
	ChangeMessageVisibilityBatchRequest struct {
		Id                string
		ReceiptHandle     string
		VisibilityTimeout int
	}
	DeleteMessageBatchRequest struct {
		Id            string
		ReceiptHandle string
	}
	ReceiveMessageRequest struct {
		Attributes          []string
		MaxNumberOfMessages int
		VisibilityTimeout   int
		WaitTimeSeconds     int
	}
	SendMessageBatchRequest struct {
		Id           string
		MessageBody  string
		DelaySeconds int
	}
)

// Responses
type (
	AddPermissionResponse struct {
		ResponseMetadata ResponseMetadata
	}
	ChangeMessageVisibilityResponse struct {
		ResponseMetadata ResponseMetadata
	}
	ChangeMessageVisiblityBatchResponse struct {
		Successful       []string                `xml:"ChangeMessageVisibilityBatchResult>ChangeMessageVisibilityBatchResultEntry>Id"`
		Failed           []BatchResultErrorEntry `xml:"ChangeMessageVisibilityBatchResult>BatchResultErrorEntry"`
		ResponseMetadata ResponseMetadata
	}
	CreateQueueResponse struct {
		QueueUrl         string `xml:"CreateQueueResult>QueueUrl"`
		ResponseMetadata ResponseMetadata
	}
	DeleteMessageResponse struct {
		ResponseMetadata ResponseMetadata
	}
	DeleteMessageBatchResponse struct {
		Successful       []string                `xml:"DeleteMessageBatchResult>DeleteMessageBatchResultEntry>Id"`
		Failed           []BatchResultErrorEntry `xml:"DeleteMessageBatchResult>BatchResultErrorEntry"`
		ResponseMetadata ResponseMetadata
	}
	DeleteQueueResponse struct {
		ResponseMetadata ResponseMetadata
	}
	GetQueueAttributesResponse struct {
		Attributes       map[string]string
		ResponseMetadata ResponseMetadata
	}
	GetQueueUrlResponse struct {
		QueueUrl         string `xml:"GetQueueUrlResult>QueueUrl"`
		ResponseMetadata ResponseMetadata
	}
	ListDeadLetterSourceQueuesResponse struct {
		QueueUrls        []string `xml:"ListDeadLetterSourceQueuesResult>QueueUrl"`
		ResponseMetadata ResponseMetadata
	}
	ListQueuesResponse struct {
		QueueUrls        []string `xml:"ListQueuesResult>QueueUrl"`
		ResponseMetadata ResponseMetadata
	}
	ReceiveMessageResponse struct {
		Messages         []Message `xml:"ReceiveMessageResult>Message"`
		ResponseMetadata ResponseMetadata
	}
	RemovePermissionResponse struct {
		ResponseMetadata ResponseMetadata
	}
	SendMessageResponse struct {
		MD5OfMessageBody string `xml:"SendMessageResult>MD5OfMessageBody"`
		MessageId        string `xml:"SendMessageResult>MessageId"`
		ResponseMetadata ResponseMetadata
	}
	SendMessageBatchResponse struct {
		Successful       []SendMessageBatchResultEntry `xml:"SendMessageBatchResult>SendMessageBatchResultEntry>Id"`
		Failed           []BatchResultErrorEntry       `xml:"SendMessageBatchResult>BatchResultErrorEntry"`
		ResponseMetadata ResponseMetadata
	}
	SetQueueAttributesResponse struct {
		ResponseMetadata ResponseMetadata
	}
)

// custom xml unmarshaling
func (this *GetQueueAttributesResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local != "GetQueueAttributesResponse" {
		return fmt.Errorf("unable to unmarshal %s into GetQueueAttributesResponse", start.Name.Local)
	}
	type Temp struct {
		GetQueueAttributesResult struct {
			Attributes []struct {
				Name, Value string
			} `xml:"Attribute"`
		}
		ResponseMetadata ResponseMetadata
	}
	var tmp Temp
	err := d.DecodeElement(&tmp, &start)
	if err != nil {
		return err
	}
	this.Attributes = make(map[string]string)
	for _, attr := range tmp.GetQueueAttributesResult.Attributes {
		this.Attributes[attr.Name] = attr.Value
	}
	this.ResponseMetadata = tmp.ResponseMetadata
	return nil
}
