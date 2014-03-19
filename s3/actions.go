package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	GetObjectRequest struct {
		Bucket                     string
		Key                        string
		Range                      string
		IfModifiedSince            string
		IfUnmodifiedSince          string
		IfMatch                    string
		IfNoneMatch                string
		ResponseContentType        string
		ResponseContentLanguage    string
		ResponseExpires            string
		ResponseCacheControl       string
		ResponseContentDisposition string
		ResponseContentEncoding    string
	}
	GetObjectResponse struct {
		io.ReadCloser
		DeleteMarker            bool
		Expiration              string
		ServerSideEncryption    string
		Restore                 string
		VersionId               string
		WebsiteRedirectLocation string
	}
)

func (req GetObjectRequest) Encode(h *http.Request) {
	if req.Range != "" {
		h.Header.Set("Range", req.Range)
	}
	if req.IfModifiedSince != "" {
		h.Header.Set("If-Modified-Since", req.IfModifiedSince)
	}
	if req.IfUnmodifiedSince != "" {
		h.Header.Set("If-Unmodified-Since", req.IfUnmodifiedSince)
	}
	if req.IfMatch != "" {
		h.Header.Set("If-Match", req.IfMatch)
	}
	if req.IfNoneMatch != "" {
		h.Header.Set("If-None-Match", req.IfNoneMatch)
	}
	h.URL.Path = "/" + req.Bucket + "/" + req.Key
}

func (res *GetObjectResponse) Decode(h *http.Response) error {
	if h.StatusCode != 200 {
		if h.Body == nil {
			return fmt.Errorf("%s", h.Status)
		}
		bs, _ := ioutil.ReadAll(h.Body)
		h.Body.Close()
		return fmt.Errorf("%s", string(bs))
	}
	res.ReadCloser = h.Body
	return nil
}

func (c *Client) GetObject(req GetObjectRequest) (GetObjectResponse, error) {
	var res GetObjectResponse
	return res, c.Do(req, &res)
}

func (c *Client) ListBucket(req ListBucketRequest) (ListBucketResponse, error) {
	var res ListBucketResponse
	return res, c.Do(req, &res)
}
