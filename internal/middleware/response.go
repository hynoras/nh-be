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

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the ResponseWriter
		writer := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next() // Run handler

		status := c.Writer.Status()

		// Skip if already JSON error or Gin handled response
		if c.IsAborted() || writer.body.Len() == 0 {
			return
		}

		var original interface{}
		if err := json.Unmarshal(writer.body.Bytes(), &original); err != nil {
			original = gin.H{"raw": writer.body.String()}
		}

		// Build final response format
		var response gin.H
		if status >= http.StatusBadRequest {
			response = gin.H{
				"success": false,
				"message": http.StatusText(status),
				"error":   original,
			}
		} else {
			response = gin.H{
				"success": true,
				"message": "OK",
				"data":    original,
			}
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeaderNow()
		json.NewEncoder(c.Writer).Encode(response)
	}
}
