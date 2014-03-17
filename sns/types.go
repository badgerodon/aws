package sns

type (
	ResponseMetadata struct {
		RequestId string
	}

	// Actions

	AddPermissionRequest struct {
		AWSAccountIds []string
		ActionNames   []string
		Label         string
		TopicArn      string
	}
	AddPermissionResponse struct {
		ResponseMetadata
	}

	ConfirmSubscriptionRequest struct {
		AuthenticateOnUnsubscribe bool
		Token                     string
		TopicArn                  string
	}
	ConfirmSubscriptionResponse struct {
		SubscriptionArn string `xml:"ConfirmSubscriptionResult>SubscriptionArn"`
		ResponseMetadata
	}

	CreateTopicRequest struct {
		Name string
	}
	CreateTopicResponse struct {
		TopicArn string `xml:"CreateTopicResult>TopicArn"`
		ResponseMetadata
	}
)
