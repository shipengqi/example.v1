package middleware

import (
	"bytes"
	"github.com/shipengqi/example.v1/apps/blog/pkg/app"
	e2 "github.com/shipengqi/example.v1/apps/blog/pkg/e"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type bodyLogWriter struct {
	gin.ResponseWriter

	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs the each request.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// The basic information
		start := time.Now().UTC()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Skip for the health check requests in release mode.
		if path == "/healthz" {
			return
		}

		// Read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		blw := &bodyLogWriter{ ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = blw
		// Continue.
		c.Next()

		// Get code and message
		var response app.Response
		code, message := e2.OK.Code(), e2.OK.Message()

		// Logging for API
		contentType := blw.Header().Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			return
		}

		// Calculates the latency time
		end := time.Now().UTC()
		latency := end.Sub(start)

		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			error2 := e2.Wrap(err, "unmarshal")
			log.Error().Err(err).Msgf("body: %s", blw.body.Bytes())

			code = error2.Code()
			message = error2.Message()
		} else {
			code = response.Code
			message = response.Msg
		}

		log.Info().Msgf("%s %s | %s | code: %d, message: %s", method, path, latency, code, message)
	}
}
