/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package customresponseheaders

import (
	"regexp"
	"strings"

	networking "k8s.io/api/networking/v1"
	"k8s.io/ingress-nginx/internal/ingress/annotations/parser"

	ing_errors "k8s.io/ingress-nginx/internal/ingress/errors"
	"k8s.io/ingress-nginx/internal/ingress/resolver"
)

var headerRegexp = regexp.MustCompile(`^[a-zA-Z\d\-_]+$`)

// Config returns the custom response headers for an Ingress rule
type Config struct {
	ResponseHeaders map[string]string `json:"custom-response-headers,omitempty"`
}

type customresponseheaders struct {
	r resolver.Resolver
}

// NewParser creates a new custom response headers annotation parser
func NewParser(r resolver.Resolver) parser.IngressAnnotation {
	return customresponseheaders{r}
}

// Parse parses the annotations contained in the ingress to use
// custom response headers
func (e customresponseheaders) Parse(ing *networking.Ingress) (interface{}, error) {
	headersMap := map[string]string{}
	responseHeader, err := parser.GetStringAnnotation("custom-response-headers", ing)
	if err != nil {
		return nil, err
	}

	headers := strings.Split(responseHeader, "||")
	for i := 0; i < len(headers); i++ {

		if !strings.Contains(headers[i], ":") {
			return nil, ing_errors.NewLocationDenied("Invalid header format")
		}

		headerSplit := strings.SplitN(headers[i], ":", 2)
		for j := range headerSplit {
			headerSplit[j] = strings.TrimSpace(headerSplit[j])
		}

		if len(headerSplit) < 2 {
			return nil, ing_errors.NewLocationDenied("Invalid header size")
		}

		if !ValidHeader(headerSplit[0]) {
			return nil, ing_errors.NewLocationDenied("Invalid header name")
		}

		headersMap[strings.TrimSpace(headerSplit[0])] = strings.TrimSpace(headerSplit[1])
	}
	return &Config{headersMap}, nil
}

// ValidHeader checks is the provided string satisfies the header's name regex
func ValidHeader(header string) bool {
	return headerRegexp.Match([]byte(header))
}
