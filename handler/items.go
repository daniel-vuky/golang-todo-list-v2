package handler

import (
	"fmt"
	"github.com/daniel-vuky/golang-todo-list-and-chat/model"
	"github.com/daniel-vuky/golang-todo-list-and-chat/repository"
	"github.com/gin-gonic/gin"
	ginSession "github.com/go-session/gin-session"
	"net/http"
	"strconv"
)

const DEFAULT_SIZE = 100
const DEFAULT_PAGE = 1

const MISSING_INPUT_ID = "Please enter ID"
const BIND_INPUT_ERR = "Can not bind the input params"
const CREATE_ITEM_ERR = "Can not create to do item, %s"
const UPDATE_ITEM_ERR = "Can not update item, %s"
const DELETE_ITEM_ERR = "Can not delete item, %s"
const FIND_ALL_ITEM_ERR = "Can not get the list items, %s"
const FIND_ITEM_ERR = "Can not find the item with ID, %d, %s"

type Items struct {
	Repository *repository.ItemsRepository
}

// Create create to do item
func (items Items) Create(c *gin.Context) {
	var itemInput model.ItemInput
	bindErr := c.ShouldBindJSON(&itemInput)
	if bindErr != nil || len(itemInput.Title) == 0 || itemInput.Status == 0 {
		WriteResult(http.StatusBadRequest, BIND_INPUT_ERR, c)
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("you need to login first"))
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
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(CREATE_ITEM_ERR, insertErr.Error()), c)
		return
	}
	WriteResultWithItem(http.StatusOK, newItem, c)
}

// List get list to do items
func (items Items) List(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("size"))
	if pageSize == 0 {
		pageSize = DEFAULT_SIZE
	}
	currentPage, _ := strconv.Atoi(c.Query("p"))
	if currentPage == 0 {
		currentPage = DEFAULT_PAGE
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("you need to login first"))
		return
	}
	listItems, findAllErr := items.Repository.FindAll(
		pageSize,
		(currentPage-1)*pageSize,
		userID.(uint64),
	)
	if findAllErr != nil {
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(FIND_ALL_ITEM_ERR, findAllErr.Error()), c)
		return
	}
	WriteResultWithListItems(http.StatusOK, listItems, c)
}

// GetByID get to do item by ID
func (items Items) GetByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		WriteResult(http.StatusBadRequest, MISSING_INPUT_ID, c)
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("you need to login first"))
		return
	}
	item, findItemErr := items.Repository.Find(itemId, userID.(uint64))
	if findItemErr != nil || item.ItemId == 0 {
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(FIND_ITEM_ERR, itemId, findItemErr.Error()), c)
		return
	}
	WriteResultWithItem(http.StatusOK, item, c)
}

// UpdateByID update to do item by ID
func (items Items) UpdateByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		WriteResult(http.StatusBadRequest, MISSING_INPUT_ID, c)
		return
	}
	var itemInput model.ItemInput
	bindErr := c.ShouldBindJSON(&itemInput)
	if bindErr != nil || len(itemInput.Title) == 0 || itemInput.Status == 0 {
		WriteResult(http.StatusBadRequest, BIND_INPUT_ERR, c)
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("you need to login first"))
		return
	}
	item, findItemErr := items.Repository.Find(itemId, userID.(uint64))
	if findItemErr != nil || item.ItemId == 0 {
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(FIND_ITEM_ERR, itemId, findItemErr.Error()), c)
		return
	}
	if updatedErr := items.Repository.Update(&item, &itemInput); updatedErr != nil {
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(UPDATE_ITEM_ERR, updatedErr.Error()), c)
		return
	}
	WriteResult(http.StatusOK, "Updated", c)
}

// DeleteByID delete to do item by item ID
func (items Items) DeleteByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		WriteResult(http.StatusBadRequest, MISSING_INPUT_ID, c)
		return
	}
	session := ginSession.FromContext(c)
	userID, userIDExisted := session.Get("user_id")
	if !userIDExisted {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("you need to login first"))
		return
	}
	item, findItemErr := items.Repository.Find(itemId, userID.(uint64))
	if findItemErr != nil || item.ItemId == 0 {
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(FIND_ITEM_ERR, itemId, findItemErr.Error()), c)
		return
	}
	if deletedErr := items.Repository.Delete(itemId); deletedErr != nil {
		WriteResult(
			http.StatusInternalServerError,
			fmt.Sprintf(DELETE_ITEM_ERR, deletedErr.Error()),
			c,
		)
		return
	}
	WriteResult(http.StatusOK, "Deleted record", c)
}
