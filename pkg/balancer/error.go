package balancer

import "strings"

func handlerSenderError(err error) bool {
	switch {
	case strings.Contains(err.Error(), "CHAT_WRITE_FORBIDDEN"):
		return true
	//case strings.Contains(err.Error(), "Forbidden: bot can't initiate conversation with a user"):
	//	return true
	default:
		return false
	}
}
