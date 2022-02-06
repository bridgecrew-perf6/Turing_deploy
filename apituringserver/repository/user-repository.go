package repository

import (
	"log"

	"github.com/xavimg/Turing/apituringserver/dto"
	"github.com/xavimg/Turing/apituringserver/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository is a contract what UserRepository can do to db.
type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User, userID string, newInfo dto.UserUpdateDTO) entity.User
	VerifyCredential(email, password string) interface{}
	VerifyUserExist(id string) interface{}
	IsDuplicateEmail(email string) (ctx *gorm.DB)
	FindByEmail(username string) entity.User
	ProfileUser(userID string) entity.User
	SaveToken(user entity.User, token string)
	DeleteToken(user entity.User, token string)
	GetToken(userID string) entity.User
}

type userConnection struct {
	connection *gorm.DB
}

// NewUserRepository is creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))

	db.connection.Save(&user)
	db.connection.Preload("Characters").Find(&user)

	return user
}

func (db *userConnection) UpdateUser(user entity.User, userID string, newInfo dto.UserUpdateDTO) entity.User {

	if newInfo.Name != "" {

		db.connection.Model(user).Where("id = ?", userID).Update("name", newInfo.Name)
	}

	if newInfo.Email != "" {

		db.connection.Model(user).Where("id = ?", userID).Update("email", newInfo.Email)
	}

	if newInfo.Password != "" {

		user.Password = hashAndSalt([]byte(newInfo.Password))

		db.connection.Model(user).Where("id = ?", userID).Update("password", user.Password)
	}

	db.connection.Preload("Characters").Preload("Characters.User").Find(&user)

	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}

	return string(hash)
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User

	res := db.connection.Where("email = ?", email).Take(&user)

	if res == nil {
		return res.Error
	}
	return user
}

func (db *userConnection) VerifyUserExist(id string) interface{} {
	var user entity.User

	res := db.connection.Where("id = ?", id).Take(&user)

	if res == nil {
		return res.Error
	}
	return user
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User

	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User

	db.connection.Find(&user, userID)

	return user
}

func (db *userConnection) FindByEmail(username string) entity.User {
	var user entity.User

	db.connection.Where("email = ? ", username).Take(&user)

	return user
}

func (db *userConnection) SaveToken(user entity.User, token string) {

	user.Token = token

	db.connection.Save(&user)
}

func (db *userConnection) DeleteToken(user entity.User, s string) {

	user.Token = s

	db.connection.Save(&user)

}

func (db *userConnection) GetToken(UserID string) entity.User {
	var user entity.User

	db.connection.Find(&user, UserID)

	return user
}
