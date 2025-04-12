package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type RecipientController struct {
	interactor domain.DocumentRecipientInteractor
}

func NewRecipientController(interactor domain.DocumentRecipientInteractor) *RecipientController {
	return &RecipientController{interactor: interactor}
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
	userID, _ := ctx.Keys["userID"].(float64)
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}
	if err := c.interactor.SignDocument(ctx, idStr, uint(userID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject", "details": err.Error()})
		return
	}
}
