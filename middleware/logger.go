package middleware

import (
	"bytes"
	"io"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go_demo_api/logger"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger(skipPaths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if slices.Contains(skipPaths, path) {
			c.Next()
			return
		}

		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()
		method := c.Request.Method

		var reqbody []byte
		if c.Request.Body != nil {
			reqbody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqbody))
		}

		logger.Log.Info().
			Str("requestId", requestID).
			Str("method", method).
			Str("path", path).
			Str("reqbody", string(reqbody)).
			Msg("Request")

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = rw

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		resbody := rw.body.String()

		logger.Log.Info().
			Str("requestId", requestID).
			Int("status", statusCode).
			Dur("duration", duration).
			Str("respbody", resbody).
			Msg("Response")
	}
}
