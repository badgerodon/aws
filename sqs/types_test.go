package sqs

import (
	"encoding/xml"
	"testing"
)

func TestUnmarshaling(t *testing.T) {
	src := `<GetQueueAttributesResponse>
  <GetQueueAttributesResult>
    <Attribute>
      <Name>VisibilityTimeout</Name>
      <Value>30</Value>
    </Attribute>
    <Attribute>
      <Name>DelaySeconds</Name>
      <Value>0</Value>
    </Attribute>
    <Attribute>
      <Name>ReceiveMessageWaitTimeSeconds</Name>
      <Value>2</Value>
    </Attribute>
  </GetQueueAttributesResult>
  <ResponseMetadata>
    <RequestId>6fde8d1e-52cd-4581-8cd9-c512f4c64223</RequestId>
  </ResponseMetadata>
</GetQueueAttributesResponse>`
	var dst GetQueueAttributesResponse
	err := xml.Unmarshal([]byte(src), &dst)
	if err != nil {
		t.Errorf("error unmarshaling GetQueueAttributesResponse: %v", err)
	}
	if dst.ResponseMetadata.RequestId != "6fde8d1e-52cd-4581-8cd9-c512f4c64223" {
		t.Errorf("error unmarshaling GetQueueAttributesResponse. expected ResponseMetadata of 6fde8d1e-52cd-4581-8cd9-c512f4c64223, got: %v", dst.ResponseMetadata)
	}
	expected := "2"
	result := dst.Attributes["ReceiveMessageWaitTimeSeconds"]
	if expected != result {
		t.Errorf("error unmarshaling GetQueueAttributesResponse. expected ReceiveMessageWaitTimeSeconds of %v, got %v", expected, result)
	}
}
