package aws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// CloudWatchValue stores a single value to push to CloudWatch.
type CloudWatchValue struct {
	MetricName string
	Unit       string
	Value      string
}

// CloudwatchReporter sends metrics to Cloudwatch using the AWS API.
type CloudWatchReporter struct {
	auth      Auth
	namespace string
	debug     bool
}

// NewCloudWatchReporter creates and returns a new CloudWatchReporter.
func NewCloudWatchReporter(namespace string, debug bool) (cw CloudWatchReporter, err error) {
	// Authenticate using the environment variables.
	auth, err := EnvAuth()
	if err != nil {
		return
	}

	// Return a new CloudWatchReporter which can be utilized to report metrics to Cloudwatch.
	cw.auth = auth
	cw.namespace = namespace
	cw.debug = debug

	return
}

// Report reports the data to AWS Cloudwatch.
func (cwr *CloudWatchReporter) Report(values ...CloudWatchValue) error {
	req, err := http.NewRequest("PUT", CloudwatchURL, nil)
	if err != nil {
		return err
	}

	req.ParseForm()
	req.Form.Set("Action", "PutMetricData")
	req.Form.Set("Namespace", cwr.namespace)

	for i, value := range values {
		metricFormat := fmt.Sprintf("MetricData.member.%d.", i+1)

		req.Form.Set(metricFormat+"MetricName", value.MetricName)
		req.Form.Set(metricFormat+"Unit", value.Unit)
		req.Form.Set(metricFormat+"Value", value.Value)
	}

	return cwr.Push(req)
}

// Push sends the reported metrics to Cloudwatch.
func (cwr *CloudWatchReporter) Push(req *http.Request) (err error) {
	client := &http.Client{}

	req.URL.RawQuery = req.Form.Encode()
	cwr.auth.Sign(req)

	// Print debug info, if debug is true.
	if cwr.debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return err
		}

		log.Println(string(dump))
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// Print debug info, if debug is true.
	if cwr.debug || resp.StatusCode != 200 {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			log.Println(string(dump))
		} else {
			return errors.New(
				fmt.Sprintf(
					"Error while reporting metrics!\nStatus code: %d\nResponse body: %s",
					resp.StatusCode,
					string(dump),
				),
			)
		}
	}

	return
}
