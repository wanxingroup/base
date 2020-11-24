package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/response"
)

func Params(c *gin.Context, v Validator) error {

	if err := c.ShouldBindJSON(v); err != nil {
		return response.ResponseData{
			Status: http.StatusOK,
			Code:   errors.CodeRequestParamError,
		}.WithError(err)
	}

	reqId := requestid.GetRequestId(c)
	logrus.Debugf("reqId: %s, v = %+v\n", reqId, v)
	if err := v.Validate(); err != nil {
		return response.ResponseData{
			Status: http.StatusOK,
			Code:   errors.CodeRequestParamError,
		}.WithError(err)
	}
	return nil
}
