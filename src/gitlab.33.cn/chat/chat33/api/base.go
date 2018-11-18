package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WriteMsg 给接口请求返回消息
func WriteMsg(w io.Writer, msg []byte) {
	var err error
	_, err = w.Write(msg)
	if err != nil {
		//logAPI.Warn("Write err", "err ", err)
	}
}

func WriteJson(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusOK, msg)
}
