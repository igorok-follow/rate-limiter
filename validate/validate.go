package validate

import (
	"errors"
	"net"
)

var (
	errorEmptyHandlerName = errors.New("empty handler name")
	errorZeroHandlerTTL   = errors.New("error zero handler ttl")
	errorInvalidIP        = errors.New("invalid ip in white address")
)

func ValidateHandlers(handlers map[string]int) error {
	for k, v := range handlers {
		switch {
		case len(k) == 0:
			return errorEmptyHandlerName
		case v == 0:
			return errorZeroHandlerTTL
		}
	}

	return nil
}

func ValidateWhiteList(whitelist map[string]struct{}) error {
	for k := range whitelist {
		switch {
		case net.ParseIP(k) == nil:
			return errorInvalidIP
		}
	}

	return nil
}
