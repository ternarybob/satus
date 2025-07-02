package satus

import "strings"

const (
	CONFIG_AUTHENTICATION_URL = "auth_url"
	CONFIG_PRIVATEADR         = "private_addr"
)

func contains(a string, list []string) bool {

	for _, b := range list {
		if strings.EqualFold(a, b) {
			return true
		}
	}
	return false

}
