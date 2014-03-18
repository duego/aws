// Package aws handles signing of API requests to Amazon
// Partly stolen from launchpad.net/goamz/aws
package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"time"
)

const debug = false

var b64 = base64.StdEncoding

type Auth struct {
	AccessKey, SecretKey string
}

// EnvAuth creates an Auth based on environment information.
func EnvAuth() (auth Auth, err error) {
	auth.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	auth.SecretKey = os.Getenv("AWS_SECRET_KEY")

	if auth.AccessKey == "" {
		err = errors.New("AWS_ACCESS_KEY not found in environment")
	}
	if auth.SecretKey == "" {
		err = errors.New("AWS_SECRET_KEY not found in environment")
	}
	return
}

// Sign the request to be sent to Amazon.
func (a *Auth) Sign(req *http.Request) {
	req.ParseForm()

	req.Form.Set("Timestamp", time.Now().In(time.UTC).Format(time.RFC3339))
	req.Form.Set("AWSAccessKeyId", a.AccessKey)
	req.Form.Set("SignatureVersion", "2")
	req.Form.Set("SignatureMethod", "HmacSHA256")
	req.URL.RawQuery = req.Form.Encode()

	if req.URL.Path == "" {
		req.URL.Path = "/"
	}

	payload := req.Method + "\n" + req.URL.Host + "\n" + req.URL.Path + "\n" + req.URL.RawQuery
	hash := hmac.New(sha256.New, []byte(a.SecretKey))
	hash.Write([]byte(payload))
	signature := make([]byte, b64.EncodedLen(hash.Size()))
	b64.Encode(signature, hash.Sum(nil))

	req.Form.Set("Signature", string(signature))
	req.URL.RawQuery = req.Form.Encode()
}
