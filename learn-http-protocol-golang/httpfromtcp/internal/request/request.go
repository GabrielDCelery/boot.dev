package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	lines := strings.Split(string(b), "\r\n")

	if len(lines) == 0 {
		return &Request{}, fmt.Errorf("expected request to have request line")
	}
	requestLineParts := strings.Split(lines[0], " ")

	if len(requestLineParts) != 3 {
		return &Request{}, fmt.Errorf("invalid request line, should have three parts")
	}

	method := requestLineParts[0]
	requestTarget := requestLineParts[1]
	httpVersionRaw := requestLineParts[2]

	if err = validateMethod(method); err != nil {
		return &Request{}, fmt.Errorf("invalid request: %v", err)
	}

	httpVersion, err := validteHttpVersion(httpVersionRaw)

	if err != nil {
		return &Request{}, fmt.Errorf("invalid http version: %v", err)
	}

	requestLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}

	request := &Request{RequestLine: requestLine}

	return request, nil
}

func validateMethod(method string) error {
	validMethods := []string{"GET", "POST", "PUT", "DELETE"}
	if slices.Contains(validMethods, method) {
		return nil
	}
	return fmt.Errorf("invalid method, received: '%s', valid values are: %v", method, validMethods)
}

func validteHttpVersion(httpVersion string) (string, error) {
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
	if !strings.HasPrefix(target, "/") || !strings.HasPrefix(target, "http://") || !strings.HasPrefix(target, "https://") {
		return fmt.Errorf("invalid request target '%s', must start with '/', 'http://' or 'https://'", target)
	}
	return nil
}
