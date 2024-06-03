// socialmedia_handler.go

package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AhmadSafrizal/myGram/model"
	"github.com/AhmadSafrizal/myGram/repository"
	"github.com/gin-gonic/gin"
)

type SocialMediaHandler struct {
	Repository *repository.SocialMediaRepository
}

func NewSocialMediaHandler(repo *repository.SocialMediaRepository) *SocialMediaHandler {
	return &SocialMediaHandler{Repository: repo}
}

func (s *SocialMediaHandler) CreateSocialMedia(ctx *gin.Context) {
	socialMedia := &model.SocialMedia{}
	if err := ctx.Bind(socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body request"})
		return
	}

	if err := socialMedia.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set creation time
	socialMedia.CreatedAt = time.Now()

	if err := s.Repository.Create(socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not create social media"})
		return
	}

	ctx.JSON(http.StatusCreated, socialMedia)
}

func (s *SocialMediaHandler) GetSocialMedias(ctx *gin.Context) {
	socialMedias, err := s.Repository.Get()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, socialMedias)
}

func (s *SocialMediaHandler) UpdateSocialMedia(ctx *gin.Context) {
	var socialMedia model.SocialMedia
	if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body request"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid social media ID"})
		return
	}

	socialMedia.ID = uint(id)
	socialMedia.UpdatedAt = time.Now()

	if err := s.Repository.Update(&socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not update social media"})
		return
	}

	ctx.JSON(http.StatusOK, socialMedia)
}

func (s *SocialMediaHandler) DeleteSocialMedia(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid social media ID"})
		return
	}

	socialMedia := &model.SocialMedia{ID: uint(id)}

	if err := s.Repository.Delete(socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not delete social media"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "social media deleted"})
}
