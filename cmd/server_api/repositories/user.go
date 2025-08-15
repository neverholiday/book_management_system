package repositories

import (
	"book-management-system/cmd/server_api/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	now := time.Now().UTC()
	user.CreatedDate = now
	user.UpdatedDate = now
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND deleted_date IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_date IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("deleted_date IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByRole(role string, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ? AND deleted_date IS NULL", role).
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByStatus(status string, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("status = ? AND deleted_date IS NULL", status).
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	user.UpdatedDate = time.Now().UTC()
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.User{}).
		Where("id = ? AND deleted_date IS NULL", id).
		Update("deleted_date", now).Error
}

func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("deleted_date IS NULL").Count(&count).Error
	return count, err
}

func (r *UserRepository) CountByRole(role string) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("role = ? AND deleted_date IS NULL", role).
		Count(&count).Error
	return count, err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ? AND deleted_date IS NULL", email).
		Count(&count).Error
	return count > 0, err
}