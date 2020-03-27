package hello

import (
	"github.com/gin-gonic/gin"
)

// Hello is a demonstration route handler for output "Hello World!".
func Hello(c *gin.Context) {
	c.Writer.WriteString("HelloWorld")
}
