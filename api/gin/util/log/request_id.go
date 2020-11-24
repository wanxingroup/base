package log

import (
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
)

// usage: RequestEntry(c).Debug(".....")
func RequestEntry(c *gin.Context) *logrus.Entry {

	return WithRequestId(c, logrus.NewEntry(logrus.New()))
}

func WithRequestId(c *gin.Context, entry *logrus.Entry) *logrus.Entry {

	requestId := requestid.GetRequestId(c)
	if xstrings.IsNotEmpty(requestId) {
		return entry.WithField("reqId", requestId)
	} else {
		return entry.WithField("reqId", "unknown")
	}
}
