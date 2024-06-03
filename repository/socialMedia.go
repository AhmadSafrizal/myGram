package repository

import (
	"log"

	"github.com/AhmadSafrizal/myGram/model"
	"gorm.io/gorm"
)

type SocialMediaRepository struct {
	DB *gorm.DB
}

func (s *SocialMediaRepository) Migrate() {
	err := s.DB.AutoMigrate(&model.SocialMedia{})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SocialMediaRepository) Create(socialMedia *model.SocialMedia) error {
	err := s.DB.Debug().Model(&model.SocialMedia{}).Create(socialMedia).Error
	return err
}

func (s *SocialMediaRepository) Get() ([]*model.SocialMedia, error) {
	socialMedia := []*model.SocialMedia{}
	err := s.DB.Debug().Model(&model.SocialMedia{}).Find(&socialMedia).Error
	return socialMedia, err
}

func (s *SocialMediaRepository) Update(socialMedia *model.SocialMedia) error {
	err := s.DB.Debug().Model(&model.SocialMedia{}).Save(socialMedia).Error
	return err
}

func (s *SocialMediaRepository) Delete(socialMedia *model.SocialMedia) error {
	err := s.DB.Debug().Model(&model.SocialMedia{}).Delete(socialMedia).Error
	return err
}
