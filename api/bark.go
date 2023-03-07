package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/model"
)

type BarkDatabase interface {
	FindBarkByToken(token string) (*model.Bark, error)
	DeleteByKey(key string) error
}

type BarkAPI struct {
	DB BarkDatabase
}

type BarkReply struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

func NewReply(data interface{}) BarkReply {
	return BarkReply{
		Code:      200,
		Message:   "success",
		Timestamp: timeNow().Unix(),
		Data:      data,
	}
}

func (a *BarkAPI) Ping(ctx *gin.Context) {
	reply := BarkReply{
		Code:      200,
		Message:   "pong",
		Timestamp: timeNow().Unix(),
	}

	ctx.JSON(200, reply)
}

func (a *BarkAPI) Register(ctx *gin.Context) {
	key, _ := ctx.GetQuery("key")
	token, _ := ctx.GetQuery("devicetoken")

	if token == "deleted" {
		a.DB.DeleteByKey(key)
		ctx.JSON(200, NewReply(nil))

		return
	}

	bark, err := a.DB.FindBarkByToken(token)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}

	ctx.JSON(200, NewReply(map[string]string{
		"key":          bark.Key,
		"device_key":   bark.Key,
		"device_token": bark.Token,
	}))

}
