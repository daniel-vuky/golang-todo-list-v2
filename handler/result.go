package handler

import (
	"github.com/daniel-vuky/golang-todo-list-and-chat/model"
	"github.com/gin-gonic/gin"
)

type ListItem []model.Item

// WriteResult write the result code and message to the gin context
func WriteResult(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

// WriteResultWithListItems write the result code and result to the gin context
func WriteResultWithListItems(code int, result ListItem, c *gin.Context) {
	c.JSON(code, result)
}

// WriteResultWithItem write the result code and result to the gin context
func WriteResultWithItem(code int, result model.Item, c *gin.Context) {
	c.JSON(code, result)
}
