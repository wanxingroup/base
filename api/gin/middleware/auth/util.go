package auth

import (
	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (u64 uint64) {
	if val, ok := c.Get(UserID); ok && val != nil {
		u64, _ = val.(uint64)
	}
	return
}

func GetUserInformation(c *gin.Context) (info map[string]string) {
	if val, ok := c.Get(UserInformation); ok && val != nil {
		info, _ = val.(map[string]string)
	}
	return
}
