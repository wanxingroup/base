package request

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
)

func ReqLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339})

		requestId := requestid.SetRequestIdIfNotExist(c)
		requestid.InitContext(c)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		entry := logrus.WithFields(logrus.Fields{
			"reqId":      requestId,
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL,
			"size":       c.Writer.Size(),
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Info(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}
