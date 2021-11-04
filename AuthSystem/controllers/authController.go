package controllers

import (
	"authSystem/database"
	"authSystem/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

// It will be used to check data in handlers.
var user models.User

func Register(c *fiber.Ctx) error {
	// The body data
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check if email already exists
	if user.Email == data["email"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "This email already exists",
		})
	}
	// Check if body-passwords are same
	if data["password"] != data["password_confirm"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}
	// Check password length
	if len(data["password"]) < 5 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Password must be minimum 5 characters",
		})
	}

	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user = models.User{
		Name:            data["name"],
		Email:           data["email"],
		Password:        passwordHashed,
		PasswordConfirm: passwordHashed,
	}

	database.DB.Create(&user)

	return c.JSON(fiber.Map{
		"message": "User registered"})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	database.DB.Where("email = ?", data["email"]).First(&user)

	// Check if user exists in database
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}
	// Check the password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	// Create JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	// Store tokenJWT to cookies
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "User logged",
	})

}

// Nos devuele el usuario conectado en ese momento
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	database.DB.Where("ID = ? ", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "User logged out",
	})
}
