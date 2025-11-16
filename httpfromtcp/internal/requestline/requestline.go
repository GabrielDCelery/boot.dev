package requestline

import (
	"fmt"
	"slices"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func NewRequestLine() RequestLine {
	return RequestLine{}
}

func (rl *RequestLine) ParseLine(line string) error {
	requestLineParts := strings.Split(line, " ")

	if len(requestLineParts) != 3 {
		return fmt.Errorf("invalid request line '%s', should have three parts", line)
	}

	method := requestLineParts[0]
	requestTarget := requestLineParts[1]
	httpVersionRaw := requestLineParts[2]

	if err := validateMethod(method); err != nil {
		return fmt.Errorf("invalid request: %v", err)
	}

	httpVersion, err := validateHttpVersion(httpVersionRaw)

	if err != nil {
		return fmt.Errorf("invalid http version: %v", err)
	}

	if err = validateRequestTarget(requestTarget); err != nil {
		return fmt.Errorf("invalid request target: %v", err)
	}

	rl.HttpVersion = httpVersion
	rl.RequestTarget = requestTarget
	rl.Method = method

	return nil
}

func validateMethod(method string) error {
	validMethods := []string{"GET", "POST", "PUT", "DELETE"}
	if slices.Contains(validMethods, method) {
		return nil
	}
	return fmt.Errorf("invalid method, received: '%s', valid values are: %v", method, validMethods)
}

func validateHttpVersion(httpVersion string) (string, error) {
	validHttpVersions := []string{"HTTP/1.1"}
	if slices.Contains(validHttpVersions, httpVersion) {
		return strings.Replace(httpVersion, "HTTP/", "", 1), nil
	}
	return "", fmt.Errorf("invalid http version, received: '%s', valid values are: %v", httpVersion, validHttpVersions)
}

func validateRequestTarget(target string) error {
	if target == "" {
		return fmt.Errorf("request target can not be empty")
	}
	if !strings.HasPrefix(target, "/") && !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		return fmt.Errorf("invalid request target '%s', must start with '/', 'http://' or 'https://'", target)
	}
	return nil
}
