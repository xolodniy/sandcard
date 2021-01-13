package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// @Summary Join table
// @Description Subscribe to table and collaborate with other players
// @Description You will receive notifications about all changes on the table.
// @Description You also are allowed to initiate some changes(events) by yourself
// @Description More info about table events in https://github.com/xolodniy/sandcard/blob/master/application/events.md
// @Description
// @Param tableID path int true "TableID" Default(1)
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /api/v1/table/id{tableID}/join [post]
func (c *Controller) joinTable(ctx *gin.Context) {
	ctx.Request.Header.Del("Origin")
	tableID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || tableID < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный идентиффикатор стола"})
		return
	}

	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("can't initialize webSocket connection")
		return
	}
	if err := c.app.JoinTable(conn, tableID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)

}
