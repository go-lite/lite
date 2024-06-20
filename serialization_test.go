package lite

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"reflect"
	"testing"
)

func TestSerializeResponse(t *testing.T) {
	tests := []struct {
		name            string
		src             interface{}
		contentType     string
		expectedBody    []byte
		expectedErr     error
		expectedErrCode int
	}{
		{
			name:         "nil src",
			src:          nil,
			contentType:  "",
			expectedBody: nil,
			expectedErr:  nil,
		},
		{
			name:         "string src",
			src:          "Hello, World!",
			contentType:  "",
			expectedBody: []byte("Hello, World!"),
			expectedErr:  nil,
		},
		{
			name:            "unsupported type",
			src:             make(chan int),
			contentType:     "",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("unsupported type: chan"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "JSON serialization",
			src:          map[string]string{"key": "value"},
			contentType:  "application/json",
			expectedBody: []byte(`{"key":"value"}` + "\n"),
			expectedErr:  nil,
		},
		{
			name:         "XML serialization",
			src:          map[string]string{"key": "value"},
			contentType:  "application/xml",
			expectedBody: []byte(`<key>value</key>`),
			expectedErr:  fmt.Errorf("xml: unsupported type: map[string]string"),
		},
		{
			name:         "Form serialization",
			src:          map[string]string{"key": "value"},
			contentType:  "application/x-www-form-urlencoded",
			expectedBody: []byte("key=value"),
			expectedErr:  nil,
		},
		{
			name:            "Form serialization with wrong type",
			src:             map[string]int{"key": 1},
			contentType:     "application/x-www-form-urlencoded",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected map[string]string for form data serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "Binary data",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "application/octet-stream",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:            "Binary data with wrong type",
			src:             1,
			contentType:     "application/octet-stream",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected []byte for octet-stream serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "Unsupported content type",
			src:          map[string]string{"key": "value"},
			contentType:  "unsupported/type",
			expectedBody: nil,
			expectedErr:  fmt.Errorf("unsupported content type: unsupported/type"),
		},
		{
			name:         "Binary data (image/png)",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "image/png",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:            "Binary data (image/png) with wrong type",
			src:             1,
			contentType:     "image/png",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected []byte for binary file serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := new(fasthttp.RequestCtx)
			ctx.Response.Header.SetContentType(tt.contentType)

			err := serializeResponse(ctx, tt.src)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				if tt.expectedErrCode != 0 {
					assert.Equal(t, tt.expectedErrCode, ctx.Response.StatusCode())
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedBody, ctx.Response.Body())
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	type Plant struct {
		XMLName xml.Name `xml:"plant"`
		Id      int      `xml:"id,attr"`
		Name    string   `xml:"name"`
	}

	tests := []struct {
		name            string
		src             interface{}
		contentType     string
		expectedBody    []byte
		expectedErr     error
		expectedErrCode int
	}{
		{
			name:         "JSON serialization",
			src:          map[string]string{"key": "value"},
			contentType:  "application/json",
			expectedBody: []byte(`{"key":"value"}` + "\n"),
			expectedErr:  nil,
		},
		{
			name:            "JSON serialization error",
			src:             func() {},
			contentType:     "application/json",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("json: unsupported type: func()"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "XML serialization",
			src:          Plant{Id: 1, Name: "Rose"},
			contentType:  "application/xml",
			expectedBody: []byte(`<plant id="1"><name>Rose</name></plant>`),
			expectedErr:  nil,
		},
		{
			name:         "Form serialization",
			src:          map[string]string{"key": "value"},
			contentType:  "application/x-www-form-urlencoded",
			expectedBody: []byte("key=value"),
			expectedErr:  nil,
		},
		{
			name:            "Form serialization with wrong type",
			src:             map[string]int{"key": 1},
			contentType:     "application/x-www-form-urlencoded",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected map[string]string for form data serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "Binary data",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "application/octet-stream",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:            "Binary data with wrong type",
			src:             "string instead of []byte",
			contentType:     "application/octet-stream",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected []byte for octet-stream serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "Unsupported content type",
			src:          map[string]string{"key": "value"},
			contentType:  "unsupported/type",
			expectedBody: nil,
			expectedErr:  fmt.Errorf("unsupported content type: unsupported/type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := new(fasthttp.RequestCtx)
			ctx.Response.Header.SetContentType(tt.contentType)

			srcVal := reflect.ValueOf(tt.src)
			err := serialize(ctx, srcVal)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				if tt.expectedErrCode != 0 {
					assert.Equal(t, tt.expectedErrCode, ctx.Response.StatusCode())
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedBody, ctx.Response.Body())
			}
		})
	}
}
