package commands

import (
	"github.com/Leixb/jutge/auth"

	"github.com/imroc/req"
)

var creds *auth.Credentials

func getCredentials(credentials ...string) *auth.Credentials {
	if creds == nil {
		creds = auth.GetInstance(credentials...)
	}
	return creds
}

func getToken(credentials ...string) string {
	return getCredentials(credentials...).TokenUID
}

func getReq(credentials ...string) *req.Req {
	return getCredentials(credentials...).R
}
