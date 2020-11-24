package response

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestResponseData_Error(t *testing.T) {

	tests := []struct {
		Input error
		Want  string
	}{
		{
			Input: &ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			Want: `{"code":400,"msg":"BadRequest"}`,
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			Want: `{"code":400,"msg":"BadRequest"}`,
		},
		{
			Input: &ResponseData{
				Code:    http.StatusBadRequest,
				Message: "Parameters Error",
			},
			Want: `{"code":400,"msg":"Parameters Error"}`,
		},
		{
			Input: &ResponseData{Data: make(chan int)},
			Want:  `{"code":500000,"msg":"Unknown Error"}`,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Want, test.Input.Error())
	}
}

func TestResponseData_String(t *testing.T) {

	tests := []struct {
		Input *ResponseData
		Want  string
	}{
		{
			Input: &ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			Want: `{"code":400,"msg":"BadRequest"}`,
		},
		{
			Input: &ResponseData{
				Code:    http.StatusBadRequest,
				Message: "Parameters Error",
			},
			Want: `{"code":400,"msg":"Parameters Error"}`,
		},
		{
			Input: &ResponseData{Data: make(chan int)},
			Want:  ``,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Want, test.Input.String())
	}
}

func TestResponseData_WithCode(t *testing.T) {

	tests := []struct {
		Input     ResponseData
		InputCode int
		Want      *ResponseData
	}{
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputCode: 500,
			Want: &ResponseData{
				Code:    500,
				Message: "BadRequest",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Want, test.Input.WithCode(test.InputCode))
	}
}

func TestResponseData_WithError(t *testing.T) {

	tests := []struct {
		Input      ResponseData
		InputError error
		Want       *ResponseData
	}{
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: errors.New("BadGateway"),
			Want: &ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadGateway",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: &ResponseData{
				Code:    500001,
				Message: "Parameters error",
			},
			Want: &ResponseData{
				Code:    500001,
				Message: "Parameters error",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: validation.NewError("500001", "Parameters error"),
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    500001,
				Message: "Parameters error",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: validation.NewError("invalid code", "Parameters error"),
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    http.StatusBadRequest,
				Message: "Parameters error",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: func() error {

				validationErrors := validation.Errors{}
				validationErrors["username"] = validation.NewError("500001", "Parameters error")
				return validationErrors
			}(),
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    500001,
				Message: "Parameters error",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: func() error {

				validationErrors := validation.Errors{}
				validationErrors["username"] = validation.NewError("500001", "Parameters error")
				validationErrors["password"] = validation.NewError("500001", "Parameters error")
				return validationErrors
			}(),
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    http.StatusBadRequest,
				Message: "password: Parameters error; username: Parameters error.",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Want, test.Input.WithError(test.InputError), test)
	}
}

func TestResponseData_withValidationError(t *testing.T) {

	tests := []struct {
		Input      ResponseData
		InputError error
		Want       *ResponseData
	}{
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: validation.NewError("invalid code", "Parameters error"),
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    http.StatusBadRequest,
				Message: "Parameters error",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: validation.NewError("500001", "Parameters error"),
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    500001,
				Message: "Parameters error",
			},
		},
		{
			Input: ResponseData{
				Code:    http.StatusBadRequest,
				Message: "BadRequest",
			},
			InputError: fmt.Errorf("normal error"),
			Want: &ResponseData{
				Code:    http.StatusBadRequest,
				Message: "normal error",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Want, test.Input.withValidationError(test.InputError), test)
	}
}

func TestResponseData_WithMessage(t *testing.T) {

	tests := []struct {
		Input        ResponseData
		InputMessage string
		Want         *ResponseData
	}{
		{
			Input: ResponseData{
				Status:  http.StatusBadRequest,
				Code:    400000,
				Data:    nil,
				Message: "BadRequest",
			},
			InputMessage: "Bad Request",
			Want: &ResponseData{
				Status:  http.StatusBadRequest,
				Code:    400000,
				Data:    nil,
				Message: "Bad Request",
			},
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.Want, test.Input.WithMessage(test.InputMessage))
	}
}

