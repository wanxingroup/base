package errors

const CodeServerInternalError = 500000     // 服务内部错误
const CodeRequestParamError = 400000       // 请求参数错误
const CodeRequestPathError = 400001        // 请求path错误
const CodeRequestTokenInvalid = 400002     // 请求token无效
const CodeRequestTokenExpired = 400003     // 请求token过期
const CodeRequestJSONDecodeFailed = 400004 // 请求的 JSON 解释失败

var ServerInternalError = Error{
	ErrorCode:    CodeServerInternalError,
	ErrorMessage: "服务内部错误",
}

var RequestParamError = Error{
	ErrorCode:    CodeRequestParamError,
	ErrorMessage: "请求参数错误",
}

var RequestPathError = Error{
	ErrorCode:    CodeRequestPathError,
	ErrorMessage: "请求path错误",
}

var RequestTokenInvalid = Error{
	ErrorCode:    CodeRequestTokenInvalid,
	ErrorMessage: "请求token无效",
}

var RequestTokenExpired = Error{
	ErrorCode:    CodeRequestTokenExpired,
	ErrorMessage: "请求token过期",
}

var RequestJSONDecodeFailed = Error{
	ErrorCode:    CodeRequestJSONDecodeFailed,
	ErrorMessage: "请求的 JSON 解释失败",
}
