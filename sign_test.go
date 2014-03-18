package aws

import (
	"net/http"
	"net/http/httputil"
	"testing"
)

func TestSign(t *testing.T) {
	auth, err := EnvAuth()
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://monitoring.us-east-1.amazonaws.com/?Version=2010-08-01&Action=ListMetrics",
		nil,
	)
	auth.Sign(req)

	dump, _ := httputil.DumpRequest(req, true)
	t.Log(string(dump))

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	dump, _ = httputil.DumpResponse(resp, true)
	t.Log(string(dump))

	if resp.StatusCode != 200 {
		t.Fail()
	}
}
