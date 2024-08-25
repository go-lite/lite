package lite

import (
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
			expectedErr:     fmt.Errorf("expected []byte for binary file serialization"),
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
	tests := []struct {
		name            string
		src             interface{}
		contentType     string
		expectedBody    []byte
		expectedErr     error
		expectedErrCode int
	}{
		{
			name:         "Text (plain) serialization",
			src:          "Hello, World!",
			contentType:  "text/plain",
			expectedBody: []byte("Hello, World!"),
			expectedErr:  nil,
		},
		{
			name:         "HTML serialization",
			src:          "<h1>Hello, World!</h1>",
			contentType:  "text/html",
			expectedBody: []byte("<h1>Hello, World!</h1>"),
			expectedErr:  nil,
		},
		{
			name:         "CSS serialization",
			src:          "body {background-color: #f3f3f3;}",
			contentType:  "text/css",
			expectedBody: []byte("body {background-color: #f3f3f3;}"),
			expectedErr:  nil,
		},
		{
			name:         "JavaScript serialization",
			src:          "console.log('Hello, World!');",
			contentType:  "text/javascript",
			expectedBody: []byte("console.log('Hello, World!');"),
			expectedErr:  nil,
		},
		{
			name:         "JSON serialization Error",
			src:          func() {},
			contentType:  "application/json",
			expectedBody: nil,
			expectedErr:  fmt.Errorf("json: unsupported type: func()"),
		},
		{
			name:         "Atom feed serialization",
			src:          "<feed><title>Example Feed</title></feed>",
			contentType:  "application/atom+xml",
			expectedBody: []byte("<feed><title>Example Feed</title></feed>"),
			expectedErr:  nil,
		},
		{
			name:         "RSS feed serialization",
			src:          "<rss><channel><title>Example Channel</title></channel></rss>",
			contentType:  "application/rss+xml",
			expectedBody: []byte("<rss><channel><title>Example Channel</title></channel></rss>"),
			expectedErr:  nil,
		},
		{
			name:         "MathML serialization",
			src:          "<math><mrow><mn>1</mn><mo>+</mo><mn>1</mn></mrow></math>",
			contentType:  "text/mathml",
			expectedBody: []byte("<math><mrow><mn>1</mn><mo>+</mo><mn>1</mn></mrow></math>"),
			expectedErr:  nil,
		},
		{
			name:         "Java descriptor serialization",
			src:          "MIDlet-Name: HelloWorld",
			contentType:  "text/vnd.sun.j2me.app-descriptor",
			expectedBody: []byte("MIDlet-Name: HelloWorld"),
			expectedErr:  nil,
		},
		{
			name:         "WML serialization",
			src:          "<wml><card id=\"card1\"><p>Hello, World!</p></card></wml>",
			contentType:  "text/vnd.wap.wml",
			expectedBody: []byte("<wml><card id=\"card1\"><p>Hello, World!</p></card></wml>"),
			expectedErr:  nil,
		},
		{
			name:         "HTC serialization",
			src:          "behavior: url(#default#homePage);",
			contentType:  "text/x-component",
			expectedBody: []byte("behavior: url(#default#homePage);"),
			expectedErr:  nil,
		},
		{
			name:            "Text serialization with wrong type",
			src:             []byte("Hello, World!"),
			contentType:     "text/plain",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected string for text serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
		},
		{
			name:         "MIDI serialization",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "audio/midi",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:         "MP3 serialization",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "audio/mpeg",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:         "OGG serialization",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "audio/ogg",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:         "MP4 serialization",
			src:          []byte{0x01, 0x02, 0x03},
			contentType:  "video/mp4",
			expectedBody: []byte{0x01, 0x02, 0x03},
			expectedErr:  nil,
		},
		{
			name:            "Binary data serialization with wrong type",
			src:             "string instead of []byte",
			contentType:     "audio/mpeg",
			expectedBody:    nil,
			expectedErr:     fmt.Errorf("expected []byte for binary file serialization"),
			expectedErrCode: fasthttp.StatusInternalServerError,
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
