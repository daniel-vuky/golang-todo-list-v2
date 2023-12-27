package handler

import (
	"github.com/daniel-vuky/golang-todo-list-v2/model"
	"github.com/daniel-vuky/golang-todo-list-v2/repository"
	"github.com/gin-gonic/gin"
	ginSession "github.com/go-session/gin-session"
	"net/http"
	"regexp"
)

const EmailRegex string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

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

// GetUsernameFromContext retrieves the username from the JWT token
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
	re := regexp.MustCompile(EmailRegex)
	if len(password) < 6 || !re.MatchString(email) {
		Redirect("register", repository.InputErrorCode, c)
		return
	}
	if users.Auth.UserExisted(username) {
		Redirect("register", repository.UserExistedErrorCode, c)
		return
	}
	passwordHashed, hashedPasswordError := users.Auth.Hash(password)
	if hashedPasswordError != nil {
		Redirect("register", repository.ErrorEncounteredErrorCode, c)
		return
	}
	newUser := model.User{
		Username: username,
		Email:    email,
		Password: string(passwordHashed),
	}
	createUserErr := users.Repository.CreateNewUser(&newUser)
	if createUserErr != nil {
		Redirect("register", repository.ErrorEncounteredErrorCode, c)
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
	user := model.User{
		Username: username,
	}
	getUserErr := users.Repository.GetUser(&user)
	if getUserErr != nil || user.UserId == 0 {
		Redirect("login", repository.UsernamePasswordErrorCode, c)
		return
	}
	hashedPassword := user.Password
	if hashedError := users.Auth.ComparePasswordHash(hashedPassword, password); hashedError != nil {
		Redirect("login", repository.UsernamePasswordErrorCode, c)
		return
	}
	token, tokenErr := users.Auth.CreateToken(user.Username)
	if tokenErr != nil {
		Redirect("login", repository.ErrorEncounteredErrorCode, c)
		return
	}
	session := ginSession.FromContext(c)
	session.Set("token", token)
	session.Set("user_id", user.UserId)
	sessionErr := session.Save()
	if sessionErr != nil {
		Redirect("login", repository.ErrorEncounteredErrorCode, c)
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
