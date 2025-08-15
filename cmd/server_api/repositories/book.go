package repositories

import (
	"book-management-system/cmd/server_api/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) Create(book *models.Book) error {
	now := time.Now().UTC()
	book.CreatedDate = now
	book.UpdatedDate = now
	return r.db.Create(book).Error
}

func (r *BookRepository) GetByID(id string) (*models.Book, error) {
	var book models.Book
	err := r.db.Where("id = ? AND deleted_date IS NULL", id).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) GetAll(limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("deleted_date IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) GetByStatus(status string, limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("status = ? AND deleted_date IS NULL", status).
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) GetByGenre(genre string, limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("genre = ? AND deleted_date IS NULL", genre).
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) GetByAuthor(author string, limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("LOWER(author) LIKE LOWER(?) AND deleted_date IS NULL", "%"+author+"%").
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) SearchByTitle(title string, limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("LOWER(title) LIKE LOWER(?) AND deleted_date IS NULL", "%"+title+"%").
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) SearchBooks(query string, limit, offset int) ([]models.Book, error) {
	var books []models.Book
	searchTerm := "%" + strings.ToLower(query) + "%"
	err := r.db.Where(
		"(LOWER(title) LIKE ? OR LOWER(author) LIKE ? OR LOWER(genre) LIKE ? OR isbn LIKE ?) AND deleted_date IS NULL",
		searchTerm, searchTerm, searchTerm, "%"+query+"%",
	).
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) GetAvailable(limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("available_quantity > 0 AND status = 'active' AND deleted_date IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&books).Error
	return books, err
}

func (r *BookRepository) Update(book *models.Book) error {
	book.UpdatedDate = time.Now().UTC()
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.Book{}).
		Where("id = ? AND deleted_date IS NULL", id).
		Update("deleted_date", now).Error
}

func (r *BookRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Book{}).Where("deleted_date IS NULL").Count(&count).Error
	return count, err
}

func (r *BookRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Book{}).
		Where("status = ? AND deleted_date IS NULL", status).
		Count(&count).Error
	return count, err
}

func (r *BookRepository) CountAvailable() (int64, error) {
	var count int64
	err := r.db.Model(&models.Book{}).
		Where("available_quantity > 0 AND status = 'active' AND deleted_date IS NULL").
		Count(&count).Error
	return count, err
}

func (r *BookRepository) ISBNExists(isbn string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Book{}).
		Where("isbn = ? AND deleted_date IS NULL", isbn).
		Count(&count).Error
	return count > 0, err
}

func (r *BookRepository) UpdateQuantity(id string, quantity, availableQuantity int) error {
	return r.db.Model(&models.Book{}).
		Where("id = ? AND deleted_date IS NULL", id).
		Updates(map[string]any{
			"quantity":           quantity,
			"available_quantity": availableQuantity,
			"updated_date":       time.Now().UTC(),
		}).Error
}