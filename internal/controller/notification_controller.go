package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type NotificationController struct {
	interactor domain.NotificationInteractor
}

func NewNotificationController(interactor domain.NotificationInteractor) *NotificationController {
	return &NotificationController{interactor: interactor}
}

func (c *NotificationController) MyNotifications(ctx *gin.Context) {
	userID, _ := ctx.Keys["userID"].(float64)
	invitations, err := c.interactor.MyNotifications(ctx, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"body": invitations,
	})
}
func (c *NotificationController) AcceptInvite(ctx *gin.Context) {
	type AcceptInviteRequest struct {
		InvitationID uint `json:"invitation_id" binding:"required"`
		EnterpriseID uint `json:"enterprise_id" binding:"required"`
	}
	var req AcceptInviteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	userID, _ := ctx.Keys["userID"].(float64)
	err := c.interactor.AcceptInvite(ctx, uint(userID), req.InvitationID, req.EnterpriseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
