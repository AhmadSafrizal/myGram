package repository

import (
	"github.com/AhmadSafrizal/myGram/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *CommentRepository) Get() ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.DB.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) FindByID(comment *model.Comment, id uint) error {
	return r.DB.First(comment, id).Error
}

func (r *CommentRepository) Update(comment *model.Comment) error {
	return r.DB.Model(&model.Comment{}).Where("id = ?", comment.ID).Updates(comment).Error
}

func (r *CommentRepository) Delete(comment *model.Comment) error {
	return r.DB.Delete(comment).Error
}

func (r *CommentRepository) Migrate() {
	r.DB.AutoMigrate(&model.Comment{})
}
