package s3

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	ListBucketRequest struct {
		Bucket       string
		Delimiter    string
		EncodingType string
		Marker       string
		MaxKeys      int
		Prefix       string
	}
	ListBucketResponse struct {
		Name           string            `xml:"Name"`
		Prefix         string            `xml:"Prefix"`
		Delimiter      string            `xml:"Delimiter"`
		Marker         string            `xml:"Marker"`
		NextMarker     string            `xml:"NextMarker"`
		MaxKeys        int               `xml:"MaxKeys"`
		IsTruncated    bool              `xml:"IsTruncated"`
		Contents       []ListBucketEntry `xml:"Contents"`
		CommonPrefixes []string          `xml:"CommonPrefixes>Prefix"`
	}
	ListBucketEntry struct {
		Key          string    `xml:"Key"`
		LastModified time.Time `xml:"LastModified"`
		ETag         string    `xml:"ETag"`
		Size         int       `xml:"Size"`
		StorageClass string    `xml:"StorageClass"`
		Owner        struct {
			ID          string `xml:"ID"`
			DisplayName string `xml:"DisplayName"`
		} `xml:"Owner"`
	}
)

func (req ListBucketRequest) Encode(h *http.Request) {
	q := h.URL.Query()
	if req.Delimiter != "" {
		q.Add("delimiter", req.Delimiter)
	}
	if req.EncodingType != "" {
		q.Add("encoding-type", req.EncodingType)
	}
	if req.Marker != "" {
		q.Add("marker", req.Marker)
	}
	if req.MaxKeys != 0 {
		q.Add("max-keys", fmt.Sprint(req.MaxKeys))
	}
	if req.Prefix != "" {
		q.Add("prefix", req.Prefix)
	}
	h.URL.Path = "/" + req.Bucket
	h.URL.RawQuery = q.Encode()
}

func (res *ListBucketResponse) Decode(h *http.Response) error {
	if h.StatusCode != 200 {
		return fmt.Errorf("%v", h)
	}
	bs, err := ioutil.ReadAll(h.Body)
	if err != nil {
		return err
	}
	h.Body.Close()
	err = xml.Unmarshal(bs, res)
	if err != nil {
		return err
	}
	return nil
}
