package repository

import (
	"database/sql"
	"github.com/daniel-vuky/golang-todo-list-v2/auth"
	"github.com/daniel-vuky/golang-todo-list-v2/model"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const MissingInputErrorCode int = 1
const MissingInputError string = "Please fill all the information input"
const InputErrorCode int = 2
const InputError string = "Please check the input"
const UserExistedErrorCode int = 3
const UserExistedError string = "An account with this username or email was existed"
const UsernamePasswordErrorCode int = 4
const UsernamePasswordError string = "Username or Password is not correct"
const CreateUserErrorCode int = 5
const CreateUserError string = "Fail to create new user, %s"
const ErrorEncounteredErrorCode int = 6
const ErrorEncounteredError = "An error encountered white performing operations"

type AuthRepository struct {
	Db *sql.DB
}

func (authRepository AuthRepository) UserExisted(value string) bool {
	var user model.User
	authRepository.Db.QueryRow("SELECT user_id FROM users WHERE username = ?", value).Scan(
		&user.UserId,
	)

	return user.UserId != 0
}

// Hash encrypt the password
func (authRepository AuthRepository) Hash(password string) ([]byte, error) {
	bytes, encryptError := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, encryptError
}

// ComparePasswordHash compare hashed password and input password
func (authRepository AuthRepository) ComparePasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CreateToken Create a token base on username
func (authRepository AuthRepository) CreateToken(username string) (string, error) {
	return auth.Create(username)
}

// ParseToken Parse the token
func (authRepository AuthRepository) ParseToken(token string) (*jwtGo.Token, error) {
	return auth.ValidateToken(token)
}

// GetUsernameFromToken Parse the token
func (authRepository AuthRepository) GetUsernameFromToken(token string) (string, error) {
	return auth.GetUsernameFromToken(token)
}

// GetUserIDFromToken Parse the token
func (authRepository AuthRepository) GetUserIDFromToken(token string) (uint64, error) {
	return auth.GetUserIDFromToken(token)
}

// GetErrorMessageByCode Get error message by code
func (authRepository AuthRepository) GetErrorMessageByCode(code int) string {
	mappingError := map[int]string{
		MissingInputErrorCode:     MissingInputError,
		InputErrorCode:            InputError,
		UserExistedErrorCode:      UserExistedError,
		UsernamePasswordErrorCode: UsernamePasswordError,
		CreateUserErrorCode:       CreateUserError,
		ErrorEncounteredErrorCode: ErrorEncounteredError,
	}
	errorMessage, existed := mappingError[code]
	if !existed {
		return ErrorEncounteredError
	}
	return errorMessage
}
