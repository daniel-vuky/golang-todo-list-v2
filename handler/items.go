package handler

import (
	"fmt"
	"github.com/daniel-vuky/golang-todo-list-and-chat/model"
	"github.com/daniel-vuky/golang-todo-list-and-chat/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Items struct {
	Repository *repository.ItemsRepository
}

// Create method of Items
// @param gin.context c
func (items Items) Create(c *gin.Context) {
	var itemInput model.ItemInput
	bindErr := c.ShouldBindJSON(&itemInput)
	if bindErr != nil || len(itemInput.Title) == 0 || itemInput.Status == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Can not bind the input params",
			c,
		)
		return
	}
	newItem := model.Item{
		UserId:      1,
		Title:       itemInput.Title,
		Description: itemInput.Description,
		Status:      itemInput.Status,
	}
	itemId, insertErr := items.Repository.Insert(newItem)
	if insertErr != nil {
		WriteResult(
			http.StatusInternalServerError,
			fmt.Sprintf("Can not create to do item, %s", insertErr.Error()),
			c,
		)
		return
	}
	WriteResultWithItem(http.StatusOK, model.Item{ItemId: uint64(itemId)}, c)
}

func (items Items) List(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("size"))
	if pageSize == 0 {
		pageSize = 9999
	}
	currentPage, _ := strconv.Atoi(c.Query("p"))
	if currentPage == 0 {
		currentPage = 1
	}
	listItems, findAllErr := items.Repository.FindAll(
		pageSize,
		(currentPage-1)*pageSize,
	)
	if findAllErr != nil {
		WriteResult(
			http.StatusInternalServerError,
			fmt.Sprintf("Can not create to do item, %s", findAllErr.Error()),
			c,
		)
		return
	}
	WriteResultWithListItems(http.StatusOK, listItems, c)
}

func (items Items) GetByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Please enter ID",
			c,
		)
		return
	}
	item, findItemErr := items.Repository.Find(itemId)
	if findItemErr != nil || item.ItemId == 0 {
		WriteResult(
			http.StatusInternalServerError,
			fmt.Sprintf("Can not find the item with ID, %d, %s", itemId, findItemErr.Error()),
			c,
		)
		return
	}
	WriteResultWithItem(http.StatusOK, item, c)
}

func (items Items) UpdateByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Please enter ID",
			c,
		)
		return
	}
	var itemInput model.ItemInput
	bindErr := c.ShouldBindJSON(&itemInput)
	if bindErr != nil || len(itemInput.Title) == 0 || itemInput.Status == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Can not bind the input params",
			c,
		)
		return
	}
	item, findItemErr := items.Repository.Find(itemId)
	if findItemErr != nil || item.ItemId == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Can not find item match with this ID",
			c,
		)
		return
	}
	if updatedErr := items.Repository.Update(item, itemInput); updatedErr != nil {
		WriteResult(
			http.StatusInternalServerError,
			fmt.Sprintf("Can not update the item, %s", updatedErr.Error()),
			c,
		)
		return
	}
	WriteResult(http.StatusOK, "Updated", c)
}

func (items Items) DeleteByID(c *gin.Context) {
	itemId, itemIdErr := strconv.Atoi(c.Param("id"))
	if itemIdErr != nil || itemId == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Please enter ID",
			c,
		)
		return
	}
	item, findItemErr := items.Repository.Find(itemId)
	if findItemErr != nil || item.ItemId == 0 {
		WriteResult(
			http.StatusBadRequest,
			"Can not find item match with this ID",
			c,
		)
		return
	}
	if deletedErr := items.Repository.Delete(itemId); deletedErr != nil {
		WriteResult(
			http.StatusInternalServerError,
			fmt.Sprintf("Can not delete item, %s", deletedErr.Error()),
			c,
		)
		return
	}
	WriteResult(http.StatusOK, "Deleted record", c)
}
