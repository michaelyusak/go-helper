package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/appconstant"

	"github.com/google/uuid"
)

func RequestIdHandlerMiddleware(c *gin.Context) {
	uuid := uuid.NewString()

	requestId := fmt.Sprintf("%v:%s", time.Now().UnixMilli(), uuid)

	c.Set(appconstant.RequestId, requestId)

	c.Next()
}
