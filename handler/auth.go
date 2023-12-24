package handler

import (
	"fmt"
	"github.com/daniel-vuky/golang-todo-list-and-chat/model"
	"github.com/daniel-vuky/golang-todo-list-and-chat/repository"
	"github.com/gin-gonic/gin"
	ginSession "github.com/go-session/gin-session"
	"net/http"
	"regexp"
)

const EmailRegex string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

const MissingInputErr = "Please fill all the information input"
const InputErr = "Please check the input"
const UserExisted = "An account with this username or email was existed"
const UsernamePasswordErr = "Username or Password is not correct"
const CreateUserErr = "Fail to create new user, %s"

type Users struct {
	Repository *repository.UsersRepository
	Auth       *repository.AuthRepository
}

func (users Users) AuthMiddleware(c *gin.Context) {
	session := ginSession.FromContext(c)
	if session == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	tokenString, tokenFine := session.Get("token")

	if tokenString == nil || !tokenFine {
		// Redirect to login page if not authenticated
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	token, err := users.Auth.ParseToken(tokenString.(string))

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	c.Next()
}

// GetUsernameFromToken retrieves the username from the JWT token
func (users Users) GetUsernameFromContext(c *gin.Context) string {
	session := ginSession.FromContext(c)
	tokenString, tokenFine := session.Get("token")

	if tokenString == nil || !tokenFine {
		return ""
	}

	username, err := users.Auth.GetUsernameFromToken(tokenString.(string))
	if err != nil {
		return ""
	}

	return username
}

// Register user
func (users Users) Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		WriteResult(http.StatusBadRequest, MissingInputErr, c)
		return
	}
	re := regexp.MustCompile(EmailRegex)
	if len(password) < 6 || !re.MatchString(email) {
		WriteResult(http.StatusBadRequest, InputErr, c)
		return
	}
	if users.Auth.UserExisted(username) {
		WriteResult(http.StatusBadRequest, UserExisted, c)
		return
	}
	passwordHashed, hashedPasswordError := users.Auth.Hash(password)
	if hashedPasswordError != nil {
		WriteResult(http.StatusBadRequest, hashedPasswordError.Error(), c)
		return
	}
	newUser := model.User{
		Username: username,
		Email:    email,
		Password: string(passwordHashed),
	}
	createUserErr := users.Repository.CreateNewUser(&newUser)
	if createUserErr != nil {
		WriteResult(http.StatusInternalServerError, fmt.Sprintf(CreateUserErr, createUserErr.Error()), c)
		return
	}
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
	return
}

// Login Post Login
func (users Users) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) == 0 || len(password) == 0 {
		WriteResult(http.StatusBadRequest, MissingInputErr, c)
		return
	}
	user := model.User{
		Username: username,
	}
	getUserErr := users.Repository.GetUser(&user)
	if getUserErr != nil || user.UserId == 0 {
		WriteResult(http.StatusBadRequest, UsernamePasswordErr, c)
		return
	}
	hashedPassword := user.Password
	if hashedError := users.Auth.ComparePasswordHash(hashedPassword, password); hashedError != nil {
		WriteResult(http.StatusBadRequest, UsernamePasswordErr, c)
		return
	}
	token, tokenErr := users.Auth.CreateToken(user.Username)
	if tokenErr != nil {
		WriteResult(http.StatusBadRequest, "Can not create the customer token", c)
		return
	}
	session := ginSession.FromContext(c)
	session.Set("token", token)
	session.Set("user_id", user.UserId)
	sessionErr := session.Save()
	if sessionErr != nil {
		c.AbortWithError(http.StatusInternalServerError, sessionErr)
		return
	}
	c.Redirect(http.StatusFound, "/")
	c.Abort()
	return
}

// Logout Post Login
func (users Users) Logout(c *gin.Context) {
	session := ginSession.FromContext(c)
	session.Delete("token")
	session.Save()
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
	return
}
