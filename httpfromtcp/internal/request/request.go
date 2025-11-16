package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"slices"
	"strings"
)

const bufferSize = 4096

const (
	RequestStateReadingRequestLine = iota
	RequestStateReadingFieldLine
	RequestStateDone
)

type Request struct {
	state       int
	RequestLine RequestLine
	Headers     headers.Headers
}

func NewRequest() *Request {
	return &Request{
		state:       RequestStateReadingRequestLine,
		RequestLine: RequestLine{},
		Headers:     make(map[string]string),
	}
}

func (r *Request) Parse(data []byte) (int, error) {
	numOfBytesParsed := 0
	pointer := 0
	for {
		if pointer == len(data) {
			break
		}
		if r.state == RequestStateDone {
			break
		}
		if r.state == RequestStateReadingRequestLine {
			char := data[pointer]
			if char == '\r' {
				if pointer == len(data)-1 {
					pointer += 1
					continue
				}

				if data[pointer+1] == '\n' {
					requestLine, err := parseLineAsRequestLine(string(data[numOfBytesParsed:pointer]))
					if err != nil {
						return 0, err
					}
					r.RequestLine = requestLine
					r.state = RequestStateDone
					pointer += 2
					numOfBytesParsed += pointer
					continue
				}
			}
			pointer += 1
		}
	}
	return numOfBytesParsed, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()
	buffer := make([]byte, bufferSize)
	parseTillIndex := 0

	for request.state != RequestStateDone {
		numOfBytesRead, errRead := reader.Read(buffer[parseTillIndex:])

		parseTillIndex += numOfBytesRead

		if parseTillIndex > bufferSize-1 {
			return &Request{}, fmt.Errorf("failed to process request: exceeded buffer size of %d", bufferSize)
		}

		numOfBytesParsed, errParse := request.Parse(buffer[:parseTillIndex])

		if errParse != nil {
			return &Request{}, fmt.Errorf("failed to process request: %v", errParse)
		}

		if errRead != nil {
			if errRead == io.EOF {
				break
			}
			return &Request{}, fmt.Errorf("failed to process request: %v", errRead)
		}

		if numOfBytesParsed > 0 {
			copy(buffer, buffer[numOfBytesParsed:parseTillIndex])
			parseTillIndex -= numOfBytesParsed
		}
	}

	if request.state != RequestStateDone {
		return &Request{}, fmt.Errorf("incomplete HTTP request: reached EOF before request completed")
	}

	return request, nil
}

func parseLineAsRequestLine(line string) (RequestLine, error) {
	requestLineParts := strings.Split(line, " ")

	if len(requestLineParts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line '%s', should have three parts", line)
	}

	method := requestLineParts[0]
	requestTarget := requestLineParts[1]
	httpVersionRaw := requestLineParts[2]

	if err := validateMethod(method); err != nil {
		return RequestLine{}, fmt.Errorf("invalid request: %v", err)
	}

	httpVersion, err := validteHttpVersion(httpVersionRaw)

	if err != nil {
		return RequestLine{}, fmt.Errorf("invalid http version: %v", err)
	}

	if err = validateRequestTarget(requestTarget); err != nil {
		return RequestLine{}, fmt.Errorf("invalid request target: %v", err)
	}

	requestLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}

	return requestLine, nil

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
	if !strings.HasPrefix(target, "/") && !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		return fmt.Errorf("invalid request target '%s', must start with '/', 'http://' or 'https://'", target)
	}
	return nil
}
