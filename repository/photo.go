package repository

import (
	"log"

	"github.com/AhmadSafrizal/myGram/model"
	"gorm.io/gorm"
)

type PhotoRepository struct {
	DB *gorm.DB
}

func (p *PhotoRepository) Migrate() {
	err := p.DB.AutoMigrate(&model.Photo{})
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PhotoRepository) Create(photo *model.Photo) error {
	return p.DB.Debug().Create(photo).Error
}

func (p *PhotoRepository) Get() ([]model.Photo, error) {
	var photos []model.Photo
	err := p.DB.Debug().Find(&photos).Error
	return photos, err
}

func (p *PhotoRepository) GetById(photo *model.Photo, id uint) error {
	return p.DB.Debug().Model(&model.Photo{}).Where("id = ?", id).First(photo).Error
}

func (p *PhotoRepository) Update(photo *model.Photo) error {
	return p.DB.Debug().Save(photo).Error
}

func (p *PhotoRepository) Delete(photo *model.Photo) error {
	return p.DB.Debug().Delete(photo).Error
}
