package requestid

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"
)

const Header = "X-Request-ID"
const Key = "requestId"

func GenerateRequestId() string {
	var b = make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		// 基本不会出现
		logrus.WithError(err).Warn("generate request id failed")
		return "fail_gen_req_id"
	}
	return hex.EncodeToString(b)
}

func GetRequestId(c *gin.Context) string {

	if requestId := c.GetString(Key); xstrings.IsNotEmpty(requestId) {
		return requestId
	}

	return c.Request.Header.Get(Header)
}

func SetRequestIdIfNotExist(c *gin.Context) string {

	requestId := GetRequestId(c)
	if requestId == "" {
		requestId = GenerateRequestId()
		c.Request.Header.Set(Header, requestId)
		setContext(c, requestId)
	}
	return requestId
}

func InitContext(c *gin.Context) {

	requestId := GetRequestId(c)
	if xstrings.IsNotEmpty(requestId) {
		setContext(c, requestId)
	}
}

func setContext(c *gin.Context, requestId string) {
	c.Set(Key, requestId)
}
