package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AhmadSafrizal/myGram/model"
	"github.com/AhmadSafrizal/myGram/repository"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Repository *repository.CommentRepository
}

func NewCommentHandler(repo *repository.CommentRepository) *CommentHandler {
	return &CommentHandler{Repository: repo}
}

func (c *CommentHandler) CreateComment(ctx *gin.Context) {
	comment := &model.Comment{}
	if err := ctx.Bind(comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body request"})
		return
	}

	if err := comment.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set creation time
	comment.CreatedAt = time.Now()

	if err := c.Repository.Create(comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not create comment"})
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

func (c *CommentHandler) GetComments(ctx *gin.Context) {
	comments, err := c.Repository.Get()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

func (c *CommentHandler) UpdateComment(ctx *gin.Context) {
	var comment model.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body request"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid comment ID"})
		return
	}

	comment.ID = uint(id)
	comment.UpdatedAt = time.Now()

	var existingComment model.Comment
	if err := c.Repository.FindByID(&existingComment, uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not find comment"})
		return
	}

	comment.CreatedAt = existingComment.CreatedAt

	fmt.Printf("Updating comment: %+v\n", comment)

	if err := c.Repository.Update(&comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not update comment"})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func (c *CommentHandler) DeleteComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid comment ID"})
		return
	}

	comment := &model.Comment{ID: uint(id)}

	if err := c.Repository.Delete(comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "could not delete comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
