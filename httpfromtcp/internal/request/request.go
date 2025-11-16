package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"slices"
	"strings"
)

const bufferSize = 4096

const (
	RequestStateReadingRequestLine = iota
	RequestStateReadingHeaders
	RequestStateDone
)

const CRLFbytes = 2

type Request struct {
	state       int
	RequestLine RequestLine
	Headers     headers.Headers
}

func NewRequest() *Request {
	return &Request{
		state:       RequestStateReadingRequestLine,
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
	}
}

func (r *Request) Parse(data []byte) (int, error) {
	numOfBytesParsed := 0
	for {
		if r.state == RequestStateDone {
			break
		}
		lineEnd, hasCompleteLine := findNextCRLF(data, numOfBytesParsed)
		if !hasCompleteLine {
			return numOfBytesParsed, nil
		}
		line := string(data[numOfBytesParsed:(lineEnd - CRLFbytes)])
		err := r.parseLine(line)
		if err != nil {
			return 0, err
		}
		numOfBytesParsed = lineEnd
	}
	return numOfBytesParsed, nil
}

func (r *Request) parseLine(line string) error {
	if r.state == RequestStateReadingRequestLine {
		err := r.RequestLine.ParseLine(line)
		if err != nil {
			return err
		}
		r.state = RequestStateReadingHeaders
		return nil
	}
	if r.state == RequestStateReadingHeaders {
		if len(line) == 0 {
			r.state = RequestStateDone
			return nil
		}
		err := r.Headers.ParseLine(line)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unhandled state")
}

func findNextCRLF(data []byte, start int) (lineEnd int, hasCompleteLine bool) {
	i := bytes.Index(data[start:], []byte("\r\n"))
	if i == -1 {
		return 0, false
	}
	return start + i + CRLFbytes, true
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
		return &Request{}, fmt.Errorf("incomplete HTTP request: reached EOF before request completed, request %+v", request)
	}

	return request, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
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
