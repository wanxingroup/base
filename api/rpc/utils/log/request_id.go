package log

import (
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
)

// usage: RequestEntry(c).Debug(".....")
func RequestEntry(c context.Context) *logrus.Entry {

	return WithRequestId(c, logrus.NewEntry(logrus.New()))
}

func WithRequestId(c context.Context, entry *logrus.Entry) *logrus.Entry {

	requestId := requestid.GetRequestIdFromRPCContext(c)
	if xstrings.IsNotEmpty(requestId) {
		return entry.WithField("reqId", requestId)
	} else {
		return entry.WithField("reqId", "unknown")
	}
}
