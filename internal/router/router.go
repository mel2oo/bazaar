package router

import (
	"bazaar/internal/domain"
	"bazaar/internal/domain/browse"
	"io"
	"net/http"
	"os"

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
	router.GET("/bazaar/v1/count", h.count)
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
	req := new(browse.QueryMeta)

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, NewReply(ErrParamVerify))
		return
	}

	res, err := h.app.MalwareQuery(req)
	if err != nil {
		c.JSON(http.StatusOK,
			NewReply(ErrSampleQuery).WithErr(err))
		return
	}

	c.JSON(http.StatusOK, NewReply(StatusOk).WithData(res))
}

// 样本数量
func (h *Handler) count(c *gin.Context) {
	req := new(browse.QueryMeta)

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, NewReply(ErrParamVerify))
		return
	}

	res, err := h.app.MalwareCount(req)
	if err != nil {
		c.JSON(http.StatusOK,
			NewReply(ErrSampleCount).WithErr(err))
		return
	}

	c.JSON(http.StatusOK, NewReply(StatusOk).WithData(res))
}

// 下载样本
func (h *Handler) download(c *gin.Context) {
	md5 := c.Query("md5")
	if len(md5) == 0 {
		c.JSON(http.StatusOK, NewReply(ErrParamVerify))
		return
	}

	res, err := h.app.MalwareDownload(md5)
	if err != nil {
		c.JSON(http.StatusOK,
			NewReply(ErrSampleDownload).WithErr(err))
		return
	}

	fi, err := os.Open(res.Path)
	if err != nil {
		c.JSON(http.StatusOK,
			NewReply(ErrSampleExists).WithErr(err))
		return
	}
	defer fi.Close()

	c.Writer.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+res.Name)
	_, err = io.Copy(c.Writer, fi)
	if err != nil {
		c.JSON(http.StatusOK,
			NewReply(ErrSampleReader).WithErr(err))
		return
	}
}
