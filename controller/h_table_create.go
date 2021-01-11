package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateTableResponse struct {
	ID int `json:"id" example:"1"`
}

// @Summary Create new table
// @Description Returns table id which should be used for joining
// @Description
// @Accept  json
// @Produce json
// @Success 200 {object} CreateTableResponse
// @Failure 500
// @Router /api/v1/table [post]
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
