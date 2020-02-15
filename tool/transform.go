package tool

import "errors"

// TransformToLogrusLevel ...
func TransformToLogrusLevel(configLevel string) (uint32, error) {
	switch configLevel {
	case "debug":
		return 5, nil
	case "info":
		return 4, nil
	case "error":
		return 2, nil
	default:
		return 5, errors.New("接受到规定外的日志等级: " + configLevel)
	}
}
