package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h *Headers) parseLine(line string) error {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return fmt.Errorf("line '%s' can not be parsed as header, incorrect spacing", line)
	}
	if !strings.HasSuffix(parts[0], ":") {
		return fmt.Errorf("line '%s' can not be parsed as header, field name should end with ':'", line)
	}
	fieldName := parts[0][:len(parts[0])-1]
	fieldValue := parts[1]
	(*h)[fieldName] = fieldValue
	return nil
}
