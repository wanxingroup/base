package requestid

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func GetRPCContext(ctx *gin.Context) context.Context {

	requestId := GetRequestId(ctx)
	rpcContext := context.WithValue(context.Background(), Key, requestId)
	return rpcContext
}

func GetRequestIdFromRPCContext(ctx context.Context) string {

	requestId := ctx.Value(Key)
	if requestId == nil {
		return ""
	}

	return requestId.(string)
}
