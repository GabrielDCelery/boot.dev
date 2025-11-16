package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/requestline"
	"io"
)

const bufferSize = 4096

const (
	StateReadingRequestLine = iota
	StateReadingHeaders
	StateDone
)

const CRLFbytes = 2

type Request struct {
	state       int
	RequestLine requestline.RequestLine
	Headers     headers.Headers
}

func NewRequest() *Request {
	return &Request{
		state:       StateReadingRequestLine,
		RequestLine: requestline.NewRequestLine(),
		Headers:     headers.NewHeaders(),
	}
}

func (r *Request) Parse(data []byte) (int, error) {
	numOfBytesParsed := 0
	for {
		if r.state == StateDone {
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
	if r.state == StateReadingRequestLine {
		err := r.RequestLine.ParseLine(line)
		if err != nil {
			return err
		}
		r.state = StateReadingHeaders
		return nil
	}
	if r.state == StateReadingHeaders {
		if len(line) == 0 {
			r.state = StateDone
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

	for request.state != StateDone {
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

	if request.state != StateDone {
		return &Request{}, fmt.Errorf("incomplete HTTP request: reached EOF before request completed, request %+v", request)
	}

	return request, nil
}
