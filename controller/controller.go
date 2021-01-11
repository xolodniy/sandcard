package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Controller struct {
	router *gin.Engine
	app    Application
}

type Application interface {
	CreateTable() (int, error)
	JoinTable(c *websocket.Conn, tableID int) error
}

func New(a Application) *Controller {
	c := &Controller{
		router: gin.Default(),
		app:    a,
	}
	c.router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	c.router.POST("/api/v1/table", c.createTable)
	c.router.GET("/api/v1/table/id:id/join", c.joinTable)
	return c
}

func (c *Controller) Serve(port int) {
	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: c.router,
	}

	fmt.Printf("http server started on :%d\n", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("can't start serving http")
	}
}
