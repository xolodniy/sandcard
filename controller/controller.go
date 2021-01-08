package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
	c.router.POST("/api/v1/table", c.createTable)
	c.router.GET("/api/v1/table/:id/join", c.joinTable)
	return c
}

func (c *Controller) Serve(port int) {
	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: c.router,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("can't start serving http")
	}
}
