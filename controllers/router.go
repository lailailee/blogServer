package controllers

import (
	"blog/core"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// HTTPAPI  http api module
type HTTPAPI struct {
	app    *gin.Engine
	lock   sync.Mutex
	logger *zap.SugaredLogger
}

// NewHTTPAPIServer new a http server
func NewHTTPAPIServer(logger *zap.SugaredLogger) *HTTPAPI {
	v := &HTTPAPI{
		app:    gin.Default(),
		logger: logger,
	}
	return v
}
func index(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/admin/")
}

func (h *HTTPAPI) initialize() {
	// v.lock.Lock()
	// defer v.lock.Unlock()
	// router := h.app
	h.app.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "OPTIONS, GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))
	h.app.POST("/v1/login", h.Login)

	h.app.GET("/v1/article", h.GetArticleList)
	h.app.GET("/v1/article/:id", h.GetArticle)
	h.app.POST("/v1/article", h.AddArticle)
	h.app.PUT("/v1/article/:id", h.EditArticle)
	h.app.DELETE("/v1/article/:id", h.DeleteArticle)

	h.app.GET("/v1/tag", h.GetTag)
	h.app.POST("/v1/tag", h.AddTag)
	h.app.PUT("/v1/tag/:id", h.EditTag)
	h.app.DELETE("/v1/tag/:id", h.DeleteTag)

	h.app.GET("/v1/category", h.GetCategory)
	h.app.POST("/v1/category", h.AddCategory)
	h.app.PUT("/v1/category/:id", h.EditCategory)
	h.app.DELETE("/v1/category/:id", h.DeleteCategory)

	h.app.GET("/v1/series", h.GetSeries)
	h.app.POST("/v1/series", h.AddSeries)
	h.app.PUT("/v1/series/:id", h.EditSeries)
	h.app.DELETE("/v1/series/:id", h.DeleteSeries)
}

func (h *HTTPAPI) Run() {
	h.initialize()

	if e0 := h.app.Run(":" + strconv.Itoa(core.Conf.HTTPAPI.Port)); e0 != nil {
		h.logger.Errorf("http api server run error, error=[%v]", e0)
		os.Exit(-1)
	}
}
