package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/immxrtalbeast/flagman-backend/internal/lib"
)

type DocumentController struct {
	interactor domain.DocumentInteractor
	recipRepo  domain.DocumentRecipientRepository
}

// curl -X POST -F "document=@/home/codys/test.pdf " http://localhost:8080/api/v1/document/create -H "Authorization: Bearer
func NewDocumentController(interactor domain.DocumentInteractor, recipRepo domain.DocumentRecipientRepository) *DocumentController {
	return &DocumentController{interactor: interactor, recipRepo: recipRepo}
}

func (c *DocumentController) CreateDocument(ctx *gin.Context) {
	file, err := ctx.FormFile("document")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "No file uploaded",
			"details": err.Error()})
		return
	}

	// 2. Валидация типа файла
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}
	tempPath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
		return
	}
	defer os.Remove(tempPath)
	// 4. Загрузка в Supabase Storage
	publicURL, err := lib.UploadToSupabase(tempPath, file.Filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to storage", "details": err.Error()})
		return
	}
	senderID, _ := ctx.Keys["userID"].(float64)
	document, err := c.interactor.CreateDocument(ctx, uint(senderID), file.Filename, publicURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to db", "details": err.Error()})
		return
	}
	recipientsJSON := ctx.PostForm("recipients")
	var recipients []int
	if err := json.Unmarshal([]byte(recipientsJSON), &recipients); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipients format", "details": err.Error()})
		return
	}
	for _, recipientID := range recipients {
		recipient := domain.DocumentRecipient{
			DocumentID: document.ID,
			UserID:     uint(recipientID),
			Status:     "pending",
		}
		if err := c.recipRepo.Create(&recipient); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to recipient", "details": err.Error()})
			return
		}
	}

	// 5. Возврат публичного URL
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "File uploaded successfully",
		"url":        publicURL,
		"document":   document,
		"recipients": recipients,
	})

}
func (c *DocumentController) DocumentByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	document, err := c.interactor.DocumentByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get document", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"document": document,
	})
}