func TestE(t *testing.T) {

	tests := []struct {
		InputStatus  int
		InputCode    int
		InputMessage string
		Want         struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			InputStatus:  UnknownError.Status,
			InputCode:    UnknownError.Code,
			InputMessage: UnknownError.Message,
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: UnknownError.Status, ResponseString: `{"code":500000,"msg":"Unknown Error"}`},
		},
		{
			InputStatus:  http.StatusOK,
			InputCode:    500,
			InputMessage: "BadRequest",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusOK, ResponseString: `{"code":500,"msg":"BadRequest"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		E(ctx, test.InputStatus, test.InputCode, test.InputMessage)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String())
	}
}

func TestResponseError(t *testing.T) {

	tests := []struct {
		Input error
		Want  struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			Input: UnknownError.Clone(),
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: UnknownError.Status, ResponseString: `{"code":500000,"msg":"Unknown Error"}`},
		},
		{
			Input: errors.New("BadRequest"),
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusInternalServerError, ResponseString: `{"code":500000,"msg":"BadRequest"}`},
		},
		{
			Input: &ResponseData{
				Status:  http.StatusUnauthorized,
				Code:    400003,
				Message: "unauthorized",
			},
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusUnauthorized, ResponseString: `{"code":400003,"msg":"unauthorized"}`},
		},
		{
			Input: ResponseData{
				Status:  http.StatusUnauthorized,
				Code:    400003,
				Message: "unauthorized",
			},
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusUnauthorized, ResponseString: `{"code":400003,"msg":"unauthorized"}`},
		},
		{
			Input: validation.NewError("400001", "parameter invalid"),
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400001,"msg":"parameter invalid"}`},
		},
		{
			Input: validation.NewError("invalid code", "parameter invalid"),
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400000,"msg":"parameter invalid"}`},
		},
		{
			Input: func() error {

				errs := validation.Errors{}
				errs["password"] = validation.NewError("400001", "password too short")
				return errs
			}(),
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400001,"msg":"password too short"}`},
		},
		{
			Input: func() error {

				errs := validation.Errors{}
				errs["password"] = validation.NewError("400044", "password too short")
				errs["username"] = validation.NewError("400045", "username too short")
				return errs
			}(),
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400000,"msg":"password: password too short; username: username too short."}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		Error(ctx, test.Input)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}

func TestData(t *testing.T) {

	tests := []struct {
		Input struct {
			Body       interface{}
			HTTPMethod string
		}
		Want struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodGet,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusOK, ResponseString: `{"code":0,"data":{"Hello":"World"}}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodPost,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusCreated, ResponseString: `{"code":0,"data":{"Hello":"World"}}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodDelete,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusNoContent, ResponseString: ``},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ctx.Request = &http.Request{Method: test.Input.HTTPMethod}
		Data(ctx, test.Input.Body)
		resp.Flush()

		assert.Equal(t, test.Want.Code, resp.Code)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String())
	}
}

func TestResponse(t *testing.T) {

	tests := []struct {
		Input struct {
			Body       interface{}
			HTTPMethod string
		}
		Want struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodGet,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusOK, ResponseString: `{"Hello":"World"}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodPost,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusCreated, ResponseString: `{"Hello":"World"}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodDelete,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusNoContent, ResponseString: ``},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ctx.Request = &http.Request{Method: test.Input.HTTPMethod}
		Response(ctx, test.Input.Body)
		resp.Flush()

		assert.Equal(t, test.Want.Code, resp.Code)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String())
	}
}

