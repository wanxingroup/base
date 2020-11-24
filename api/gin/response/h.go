package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
)

var (
	UnknownError = ResponseData{Status: http.StatusInternalServerError, Code: errors.CodeServerInternalError, Data: nil, Message: "Unknown Error"} // Status 临时改成200值，等待全部替换 Manba 后即可用 400 替换
)

type ResponseData struct {
	Status  int         `json:"-"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg,omitempty"`
}

func (data ResponseData) Error() string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logrus.WithField("error", err).Info("ResponseData.Error()")
		// 这里错误类型是定义好的，正常不应该出现Marshal错误，并且对这个情况，现在用单测覆盖这部分输出，保证输出正常。
		jsonBytes, _ = json.Marshal(UnknownError)
	}
	return string(jsonBytes)
}

func (data ResponseData) String() string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logrus.WithField("error", err).Info("ResponseData.String():", err)
		return ""
	}
	return string(jsonBytes)
}

func (data ResponseData) Clone() *ResponseData {

	clonedData := data
	return &clonedData
}

func (data ResponseData) WithCode(code int) *ResponseData {
	data.Code = code
	return &data
}

func (data ResponseData) WithError(err error) *ResponseData {
	switch err := err.(type) {
	case *ResponseData:
		data.Code = err.Code
		data.Message = err.Message
	case validation.Error:

		validationErrorData := data.withValidationError(err)

		data.Status = validationErrorData.Status
		data.Code = validationErrorData.Code
		data.Message = validationErrorData.Message

	case validation.Errors:

		data.Status = http.StatusBadRequest

		if len(err) == 1 {

			for _, validationError := range err {

				validationErrorData := data.withValidationError(validationError)

				data.Code = validationErrorData.Code
				data.Message = validationErrorData.Message
				break
			}

			break
		}

		data.Message = err.Error()
	default:
		data.Message = err.Error()
	}

	return &data
}

func (data ResponseData) WithMessage(message string) *ResponseData {
	data.Message = message
	return &data
}

func (data ResponseData) withValidationError(err error) *ResponseData {

	switch err := err.(type) {
	case validation.Error:

		data.Status = http.StatusBadRequest
		code, convertError := strconv.ParseInt(err.Code(), 10, 64)
		if convertError != nil {
			data.Message = err.Error()
			break
		}

		data.Code = int(code)
		data.Message = err.Message()
	default:
		data.Message = err.Error()
	}

	return &data
}

func E(c *gin.Context, httpStatus, errCode int, message string) {

	c.JSON(httpStatus, ResponseData{
		Code:    errCode,
		Message: message,
	})
}

func InternalErr(c *gin.Context, errCode int, message string) {
	E(c, http.StatusInternalServerError, errCode, message)
}

func Forbidden(c *gin.Context, errCode int, message string) {
	E(c, http.StatusForbidden, errCode, message)
}

func TokenInvalid(c *gin.Context) {
	E(c, http.StatusUnauthorized, errors.CodeRequestTokenInvalid, "request 'token' is invalid")
}

func TokenExpired(c *gin.Context) {
	E(c, http.StatusUnauthorized, errors.CodeRequestTokenExpired, "request 'token' is expired")
}

func ParamErr(c *gin.Context, message string) {
	E(c, http.StatusBadRequest, errors.CodeRequestParamError, message)
}

func ParamRequired(c *gin.Context, name string) {
	E(c, http.StatusBadRequest, errors.CodeRequestParamError, fmt.Sprintf("%v: %v", name, "cannot be empty"))
}

func ParamInvalid(c *gin.Context, name string) {
	E(c, http.StatusBadRequest, errors.CodeRequestParamError, fmt.Sprintf("%v: %v", name, "invalid value"))
}

func ResponseJSONP(c *gin.Context, data interface{}) {

	if c.Request.Method == "POST" {
		c.JSONP(http.StatusCreated, data)
		return
	}
	if c.Request.Method == "DELETE" {
		c.JSONP(http.StatusNoContent, data)
		return
	}
	c.JSONP(http.StatusOK, data)
}

func Response(c *gin.Context, body interface{}) {

	if c.Request.Method == "POST" {
		c.JSON(http.StatusCreated, body)
		return
	}
	if c.Request.Method == "DELETE" {
		c.JSON(http.StatusNoContent, body)
		return
	}
	c.JSON(http.StatusOK, body)
}

func Error(c *gin.Context, e error) {

	logrus.Tracef("response error: (%T)%v", e, e)
	switch err := e.(type) {
	case *ResponseData:
		c.JSON(err.Status, err)

	case ResponseData:
		c.JSON(err.Status, err)

	case validation.Error:

		errorCode, convertError := strconv.ParseInt(err.Code(), 10, 64)
		if convertError != nil {

			c.JSON(http.StatusBadRequest, ResponseData{
				Code:    errors.CodeRequestParamError,
				Message: err.Message(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, ResponseData{
			Code:    int(errorCode),
			Message: err.Message(),
		})

	case validation.Errors:

		if len(err) > 1 {

			c.JSON(http.StatusBadRequest, ResponseData{
				Code:    errors.CodeRequestParamError,
				Message: err.Error(),
			})
			return
		}

		for _, validationError := range err {

			Error(c, validationError)
			return
		}

	default:
		c.JSON(UnknownError.Status, &ResponseData{
			Code:    UnknownError.Code,
			Message: err.Error(),
		})
	}
}

func Data(c *gin.Context, data interface{}) {

	Response(c, ResponseData{
		Data: data,
	})
}
