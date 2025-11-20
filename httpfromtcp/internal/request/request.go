package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/requestline"
	"io"
	"strconv"
)

const bufferSize = 4096

const (
	RequestStateReadingRequestLine = iota
	RequestStateReadingHeaders
	RequestStateReadingBody
	RequestStateDone
)

const CRLFbytes = 2

type Request struct {
	state       int
	RequestLine requestline.RequestLine
	Headers     headers.Headers
	Body        []byte
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state:       RequestStateReadingRequestLine,
		RequestLine: requestline.NewRequestLine(),
		Headers:     headers.NewHeaders(),
		Body:        make([]byte, 0),
	}
	buffer := make([]byte, bufferSize)
	parseTillIndex := 0

	for request.state != RequestStateDone {
		numOfBytesRead, errRead := reader.Read(buffer[parseTillIndex:])

		parseTillIndex += numOfBytesRead

		if parseTillIndex > bufferSize-1 {
			return &Request{}, fmt.Errorf("failed to process request: exceeded buffer size of %d", bufferSize)
		}

		isLastChunk := errRead == io.EOF

		numOfBytesParsed, errParse := request.parse(buffer[:parseTillIndex], isLastChunk)

		if errParse != nil {
			return &Request{}, fmt.Errorf("failed to process request: %v", errParse)
		}

		if errRead != nil && errRead != io.EOF {
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

func (r *Request) parse(data []byte, isLastChunk bool) (int, error) {
	numOfBytesParsed := 0

	for r.state != RequestStateDone {
		var err error
		var bytesConsumed int

		switch r.state {
		case RequestStateReadingRequestLine:
			bytesConsumed, err = r.parseRequestLine(data[numOfBytesParsed:])
		case RequestStateReadingHeaders:
			bytesConsumed, err = r.parseHeader(data[numOfBytesParsed:])
		case RequestStateReadingBody:
			bytesConsumed, err = r.parseBody(data[numOfBytesParsed:], isLastChunk)
		}

		if err != nil {
			return 0, fmt.Errorf("failed to parse data chunk, %v", err)
		}

		numOfBytesParsed += bytesConsumed

		if bytesConsumed == 0 {
			break
		}
	}

	return numOfBytesParsed, nil
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	lineEnd, hasCompleteLine := findNextCRLF(data, 0)
	if !hasCompleteLine {
		return 0, nil
	}
	line := string(data[0:(lineEnd - CRLFbytes)])
	err := r.RequestLine.ParseLine(line)
	if err != nil {
		return 0, err
	}
	r.state = RequestStateReadingHeaders
	return lineEnd, nil
}

func (r *Request) parseHeader(data []byte) (int, error) {
	lineEnd, hasCompleteLine := findNextCRLF(data, 0)
	if !hasCompleteLine {
		return 0, nil
	}
	line := string(data[0:(lineEnd - CRLFbytes)])
	// if length is 0 then we are reading the \r\n empty line which is the indicator of end of header
	if len(line) == 0 {
		r.state = RequestStateReadingBody
		return lineEnd, nil
	}
	err := r.Headers.ParseLine(line)
	if err != nil {
		return 0, err
	}
	return lineEnd, nil
}

func (r *Request) parseBody(data []byte, isLastChunk bool) (int, error) {
	r.Body = append(r.Body, data...)
	if isLastChunk {
		r.state = RequestStateDone
	}
	headerContentLength, ok := r.Headers.Get("Content-Length")
	if !ok {
		return len(data), nil
	}
	contentLength, err := strconv.Atoi(headerContentLength)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Content-Length '%s'", headerContentLength)
	}
	if isLastChunk && len(r.Body) != contentLength {
		return 0, fmt.Errorf("received %d bytes of body when expected %d", len(r.Body), contentLength)
	}
	return len(data), nil
}

func findNextCRLF(data []byte, start int) (lineEnd int, hasCompleteLine bool) {
	i := bytes.Index(data[start:], []byte("\r\n"))
	if i == -1 {
		return 0, false
	}
	return start + i + CRLFbytes, true
}