func TestR(t *testing.T) {

	tests := []struct {
		Input struct {
			Body       interface{}
			HTTPMethod string
		}
		Want struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodGet,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusOK, ResponseString: `{"code":0,"data":{"Hello":"World"}}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodPost,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusCreated, ResponseString: `{"code":0,"data":{"Hello":"World"}}`},
		},
		{
			Input: struct {
				Body       interface{}
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				HTTPMethod: http.MethodDelete,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusNoContent, ResponseString: ``},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ctx.Request = &http.Request{Method: test.Input.HTTPMethod}
		Data(ctx, test.Input.Body)
		resp.Flush()

		assert.Equal(t, test.Want.Code, resp.Code)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String())
	}
}

func TestRJsonP(t *testing.T) {

	tests := []struct {
		Input struct {
			Body       interface{}
			URL        string
			HTTPMethod string
		}
		Want struct {
			Code           int
			ResponseString string
		}
	}{
		{
			Input: struct {
				Body       interface{}
				URL        string
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				URL:        "http://localhost/?callback=c",
				HTTPMethod: http.MethodGet,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusOK, ResponseString: `c({"Hello":"World"});`},
		},
		{
			Input: struct {
				Body       interface{}
				URL        string
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				URL:        "http://localhost/?callback=c",
				HTTPMethod: http.MethodPost,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusCreated, ResponseString: `c({"Hello":"World"});`},
		},
		{
			Input: struct {
				Body       interface{}
				URL        string
				HTTPMethod string
			}{
				Body:       struct{ Hello string }{Hello: "World"},
				URL:        "http://localhost/?callback=c",
				HTTPMethod: http.MethodDelete,
			},
			Want: struct {
				Code           int
				ResponseString string
			}{Code: http.StatusNoContent, ResponseString: ``},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		requestURL, err := url.Parse(test.Input.URL)
		assert.Nil(t, err)
		ctx.Request = &http.Request{Method: test.Input.HTTPMethod, URL: requestURL}
		ResponseJSONP(ctx, test.Input.Body)
		resp.Flush()

		assert.Equal(t, test.Want.Code, resp.Code)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String())
	}
}

func TestInternalErr(t *testing.T) {

	tests := []struct {
		InputCode    int
		InputMessage string
		Want         struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			InputCode:    500000,
			InputMessage: "Unknown Error",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusInternalServerError, ResponseString: `{"code":500000,"msg":"Unknown Error"}`},
		},
		{
			InputCode:    400000,
			InputMessage: "BadRequest",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusInternalServerError, ResponseString: `{"code":400000,"msg":"BadRequest"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		InternalErr(ctx, test.InputCode, test.InputMessage)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}

func TestForbidden(t *testing.T) {

	tests := []struct {
		InputCode    int
		InputMessage string
		Want         struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			InputCode:    500000,
			InputMessage: "Unknown Error",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusForbidden, ResponseString: `{"code":500000,"msg":"Unknown Error"}`},
		},
		{
			InputCode:    400000,
			InputMessage: "BadRequest",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusForbidden, ResponseString: `{"code":400000,"msg":"BadRequest"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		Forbidden(ctx, test.InputCode, test.InputMessage)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}

func TestTokenInvalid(t *testing.T) {

	tests := []struct {
		Want struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusUnauthorized, ResponseString: `{"code":400002,"msg":"request 'token' is invalid"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		TokenInvalid(ctx)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String())
	}
}

func TestTokenExpired(t *testing.T) {

	tests := []struct {
		Want struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusUnauthorized, ResponseString: `{"code":400003,"msg":"request 'token' is expired"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		TokenExpired(ctx)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}

func TestParamErr(t *testing.T) {

	tests := []struct {
		Input string
		Want  struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			Input: "shopId is invalid",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400000,"msg":"shopId is invalid"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ParamErr(ctx, test.Input)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}

func TestParamRequired(t *testing.T) {

	tests := []struct {
		Input string
		Want  struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			Input: "shopId",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400000,"msg":"shopId: cannot be empty"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ParamRequired(ctx, test.Input)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}

func TestParamInvalid(t *testing.T) {

	tests := []struct {
		Input string
		Want  struct {
			HTTPStatus     int
			ResponseString string
		}
	}{
		{
			Input: "shopId",
			Want: struct {
				HTTPStatus     int
				ResponseString string
			}{HTTPStatus: http.StatusBadRequest, ResponseString: `{"code":400000,"msg":"shopId: invalid value"}`},
		},
	}

	for _, test := range tests {

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ParamInvalid(ctx, test.Input)
		resp.Flush()

		assert.Equal(t, test.Want.HTTPStatus, resp.Code, test)
		assert.Equal(t, test.Want.ResponseString, resp.Body.String(), test)
	}
}
