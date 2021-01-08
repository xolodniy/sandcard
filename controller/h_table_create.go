package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controller) createTable(ctx *gin.Context) {
	id, err := c.app.CreateTable()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"tableID": id,
	})
}
