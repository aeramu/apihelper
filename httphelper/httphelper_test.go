package httphelper_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aeramu/apihelper/exception"
	"github.com/aeramu/apihelper/httphelper"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Foo string
	Bar string
}

type testCase struct {
	name    string
	handler http.HandlerFunc
}

var (
	errGeneric   = errors.New("some error")
	errException = exception.InvalidRequest("TEST_ERROR", "TEST_MESSAGE", errors.New("some error"))
	errHTTP      = errException.(httphelper.HTTPError)
	defaultCode   = "TEST_ERROR"
	defaultMessage = "TEST_MESSAGE"
)

func getTestCases() []testCase {
	httphelper.Configure(
		httphelper.WithDefaultErrorCode(defaultCode),
		httphelper.WithDefaultErrorMessage(defaultMessage),
	)

	return []testCase{
		{
			name: "return OK - standardized format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := Data{
					Foo: "foo",
					Bar: "bar",
				}
				httphelper.OK(w, data)
			},
		},
		{
			name: "return OK - unstandardized format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"foo":"bar"}`))
			},
		},
		{
			name: "return OK - unable to decode",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"success":"bar"}`))
			},
		},
		{
			name: "return OK - empty body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name: "return OK - malformed JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"foo": "bar"`))
			},
		},
		{
			name: "return error - generic error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				httphelper.Error(w, errGeneric)
			},
		},
		{
			name: "return error - exception error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				httphelper.Error(w, errException)
			},
		},
		{
			name: "return error - empty body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
		{
			name: "return error - unable to decode",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"some error"}`))
			},
		},
		{
			name: "return error - unstandardized format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"errors":"some error"}`))
			},
		},
		{
			name: "return error - malformed JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"foo": "bar"`))
			},
		},
	}
}

func TestHTTPHelper(t *testing.T) {
	tests := getTestCases()
	validate := map[string]func(*testing.T, *http.Response, []byte){
		"return OK - standardized format": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)
			errInfo := result.GetErrorInfo()
			assert.Nil(t, errInfo)
			data, err := httphelper.ReadData[Data](result)
			assert.NoError(t, err)
			assert.Equal(t, "foo", data.Foo)
			assert.Equal(t, "bar", data.Bar)
		},
		"return OK - unstandardized format": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return OK - unable to decode": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.Error(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return OK - empty body": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.Error(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return OK - malformed JSON response": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.Error(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - generic error": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, defaultCode, errInfo.Code)
			assert.Equal(t, defaultMessage, errInfo.Message)
			assert.Equal(t, errGeneric.Error(), errInfo.Details)
		},
		"return error - exception error": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, errHTTP.HTTPCode(), resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, errHTTP.Code(), errInfo.Code)
			assert.Equal(t, errHTTP.Message(), errInfo.Message)
			assert.Equal(t, errHTTP.Error(), errInfo.Details)
		},
		"return error - empty body": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.Error(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - unable to decode": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.Error(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - unstandardized format": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - malformed JSON response": func(t *testing.T, resp *http.Response, body []byte) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			var result httphelper.Response
			err := json.Unmarshal(body, &result)
			assert.Error(t, err)
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			resp, err := http.Get(ts.URL)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()

			validate[tt.name](t, resp, body)
			exampleHTTPImpl(t, tt.name, ts)
		})
	}
}

func exampleHTTPImpl(t *testing.T, name string, ts *httptest.Server) {
	t.Helper()
	fmt.Println("=====================")
	fmt.Printf("%s output:\n", name)

	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error read body: %v\n", err)
		return
	}

	var res httphelper.Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Printf("error decode response: %v\n", err)
		fmt.Printf("status code: %d\n", resp.StatusCode)
		fmt.Printf("body: %s\n", string(body))
		return
	}

	if errInfo := res.GetErrorInfo(); errInfo != nil {
		fmt.Printf("error: %v\n", errInfo.Error())
		fmt.Printf("code: %s\n", errInfo.Code)
		fmt.Printf("message: %s\n", errInfo.Message)
		fmt.Printf("details: %s\n", errInfo.Details)
		fmt.Printf("status code: %d\n", resp.StatusCode)
		fmt.Printf("body: %s\n", string(body))
		return
	}

	result, err := httphelper.ReadData[Data](res)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func TestRestyHelper(t *testing.T) {
	client := resty.New()
	tests := getTestCases()

	validate := map[string]func(*testing.T, *resty.Response, httphelper.Response){
		"return OK - standardized format": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusOK, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.Nil(t, errInfo)
			data, err := httphelper.ReadData[Data](result)
			assert.NoError(t, err)
			assert.Equal(t, "foo", data.Foo)
			assert.Equal(t, "bar", data.Bar)
		},
		"return OK - unstandardized format": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusOK, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return OK - unable to decode": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusOK, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return OK - empty body": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusOK, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return OK - malformed JSON response": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusOK, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - generic error": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, defaultCode, errInfo.Code)
			assert.Equal(t, defaultMessage, errInfo.Message)
			assert.Equal(t, errGeneric.Error(), errInfo.Details)
		},
		"return error - exception error": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, errHTTP.HTTPCode(), resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, errHTTP.Code(), errInfo.Code)
			assert.Equal(t, errHTTP.Message(), errInfo.Message)
			assert.Equal(t, errHTTP.Error(), errInfo.Details)
		},
		"return error - empty body": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - unable to decode": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - unstandardized format": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
		"return error - malformed JSON response": func(t *testing.T, resp *resty.Response, result httphelper.Response) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
			errInfo := result.GetErrorInfo()
			assert.NotNil(t, errInfo)
			assert.Equal(t, httphelper.UNKNOWN_ERROR, errInfo.Code)
			assert.Equal(t, httphelper.UNKNOWN_MESSAGE, errInfo.Message)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			var result httphelper.Response
			resp, err := client.R().
				SetResult(&result).
				SetError(&result).
				Get(ts.URL)
			assert.NoError(t, err)

			validate[tt.name](t, resp, result)
			exampleRestyImpl(t, tt.name, ts)
		})
	}
}

func exampleRestyImpl(t *testing.T, name string, ts *httptest.Server) {
	t.Helper()
	fmt.Println("=====================")
	fmt.Printf("%s output:\n", name)

	client := resty.New()
	var result httphelper.Response
	resp, err := client.R().
		SetResult(&result).
		SetError(&result).
		Get(ts.URL)
	assert.NoError(t, err)

	if errInfo := result.GetErrorInfo(); errInfo != nil {
		fmt.Printf("error: %v\n", errInfo.Error())
		fmt.Printf("code: %s\n", errInfo.Code)
		fmt.Printf("message: %s\n", errInfo.Message)
		fmt.Printf("details: %s\n", errInfo.Details)
		fmt.Printf("status code: %d\n", resp.StatusCode())
		fmt.Printf("body: %s\n", string(resp.Body()))
		return
	}

	data, err := httphelper.ReadData[Data](result)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)
}
