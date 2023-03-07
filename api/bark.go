package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/model"
	"github.com/lithammer/shortuuid/v3"
)

type BarkDatabase interface {
	CreateBark(bark *model.Bark) error
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
	key := shortuuid.New()
	token, _ := ctx.GetQuery("devicetoken")

	model := &model.Bark{
		Key:   key,
		Token: token,
	}

	if success := successOrAbort(ctx, 500, a.DB.CreateBark(model)); !success {
		return
	}

	ctx.JSON(200, NewReply(map[string]string{
		// compatible with old resp
		"key":          key,
		"device_key":   key,
		"device_token": token,
	}))
}
