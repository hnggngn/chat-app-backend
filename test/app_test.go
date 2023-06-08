package test

import (
	"bytes"
	"chat_backend/generated"
	"chat_backend/internal/delivery/router"
	"chat_backend/pkg/utils"
	"context"
	"encoding/json"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/validate"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func appTest() (*fiber.App, *generated.Queries) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error reading environment: %v", err)
	}

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"),
		AllowCredentials: true,
	}))
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(recover.New())

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	_, queries := utils.Database()

	cld, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	router.AppRouter(app, queries, cld)

	return app, queries
}

type TestCase struct {
	expected interface{}
	actual   interface{}
}

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	username       = "test-user"
	updateUsername = "test-user-updated"
	password       = "test-password"
)

func genValue(l int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, l)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func afterAll() func(t *testing.T) {
	_, queries := appTest()
	_ = queries.DeleteUserByUsername(context.Background(), username)
	_ = queries.DeleteUserByUsername(context.Background(), updateUsername)
	return func(t *testing.T) {
		t.Log("Clean up.")
	}
}

var (
	app, _ = appTest()
)

func TestSignUp(t *testing.T) {
	defer afterAll()

	t.Run("Should return error when username and password is empty", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": "",
			"password": "",
		}

		errorSchema := map[string]map[string]string{
			"password": {
				"required": "password is required to not be empty",
			},
			"username": {
				"required": "username is required to not be empty",
			},
		}

		input, _ := json.Marshal(inputSchema)
		expected, _ := json.Marshal(errorSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusForbidden,
				actual:   res.StatusCode,
			},
			{
				expected: string(expected),
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})

	t.Run("Should return error when username and password exceeds character length", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": genValue(40),
			"password": genValue(200),
		}

		errorSchema := map[string]map[string]string{
			"password": {
				"max_len": "password max length is 100",
			},
			"username": {
				"max_len": "username max length is 30",
			},
		}

		input, _ := json.Marshal(inputSchema)
		expected, _ := json.Marshal(errorSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusForbidden,
				actual:   res.StatusCode,
			},
			{
				expected: string(expected),
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})

	t.Run("Should return CREATED", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusCreated,
				actual:   res.StatusCode,
			},
			{
				expected: "Created",
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})

	t.Run("Should return error when user already exists", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		errorSchema := fiber.Map{
			"message": "User already exists.",
		}

		input, _ := json.Marshal(inputSchema)
		expected, _ := json.Marshal(errorSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusForbidden,
				actual:   res.StatusCode,
			},
			{
				expected: string(expected),
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})
}

func TestLogin(t *testing.T) {
	defer afterAll()

	t.Run("Should return error when username and password is empty", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": "",
			"password": "",
		}

		errorSchema := map[string]map[string]string{
			"password": {
				"required": "password is required to not be empty",
			},
			"username": {
				"required": "username is required to not be empty",
			},
		}

		input, _ := json.Marshal(inputSchema)
		expected, _ := json.Marshal(errorSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusForbidden,
				actual:   res.StatusCode,
			},
			{
				expected: string(expected),
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})

	t.Run("Should return error when user not exists", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": "unknown",
			"password": password,
		}

		errorSchema := fiber.Map{
			"message": "User not exists.",
		}

		input, _ := json.Marshal(inputSchema)
		expected, _ := json.Marshal(errorSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusForbidden,
				actual:   res.StatusCode,
			},
			{
				expected: string(expected),
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})

	t.Run("Should return OK", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		signUpReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		signUpReq.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(signUpReq)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusOK,
				actual:   res.StatusCode,
			},
			{
				expected: "OK",
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
		assert.NotEmpty(t, res.Cookies())
	})

	t.Run("Should return error when wrong password", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": "wrong-password",
		}

		errorSchema := fiber.Map{
			"message": "Password not correct.",
		}

		input, _ := json.Marshal(inputSchema)
		expected, _ := json.Marshal(errorSchema)

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		tests := []TestCase{
			{
				expected: fiber.StatusForbidden,
				actual:   res.StatusCode,
			},
			{
				expected: string(expected),
				actual:   string(body),
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.actual)
		}
	})
}

func TestLogout(t *testing.T) {
	defer afterAll()

	t.Run("Should return error when not logged", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/signout", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Should log out", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		signUpReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		signUpReq.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(signUpReq)

		loginReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRes, _ := app.Test(loginReq)

		cookie := loginRes.Cookies()

		req := httptest.NewRequest(fiber.MethodPost, "/api/auth/signout", nil)
		req.AddCookie(cookie[0])
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

func TestGetProfile(t *testing.T) {
	defer afterAll()

	t.Run("Should return error when not logged", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, "/api/user/profile", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Should return data", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		signUpReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		signUpReq.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(signUpReq)

		loginReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRes, _ := app.Test(loginReq)

		cookie := loginRes.Cookies()

		req := httptest.NewRequest(fiber.MethodGet, "/api/user/profile", nil)
		req.AddCookie(cookie[0])
		res, _ := app.Test(req)

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.NotEmpty(t, body)
	})
}

func TestUpdateProfile(t *testing.T) {
	defer afterAll()

	t.Run("Should return error when not logged", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPatch, "/api/user/profile/update", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Should update 1 field", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		signUpReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		signUpReq.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(signUpReq)

		loginReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRes, _ := app.Test(loginReq)

		cookie := loginRes.Cookies()

		inputUpdateSchema := fiber.Map{
			"username": updateUsername,
		}

		inputUpdate, _ := json.Marshal(inputUpdateSchema)

		req := httptest.NewRequest(fiber.MethodPatch, "/api/user/profile/update", bytes.NewReader(inputUpdate))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(cookie[0])
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("Should update 2 field", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		signUpReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		signUpReq.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(signUpReq)

		loginReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRes, _ := app.Test(loginReq)

		cookie := loginRes.Cookies()

		inputUpdateSchema := fiber.Map{
			"username": updateUsername,
			"password": "blah",
		}

		inputUpdate, _ := json.Marshal(inputUpdateSchema)

		req := httptest.NewRequest(fiber.MethodPatch, "/api/user/profile/update", bytes.NewReader(inputUpdate))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(cookie[0])
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

func TestDeleteUser(t *testing.T) {
	defer afterAll()

	t.Run("Should return error when not logged", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodDelete, "/api/user/profile/delete", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Should delete user data and log out", func(t *testing.T) {
		inputSchema := fiber.Map{
			"username": username,
			"password": password,
		}

		input, _ := json.Marshal(inputSchema)

		signUpReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/signup", bytes.NewReader(input))
		signUpReq.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(signUpReq)

		loginReq := httptest.NewRequest(fiber.MethodPost, "/api/auth/login", bytes.NewReader(input))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRes, _ := app.Test(loginReq)

		cookie := loginRes.Cookies()

		req := httptest.NewRequest(fiber.MethodDelete, "/api/user/profile/delete", nil)
		req.AddCookie(cookie[0])
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Empty(t, res.Cookies()[0].Value)
	})
}
