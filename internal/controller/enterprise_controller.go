package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type EnterpriseController struct {
	interactor domain.EnterpriseInteractor
}

func NewEnterpriseController(interactor domain.EnterpriseInteractor) *EnterpriseController {
	return &EnterpriseController{interactor: interactor}
}

func (c *EnterpriseController) CreateEnterprise(ctx *gin.Context) {
	var request struct {
		Name        string `json:"name" binding:"required,min=3,max=50"`
		Description string `json:"description" binding :"max=150"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	userID, _ := ctx.Keys["userID"].(float64)

	// 3. Вызов бизнес-логики
	enterprise, err := c.interactor.CreateEnterprise(uint(userID), request.Name, request.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Успешный ответ
	ctx.JSON(http.StatusCreated, gin.H{
		"id":          enterprise.ID,
		"name":        enterprise.Name,
		"description": enterprise.Description,
	})
}

func (c *EnterpriseController) Enterprise(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	enterprise, err := c.interactor.EnterpriseByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get enterprise", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"body": enterprise,
	})
}
