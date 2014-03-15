package sqs

type (
	ResponseMetadata struct {
		RequestId string
		BoxUsage  float64
	}
	AddPermissionResponse struct {
		ResponseMetadata ResponseMetadata
	}
	ChangeMessageVisiblityResponse struct {
		ResponseMetadata ResponseMetadata
	}
	ChangeMessageVisiblityBatchResponse struct {
		Successful       []string                `xml:"ChangeMessageVisibilityBatchResult>ChangeMessageVisibilityBatchResultEntry>Id"`
		Failed           []BatchResultErrorEntry `xml:"ChangeMessageVisibilityBatchResult>BatchResultErrorEntry"`
		ResponseMetadata ResponseMetadata
	}
	BatchResultErrorEntry struct {
		Code        string
		Id          string
		Message     string
		SenderFault bool
	}
	ChangeMessageVisibilityBatchRequest struct {
		Id                string
		ReceiptHandle     string
		VisibilityTimeout int
	}
	CreateQueueResponse struct {
		QueueUrl         string `xml:"CreateQueueResult>QueueUrl"`
		ResponseMetadata ResponseMetadata
	}
	DeleteMessageBatchRequestEntry struct {
		Id            string
		ReceiptHandle string
	}
	DeleteMessageBatchResult struct {
		Failed     []BatchResultErrorEntry
		Successful []DeleteMessageBatchResultEntry
	}
	DeleteMessageBatchResultEntry struct {
		Id string
	}
	ErrorResponse struct {
		Type      string `xml:"Error>Type"`
		Code      string `xml:"Error>Code"`
		Message   string `xml:"Error>Message"`
		RequestId string
	}
	GetQueueAttributesResult struct {
		Attributes map[string]string
	}
	GetQueueUrlResult struct {
		QueueUrl string
	}
	ListDeadLetterSourceQueuesResult struct {
		QueueUrls []string
	}
	ListQueuesResponse struct {
		QueueUrls        []string `xml:"ListQueuesResult>QueueUrl"`
		ResponseMetadata ResponseMetadata
	}
	Message struct {
		Attributes    map[string]string
		Body          string
		MD5OfBody     string
		MessageId     string
		ReceiptHandle string
	}
	ReceiveMessageResult struct {
		Messages []Message
	}
	SendMessageBatchRequestEntry struct {
		DelaySeconds int
		Id           string
		MessageBody  string
	}
	SendMessageBatchResult struct {
		Failed     []BatchResultErrorEntry
		Successful SendMessageBatchResultEntry
	}
	SendMessageBatchResultEntry struct {
		Id               string
		MD5OfMessageBody string
		MessageId        string
	}
	SendMessageResult struct {
		MD5OfMessageBody string
		MessageId        string
	}
)
