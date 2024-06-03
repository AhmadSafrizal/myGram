package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AhmadSafrizal/myGram/model"
	"github.com/AhmadSafrizal/myGram/repository"
	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	Repository *repository.PhotoRepository
}

func NewPhotoHandler(repo *repository.PhotoRepository) *PhotoHandler {
	return &PhotoHandler{Repository: repo}
}

func (h *PhotoHandler) AddPhoto(ctx *gin.Context) {
	title := ctx.PostForm("title")
	caption := ctx.PostForm("caption")

	file, err := ctx.FormFile("photo")
	if err != nil {
		log.Println("Error getting form file:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid file"})
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

	if err := ctx.SaveUploadedFile(file, filepath.Join("uploads", filename)); err != nil {
		log.Println("Error saving uploaded file:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not save file"})
		return
	}

	photoURL := fmt.Sprintf("/uploads/%s", filename)

	photo := model.Photo{
		Title:    title,
		Caption:  caption,
		PhotoURL: photoURL,
	}
	if err := h.Repository.Create(&photo); err != nil {
		log.Println("Error creating photo record:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not save photo information"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "file uploaded successfully",
		"id":        photo.ID,
		"title":     photo.Title,
		"caption":   photo.Caption,
		"url_photo": photo.PhotoURL,
	})
}

func (p *PhotoHandler) GetPhotos(ctx *gin.Context) {
	photos, err := p.Repository.Get()
	if err != nil {
		log.Println("Error getting photos:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, photos)
}

func (p *PhotoHandler) UpdatePhoto(ctx *gin.Context) {
	var photo model.Photo
	if err := ctx.ShouldBind(&photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid body request",
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println("Error parsing photo ID:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid photo ID",
		})
		return
	}

	photo.ID = uint(id)

	existingPhoto := model.Photo{}
	if err := p.Repository.GetById(&existingPhoto, photo.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{
			"message": "photo not found",
		})
		return
	}

	file, err := ctx.FormFile("photo")
	if err == nil {
		oldFilePath := filepath.Join("uploads", filepath.Base(existingPhoto.PhotoURL))
		if _, err := os.Stat(oldFilePath); err == nil {
			os.Remove(oldFilePath)
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
		if err := ctx.SaveUploadedFile(file, filepath.Join("uploads", filename)); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": "could not save file",
			})
			return
		}
		photo.PhotoURL = fmt.Sprintf("/uploads/%s", filename)
	} else {
		photo.PhotoURL = existingPhoto.PhotoURL
	}

	if err := p.Repository.Update(&photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "could not update photo",
		})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (p *PhotoHandler) DeletePhoto(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println("Error parsing photo ID:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid photo ID",
		})
		return
	}

	photo := &model.Photo{ID: uint(id)}

	existingPhoto := model.Photo{}
	if err := p.Repository.GetById(&existingPhoto, photo.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{
			"message": "photo not found",
		})
		return
	}

	if err := p.Repository.Delete(photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "could not delete photo",
		})
		return
	}

	filePath := filepath.Join("uploads", filepath.Base(existingPhoto.PhotoURL))
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "photo deleted",
	})
}
