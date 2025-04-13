package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
)

type RecipientController struct {
	interactor     domain.DocumentRecipientInteractor
	repo           domain.DocumentRecipientRepository
	redis          *redis.Client
	usrRepo        domain.UserRepository
	mailServiceURL string
}

func NewRecipientController(interactor domain.DocumentRecipientInteractor, repo domain.DocumentRecipientRepository, redis *redis.Client, usrRepo domain.UserRepository, mailServiceURL string) *RecipientController {
	return &RecipientController{interactor: interactor, repo: repo, redis: redis, usrRepo: usrRepo, mailServiceURL: mailServiceURL}
}

func (c *RecipientController) ListUserDocuments(ctx *gin.Context) {
	status := ctx.Query("status")
	userID, _ := ctx.Keys["userID"].(float64)
	list, err := c.interactor.ListUserDocuments(ctx, uint(userID), status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get documents", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}
func (c *RecipientController) RejectDocument(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, _ := ctx.Keys["userID"].(float64)
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}
	if err := c.interactor.RejectDocument(ctx, idStr, uint(userID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject", "details": err.Error()})
		return
	}
}

func (c *RecipientController) SignDocument(ctx *gin.Context) {
	idStr := ctx.Param("id")
	type SignDocumentRequest struct {
		Code string `json:"code" binding:"required"`
	}
	var req SignDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	usrID, _ := ctx.Keys["userID"].(float64)
	user, _ := c.usrRepo.User(ctx, uint(usrID))
	storedCode, err := c.redis.Get(ctx, "verification:"+user.Email).Result()

	if err == redis.Nil {
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Verification required", "details": err.Error()})
		return
	}

	if storedCode != req.Code {
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Invalid verification code"})
		return
	}
	userID, _ := ctx.Keys["userID"].(float64)
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing documentRecipient ID"})
		return
	}
	if err := c.interactor.SignDocument(ctx, idStr, uint(userID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject", "details": err.Error()})
		return
	}
	c.redis.Del(ctx, "verification:"+user.Email)
	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *RecipientController) ByID(ctx *gin.Context) {
	id := ctx.Param("id")

	recipient, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed get recipient`s document", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"body": recipient,
	})
}

// func (c *RecipientController) RequestToSign(ctx *gin.Context) {
// 	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
// 	code := r.Intn(10000) + 1000
// 	usrID, _ := ctx.Keys["userID"].(float64)
// 	user, _ := c.usrRepo.User(ctx, uint(usrID))
// 	err := c.redis.Set(ctx, "verification:"+user.Email, code, 10*time.Minute).Err()
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed redis set", "details": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"email": user.Email,
// 		"code":  code,
// 	})
// }

func (c *RecipientController) RequestToSign(ctx *gin.Context) {
	// Структура для тела запроса
	type EmailRequest struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	// idStr := ctx.Param("id")
	// id, _ := strconv.Atoi(idStr)
	userID, _ := ctx.Keys["userID"].(float64)
	user, _ := c.usrRepo.User(ctx, uint(userID))

	code := rand.Intn(10000) + 1000
	requestBody := EmailRequest{
		Email: user.Email,
		Code:  strconv.Itoa(code),
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to marshar request", "details": err.Error()})
		return
	}
	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Формируем запрос
	req, err := http.NewRequest("POST", c.mailServiceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create request", "details": err.Error()})
		return

	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Document-Service/1.0")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": err.Error()})
		return
	}
	c.redis.Set(ctx, "verification:"+user.Email, code, 10*time.Minute).Err()
	ctx.JSON(http.StatusOK, gin.H{"body": user})
}
