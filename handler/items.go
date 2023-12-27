package handler

import (
	"errors"
	"fmt"
	"github.com/daniel-vuky/golang-todo-list-v2/model"
	"github.com/daniel-vuky/golang-todo-list-v2/repository"
	"github.com/gin-gonic/gin"
	ginSession "github.com/go-session/gin-session"
	"net/http"
	"strconv"
)

const DefaultSize int = 100
const DefaultPage int = 1

const SessionError string = "you need to login first"
const MissingInputID string = "please enter ID"
const BindInputError string = "can not bind the input params"
const CreateItemError string = "can not create to do item, %s"
const UpdateItemError string = "can not update item, %s"
const DeleteItemError string = "can not delete item, %s"
const FindAllItemError string = "can not get the list items, %s"
const FindItemError string = "can not find the item with ID, %d, %s"

type Items struct {
	Repository *repository.ItemsRepository
}

// Create create to do item
func (items Items) Create(c *gin.Context) {
	var itemInput model.ItemInput
	bindErr := c.ShouldBindJSON(&itemInput)
	if bindErr != nil || len(itemInput.Title) == 0 || itemInput.Status == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(BindInputError))
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, errors.New(SessionError))
		return
	}
	newItem := model.Item{
		UserId:      userID.(uint64),
		Title:       itemInput.Title,
		Description: itemInput.Description,
		Status:      itemInput.Status,
	}
	insertErr := items.Repository.Insert(&newItem)
	if insertErr != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New(CreateItemError))
		return
	}
	WriteResultWithItem(http.StatusOK, newItem, c)
}

// List get list to do items
func (items Items) List(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("size"))
	if pageSize == 0 {
		pageSize = DefaultSize
	}
	currentPage, _ := strconv.Atoi(c.Query("p"))
	if currentPage == 0 {
		currentPage = DefaultPage
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, errors.New(SessionError))
		return
	}
	listItems, findAllErr := items.Repository.FindAll(
		pageSize,
		(currentPage-1)*pageSize,
		userID.(uint64),
	)
	if findAllErr != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(FindAllItemError, findAllErr.Error()))
		return
	}
	WriteResultWithListItems(http.StatusOK, listItems, c)
}

// GetByID get to do item by ID
func (items Items) GetByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(MissingInputID))
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, errors.New(SessionError))
		return
	}
	item, findItemErr := items.Repository.Find(itemId, userID.(uint64))
	if findItemErr != nil || item.ItemId == 0 {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(FindItemError, itemId, findItemErr.Error()))
		return
	}
	WriteResultWithItem(http.StatusOK, item, c)
}

// UpdateByID update to do item by ID
func (items Items) UpdateByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(MissingInputID))
		return
	}
	var itemInput model.ItemInput
	bindErr := c.ShouldBindJSON(&itemInput)
	if bindErr != nil || len(itemInput.Title) == 0 || itemInput.Status == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(BindInputError))
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, errors.New(SessionError))
		return
	}
	item, findItemErr := items.Repository.Find(itemId, userID.(uint64))
	if findItemErr != nil || item.ItemId == 0 {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(FindItemError, itemId, findItemErr.Error()))
		return
	}
	if updatedErr := items.Repository.Update(&item, &itemInput); updatedErr != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(UpdateItemError, updatedErr.Error()))
		return
	}
	WriteResult(http.StatusOK, "Updated", c)
}

// DeleteByID delete to do item by item ID
func (items Items) DeleteByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(MissingInputID))
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, errors.New(SessionError))
		return
	}
	item, findItemErr := items.Repository.Find(itemId, userID.(uint64))
	if findItemErr != nil || item.ItemId == 0 {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(FindItemError, itemId, findItemErr.Error()))
		return
	}
	if deletedErr := items.Repository.Delete(itemId); deletedErr != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(DeleteItemError, deletedErr.Error()))
		return
	}
	WriteResult(http.StatusOK, "Deleted record", c)
}
