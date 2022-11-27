package router

import (
	"bazaar/internal/domain"
	"bazaar/internal/domain/browse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	app *domain.Domain
}

func New(router *gin.Engine, app *domain.Domain) error {
	h := &Handler{
		app: app,
	}

	router.POST("/bazaar/v1/upload", h.upload)
	router.GET("/bazaar/v1/query", h.query)
	router.GET("/bazaar/v1/download", h.download)
	return nil
}

// 上传样本
func (h *Handler) upload(c *gin.Context) {
	req := new(browse.Malware)

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, NewReply(ErrParamVerify))
		return
	}

	res, err := h.app.MalwareCreate(req)
	if err != nil {
		c.JSON(http.StatusOK,
			NewReply(ErrSampleUpload).WithErr(err))
		return
	}

	c.JSON(http.StatusOK, NewReply(StatusOk).WithData(res))
}

// 查询样本
func (h *Handler) query(c *gin.Context) {

}

// 下载样本
func (h *Handler) download(c *gin.Context) {

}
