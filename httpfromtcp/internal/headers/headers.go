package headers

import (
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) parseLine(line string) error {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return fmt.Errorf("line '%s' can not be parsed as header, incorrect spacing", line)
	}
	if !strings.HasSuffix(parts[0], ":") {
		return fmt.Errorf("line '%s' can not be parsed as header, field name should end with ':'", line)
	}
	fieldName := parts[0][:len(parts[0])-1]
	err := validateFieldName(fieldName)
	if err != nil {
		return err
	}
	fieldName = convertFieldNameToConanocalForm(fieldName)
	fieldValue := parts[1]
	value, ok := h[fieldName]
	if ok {
		h[fieldName] = fmt.Sprintf("%s, %s", value, fieldValue)
	} else {
		h[fieldName] = fieldValue
	}
	return nil
}

func validateFieldName(fieldName string) error {
	pattern := `^[a-zA-Z0-9\!\#\$\%\&\'\*\+\-\.\^\_\|\~]+$`
	matched, err := regexp.Match(pattern, []byte(fieldName))
	if err != nil {
		return fmt.Errorf("failed to validate field name '%s', reason: %v", fieldName, err)
	}
	if !matched {
		return fmt.Errorf("field name '%s' contains invalid characters", fieldName)
	}
	return nil
}

func convertFieldNameToConanocalForm(fieldName string) string {
	parts := strings.Split(fieldName, "-")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, "-")
}
