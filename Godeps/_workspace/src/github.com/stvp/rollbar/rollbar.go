package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	NAME    = "go-rollbar"
	VERSION = "0.2.0"

	// Severity levels
	CRIT  = "critical"
	ERR   = "error"
	WARN  = "warning"
	INFO  = "info"
	DEBUG = "debug"

	FILTERED = "[FILTERED]"
)

var (
	// Rollbar access token. If this is blank, no errors will be reported to
	// Rollbar.
	Token = ""

	// All errors and messages will be submitted under this environment.
	Environment = "development"

	// API endpoint for Rollbar.
	Endpoint = "https://api.rollbar.com/api/1/item/"

	// Maximum number of errors allowed in the sending queue before we start
	// dropping new errors on the floor.
	Buffer = 1000

	// Filter GET and POST parameters from being sent to Rollbar.
	FilterFields = regexp.MustCompile("password|secret|token")

	// Queue of messages to be sent.
	bodyChannel chan map[string]interface{}
	waitGroup   sync.WaitGroup
)

type Field struct {
	Name string
	Data interface{}
}

// -- Setup

func init() {
	bodyChannel = make(chan map[string]interface{}, Buffer)

	go func() {
		for body := range bodyChannel {
			post(body)
			waitGroup.Done()
		}
	}()
}

// -- Error reporting

// Error asynchronously sends an error to Rollbar with the given severity level.
func Error(level string, err error, fields ...*Field) {
	ErrorWithStackSkip(level, err, 1, fields...)
}

// ErrorWithStackSkip asynchronously sends an error to Rollbar with the given
// severity level and a given number of stack trace frames skipped.
func ErrorWithStackSkip(level string, err error, skip int, fields ...*Field) {
	stack := BuildStack(2 + skip)
	ErrorWithStack(level, err, stack, fields...)
}

func ErrorWithStack(level string, err error, stack Stack, fields ...*Field) {
	buildAndPushError(level, err, stack, fields...)
}

// RequestError asynchronously sends an error to Rollbar with the given
// severity level and request-specific information.
func RequestError(level string, r *http.Request, err error, fields ...*Field) {
	RequestErrorWithStackSkip(level, r, err, 1, fields...)
}

// RequestErrorWithStackSkip asynchronously sends an error to Rollbar with the
// given severity level and a given number of stack trace frames skipped, in
// addition to extra request-specific information.
func RequestErrorWithStackSkip(level string, r *http.Request, err error, skip int, fields ...*Field) {
	stack := BuildStack(2 + skip)
	RequestErrorWithStack(level, r, err, stack, fields...)
}

// RequestErrorWithStack is like RequestError, but the stack is given
// as an argument
func RequestErrorWithStack(level string, r *http.Request, err error, stack Stack, fields ...*Field) {
	buildAndPushError(level, err, stack, &Field{Name: "request", Data: errorRequest(r)})
}

func buildError(level string, err error, stack Stack, fields ...*Field) map[string]interface{} {
	body := buildBody(level, err.Error())
	data := body["data"].(map[string]interface{})
	errBody, fingerprint := errorBody(err, stack)
	data["body"] = errBody
	data["fingerprint"] = fingerprint

	for _, field := range fields {
		data[field.Name] = field.Data
	}

	return body
}

func buildAndPushError(level string, err error, stack Stack, fields ...*Field) {
	push(buildError(level, err, stack, fields...))
}

// -- Message reporting

// Message asynchronously sends a message to Rollbar with the given severity
// level. Rollbar request is asynchronous.
func Message(level string, msg string) {
	body := buildBody(level, msg)
	data := body["data"].(map[string]interface{})
	data["body"] = messageBody(msg)

	push(body)
}

// -- Misc.

// Wait will block until the queue of errors / messages is empty.
func Wait() {
	waitGroup.Wait()
}

// Build the main JSON structure that will be sent to Rollbar with the
// appropriate metadata.
func buildBody(level, title string) map[string]interface{} {
	timestamp := time.Now().Unix()
	hostname, _ := os.Hostname()

	return map[string]interface{}{
		"access_token": Token,
		"data": map[string]interface{}{
			"environment": Environment,
			"title":       title,
			"level":       level,
			"timestamp":   timestamp,
			"platform":    runtime.GOOS,
			"language":    "go",
			"server": map[string]interface{}{
				"host": hostname,
			},
			"notifier": map[string]interface{}{
				"name":    NAME,
				"version": VERSION,
			},
		},
	}
}

// errorBody generate the error body with a given stack trace
func errorBody(err error, stack Stack) (map[string]interface{}, string) {
	fingerprint := stack.Fingerprint()
	errBody := map[string]interface{}{
		"trace": map[string]interface{}{
			"frames": stack,
			"exception": map[string]interface{}{
				"class":   errorClass(err),
				"message": err.Error(),
			},
		},
	}
	return errBody, fingerprint
}

// Extract error details from a Request to a format that Rollbar accepts.
func errorRequest(r *http.Request) map[string]interface{} {
	cleanQuery := filterParams(r.URL.Query())

	return map[string]interface{}{
		"url":     r.URL.String(),
		"method":  r.Method,
		"headers": flattenValues(r.Header),

		// GET params
		"query_string": url.Values(cleanQuery).Encode(),
		"GET":          flattenValues(cleanQuery),

		// POST / PUT params
		"POST": flattenValues(filterParams(r.Form)),
	}
}

// filterParams filters sensitive information like passwords from being sent to
// Rollbar.
func filterParams(values map[string][]string) map[string][]string {
	for key, _ := range values {
		if FilterFields.Match([]byte(key)) {
			values[key] = []string{FILTERED}
		}
	}

	return values
}

func flattenValues(values map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range values {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
		}
	}

	return result
}

// Build a message inner-body for the given message string.
func messageBody(s string) map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]interface{}{
			"body": s,
		},
	}
}

func errorClass(err error) string {
	class := reflect.TypeOf(err).String()
	if class == "" {
		return "panic"
	} else if class == "*errors.errorString" {
		checksum := adler32.Checksum([]byte(err.Error()))
		return fmt.Sprintf("{%x}", checksum)
	} else {
		return strings.TrimPrefix(class, "*")
	}
}

// -- POST handling

// Queue the given JSON body to be POSTed to Rollbar.
func push(body map[string]interface{}) {
	if len(bodyChannel) < Buffer {
		waitGroup.Add(1)
		bodyChannel <- body
	} else {
		stderr("buffer full, dropping error on the floor")
	}
}

// POST the given JSON body to Rollbar synchronously.
func post(body map[string]interface{}) {
	if len(Token) == 0 {
		stderr("empty token")
		return
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		stderr("failed to encode payload: %s", err.Error())
		return
	}

	resp, err := http.Post(Endpoint, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		stderr("POST failed: %s", err.Error())
	} else if resp.StatusCode != 200 {
		stderr("received response: %s", resp.Status)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

// -- stderr

func stderr(format string, args ...interface{}) {
	format = "Rollbar error: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
}
