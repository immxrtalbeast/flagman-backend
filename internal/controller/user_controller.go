package controller

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/immxrtalbeast/flagman-backend/internal/lib"
)

type UserController struct {
	interactor  domain.UserInteractor
	tokenTTL    time.Duration
	tokenSecret string
}

func NewUserController(interactor domain.UserInteractor, tokenTTL time.Duration, tokenSecret string) *UserController {
	return &UserController{interactor: interactor, tokenTTL: tokenTTL, tokenSecret: tokenSecret}
}

func (c *UserController) Register(ctx *gin.Context) {
	type RegisterRequest struct {
		FullName    string `json:"fullname" binding:"required,min=3,max=50"`
		Email       string `json:"email" binding:"required"`
		PhoneNumber string `json:"phonenumber" binding:"required"`
		Pass        string `json:"password" binding:"required,min=8,max=50"`
	}

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Валидация имени
	nameRegex := regexp.MustCompile(`^[а-яА-ЯёЁ\s]+$`)
	if !nameRegex.MatchString(req.FullName) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid name",
			"details": "Name must contain only russian letters and spaces",
		})
		return
	}
	// Валидация email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid email",
			"details": "Email must be in a valid format (e.g. user@example.com)",
		})
		return
	}

	// Валидация номера телефона (+7 и 10 цифр)
	phoneRegex := regexp.MustCompile(`^\+[78]\d{10}$`)
	if !phoneRegex.MatchString(req.PhoneNumber) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid phone number",
			"details": "Phone must be in format +7XXXXXXXXXX (11 digits after +7)",
		})
		return
	}
	// Валидация пароля
	passRegex := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\[\]{};:<>,./?~\\-]+$`)
	if !passRegex.MatchString(req.Pass) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid password",
			"details": "Password contains forbidden characters",
		})
		return
	}

	// Если все проверки пройдены
	id, err := c.interactor.CreateUser(ctx, req.FullName, req.Email, req.PhoneNumber, req.Pass)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"details": err.Error(),
		})
		return
	}
	token, err := c.interactor.Login(ctx, req.Email, req.Pass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to login",
			"details": err.Error(),
		})
		return
	}
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"jwt",                     // Имя куки
		token,                     // Значение токена
		int(c.tokenTTL.Seconds()), // Макс возраст в секундах
		"/",                       // Путь
		"",                        // Домен (пусто для текущего домена)
		true,                      // Secure (использовать true в production для HTTPS)
		true,                      // HttpOnly
	)

	userIDStr := strconv.FormatUint(uint64(id), 10)

	// Устанавливаем куку
	ctx.SetCookie(
		"user_id",
		userIDStr,
		int(c.tokenTTL.Seconds()),
		"/",
		"",
		true,
		false, // HttpOnly=false чтобы клиент мог читать JS
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created",
	})

}

func (c *UserController) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Email string `json:"email" binding:"required"`
		Pass  string `json:"password" binding:"required"`
	}
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	token, err := c.interactor.Login(ctx, req.Email, req.Pass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to login",
			"details": err.Error(),
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"jwt",                     // Имя куки
		token,                     // Значение токена
		int(c.tokenTTL.Seconds()), // Макс возраст в секундах
		"/",                       // Путь
		"",                        // Домен (пусто для текущего домена)
		true,                      // Secure (использовать true в production для HTTPS)
		false,                     // HttpOnly
	)
	id, err := lib.IdFromToken(token, c.tokenSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	userIDStr := strconv.FormatUint(uint64(id), 10)
	ctx.SetCookie(
		"user_id",
		userIDStr,
		int(c.tokenTTL.Seconds()),
		"/",
		"",
		true,
		false, // HttpOnly=false чтобы клиент мог читать JS
	)

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *UserController) User(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}
	id, _ := strconv.Atoi(idStr)
	user, err := c.interactor.User(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get user",
			"details": err.Error(),
		})
		return
	}

	type UserResponse struct {
		ID            uint        `json:"ID"`
		FullName      string      `json:"FullName"`
		Email         string      `json:"Email"`
		PhoneNumber   string      `json:"PhoneNumber"`
		CreatedAt     time.Time   `json:"CreatedAt"`
		Organizations interface{} `json:"Organizations"`
		Roles         interface{} `json:"Roles"`
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": UserResponse{
			ID:            user.ID,
			FullName:      user.FullName,
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			CreatedAt:     user.CreatedAt,
			Organizations: user.Organizations,
			Roles:         user.Roles,
		},
	})
}
func (c *UserController) Users(ctx *gin.Context) {
	usersDB, err := c.interactor.Users(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get users",
			"details": err.Error(),
		})
		return
	}
	type UsersResponse struct {
		ID       uint   `json:"ID"`
		FullName string `json:"FullName"`
	}
	var users []UsersResponse
	for _, user := range usersDB {
		user := UsersResponse{
			ID:       user.ID,
			FullName: user.FullName,
		}
		users = append(users, user)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
