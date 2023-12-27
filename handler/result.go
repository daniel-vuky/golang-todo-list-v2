package handler

import (
	"fmt"
	"github.com/daniel-vuky/golang-todo-list-v2/model"
	"github.com/gin-gonic/gin"
	"net/http"
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

func Redirect(path string, errorCode int, c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("/%s?error=%d", path, errorCode))
}
