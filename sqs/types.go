package sqs

type (
	BatchResultErrorEntry struct {
		Code        string
		Id          string
		Message     string
		SenderFault bool
	}
	ChangeMessageVisibilityBatchRequestEntry struct {
		Id                string
		ReceiptHandle     string
		VisibilityTimeout int
	}
	ChangeMessageVisibilityBatchResult struct {
		Failed     []BatchResultErrorEntry
		Successful []ChangeMessageVisibilityBatchResultEntry
	}
	ChangeMessageVisibilityBatchResultEntry struct {
		Id string
	}
	CreateQueueResult struct {
		QueueUrl string
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
	GetQueueAttributesResult struct {
		Attributes map[string]string
	}
	GetQueueUrlResult struct {
		QueueUrl string
	}
	ListDeadLetterSourceQueuesResult struct {
		QueueUrls []string
	}
	ListQueuesResult struct {
		QueueUrls []string
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
