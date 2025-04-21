package http

import (
	"context"
	"os"
	"pvz-backend/internal/models"
	"pvz-backend/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type DummyLoginRequest struct {
	Role string `json:"role"`
}

func DummyLogin(c *fiber.Ctx) error {
	var body DummyLoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "bad request"})
	}
	if body.Role != "moderator" && body.Role != "employee" {
		return c.Status(400).JSON(fiber.Map{"message": "invalid role"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": body.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "token error"})
	}

	return c.JSON(signed)
}

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "bad request"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	err := repository.CreateUser(context.Background(), req.Email, string(hash), req.Role)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "email already used"})
	}
	return c.SendStatus(201)
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "bad request"})
	}

	hash, role, err := repository.GetUserHashAndRole(context.Background(), req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return c.JSON(signed)
}