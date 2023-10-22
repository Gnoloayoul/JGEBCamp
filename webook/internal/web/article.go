package web

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/gin-gonic/gin"
)

var _ handler = (*ArticleHandler)(nil)
type ArticleHandler struct {
	l logger.LoggerV1
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {

}
