package emqx

import (
	"github.com/gin-gonic/gin"
)

type ACLReq struct {
	ClientID string `json:"clientid"`
	Username string `json:"username"`
	Topic    string `json:"topic"`
	Action   string `json:"action"` // publish | subscribe
}

func ACLHandler(c *gin.Context) {
	var req ACLReq
	_ = c.ShouldBindJSON(&req)

}
