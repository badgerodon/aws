package v2

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func Sign(req *http.Request) error {
	ps := []string{}
	for k, vs := range req.URL.Query() {
		for _, v := range vs {
			ps = append(ps, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}
	sort.Strings(ps)

	stringToSign := req.Method + "\n" +
		req.URL.Host + "\n" +
		req.URL.Path + "\n" +
		strings.Join(ps, "&")

}
