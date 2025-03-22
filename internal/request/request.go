package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type requestState int

const (
	stateInitialized requestState = iota // 0
	stateDone                            // 1
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		state: requestStateInitialized,
	}
	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(io.EOF, err) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request")
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateDone {
		return 0, nil
	}
	if !strings.Contains(string(data), "\r\n") {
		return 0, nil
	}
	// Split the request into lines
	lines := strings.Split(string(data), "\r\n")
	req := strings.Split(lines[0], " ")
	if len(req) != 3 {
		return 0, fmt.Errorf("invalid request line: %s", lines[0])
	}

	for _, l := range req[0] {
		if strings.ToUpper(string(l)) != string(l) {
			return 0, fmt.Errorf("invalid request line: %s", lines[0])
		}

	}
	if req[2] != "HTTP/1.1" {
		return 0, fmt.Errorf("invalid request line: %s", lines[0])
	}

	r.RequestLine.Method = req[0]
	r.RequestLine.RequestTarget = req[1]
	r.RequestLine.HttpVersion = strings.Split(req[2], "/")[1]
	r.state = stateDone
	return len(data), nil
}
