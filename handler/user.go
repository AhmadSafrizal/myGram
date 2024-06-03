package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AhmadSafrizal/myGram/helper"
	"github.com/AhmadSafrizal/myGram/model"
	"github.com/AhmadSafrizal/myGram/repository"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repository *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repository: repo}
}

func (u *UserHandler) GetGorm(ctx *gin.Context) {
	users, err := u.Repository.Get()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (u *UserHandler) CreateGorm(ctx *gin.Context) {
	user := &model.User{}
	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid body request",
		})
		return
	}

	if !helper.IsValidEmail(user.Email) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid email format",
		})
		return
	}

	if len(user.Password) < 6 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "password must be at least 6 characters",
		})
		return
	}

	userFetched, err := u.Repository.GetByEmail(user.Email)
	if err == nil && userFetched.ID != 0 {
		ctx.AbortWithStatusJSON(http.StatusConflict, map[string]any{
			"message": "email already registered",
		})
		return
	}

	passwordHashed, err := helper.HashPassword(user.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "something went wrong",
		})
		return
	}
	user.Password = passwordHashed

	err = u.Repository.Create(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "something went wrong",
		})
		return
	}

	res := map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"age":      user.Age,
	}
	ctx.JSON(http.StatusCreated, res)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	user := &model.User{}
	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid body request",
		})
		return
	}

	userFetched, err := u.Repository.GetByEmail(user.Email)
	if err != nil || userFetched.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{
			"message": "email not found",
		})
		return
	}
	valid := helper.CheckPasswordHash(user.Password, userFetched.Password)
	if !valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
			"message": "wrong password",
		})
		return
	}
	token, err := helper.GenerateUserJWT(userFetched.Username, userFetched.Email)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

func (u *UserHandler) UpdateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid body request",
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid user ID",
		})
		return
	}

	user.ID = uint(id)

	existingUser := model.User{}
	if err := u.Repository.GetById(&existingUser, user.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
			"message": "user not found",
		})
		return
	}

	if user.Password != "" {
		passwordHashed, err := helper.HashPassword(user.Password)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "could not hash password",
			})
			return
		}
		user.Password = passwordHashed
	} else {
		user.Password = existingUser.Password
	}

	user.UpdatedAt = time.Now()

	if err := u.Repository.Update(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not update user",
		})
		return
	}

	res := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"age":        user.Age,
		"updated_at": user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, res)
}

func (u *UserHandler) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "invalid user ID",
		})
		return
	}

	user := &model.User{ID: uint(id)}

	if err := u.Repository.Delete(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "could not delete user",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "Account succcs successfully deleted",
	})
}
