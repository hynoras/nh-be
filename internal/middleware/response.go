package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter captures response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	// Capture, but don't immediately write to client yet
	return w.body.Write(b)
}

func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = writer

		c.Next() // execute handler

		// If handler already aborted (e.g. via c.AbortWithStatusJSON), skip
		if c.IsAborted() || writer.body.Len() == 0 {
			return
		}

		status := c.Writer.Status()
		var original interface{}

		// Try to unmarshal the JSON body that handler wrote
		if err := json.Unmarshal(writer.body.Bytes(), &original); err != nil {
			original = gin.H{"raw": writer.body.String()}
		}

		// Read possible custom message
		customMsg, hasMsg := c.Get("message")

		// Build final formatted response
		var response gin.H
		if status >= http.StatusBadRequest {
			response = gin.H{
				"success": false,
				"message": http.StatusText(status),
				"error":   original,
			}
			if hasMsg {
				response["message"] = customMsg
			}
		} else {
			response = gin.H{
				"success": true,
				"message": "OK",
				"data":    original,
			}
			if hasMsg {
				response["message"] = customMsg
			}
		}

		// Rewrite the formatted response
		c.Writer = writer.ResponseWriter // restore original writer
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(status)
		_ = json.NewEncoder(c.Writer).Encode(response)
	}
}
