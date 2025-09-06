package apis

import (
	"book-management-system/cmd/server_api/models"
	"book-management-system/cmd/server_api/repositories"
	"book-management-system/pkg/auth"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookAPI struct {
	bookRepo *repositories.BookRepository
	authMw   *auth.Middleware
}

func NewBookAPI(bookRepo *repositories.BookRepository, authMw *auth.Middleware) *BookAPI {
	return &BookAPI{
		bookRepo: bookRepo,
		authMw:   authMw,
	}
}

func (api *BookAPI) Setup(group *echo.Group) {
	group.POST("", api.createBook, api.authMw.RequireAdmin())
	group.GET("", api.getBooks)
	group.GET("/:id", api.getBook)
	group.GET("/search", api.searchBooks)
	group.GET("/available", api.getAvailableBooks)
	group.PUT("/:id", api.updateBook, api.authMw.RequireAdmin())
	group.DELETE("/:id", api.deleteBook, api.authMw.RequireAdmin())
	group.PUT("/:id/quantity", api.updateQuantity, api.authMw.RequireAdmin())
}

func (api *BookAPI) createBook(c echo.Context) error {
	var req struct {
		Title             string   `json:"title"`
		Author            string   `json:"author"`
		ISBN              *string  `json:"isbn"`
		Publisher         *string  `json:"publisher"`
		PublicationYear   *int     `json:"publication_year"`
		Genre             *string  `json:"genre"`
		Description       *string  `json:"description"`
		Pages             *int     `json:"pages"`
		Language          string   `json:"language"`
		Price             *float64 `json:"price"`
		Quantity          int      `json:"quantity"`
		AvailableQuantity int      `json:"available_quantity"`
		Location          *string  `json:"location"`
		Status            string   `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request body",
		})
	}

	if req.Title == "" || req.Author == "" || req.Language == "" || req.Status == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Title, author, language, and status are required",
		})
	}

	if req.ISBN != nil && *req.ISBN != "" {
		exists, err := api.bookRepo.ISBNExists(*req.ISBN)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Message: "Failed to check ISBN existence",
			})
		}
		if exists {
			return c.JSON(http.StatusConflict, models.Response{
				Message: "Book with this ISBN already exists",
			})
		}
	}

	book := &models.Book{
		ID:                uuid.New().String(),
		Title:             req.Title,
		Author:            req.Author,
		ISBN:              req.ISBN,
		Publisher:         req.Publisher,
		PublicationYear:   req.PublicationYear,
		Genre:             req.Genre,
		Description:       req.Description,
		Pages:             req.Pages,
		Language:          req.Language,
		Price:             req.Price,
		Quantity:          req.Quantity,
		AvailableQuantity: req.AvailableQuantity,
		Location:          req.Location,
		Status:            req.Status,
	}

	if err := api.bookRepo.Create(book); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to create book",
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
		Data:    book,
		Message: "Book created successfully",
	})
}

func (api *BookAPI) getBooks(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	status := c.QueryParam("status")
	genre := c.QueryParam("genre")
	author := c.QueryParam("author")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var books []models.Book
	var err error

	if status != "" {
		books, err = api.bookRepo.GetByStatus(status, limit, offset)
	} else if genre != "" {
		books, err = api.bookRepo.GetByGenre(genre, limit, offset)
	} else if author != "" {
		books, err = api.bookRepo.GetByAuthor(author, limit, offset)
	} else {
		books, err = api.bookRepo.GetAll(limit, offset)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to retrieve books",
		})
	}

	total, err := api.bookRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get book count",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data: map[string]any{
			"books":  books,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
		Message: "Books retrieved successfully",
	})
}

func (api *BookAPI) getBook(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Book ID is required",
		})
	}

	book, err := api.bookRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Message: "Book not found",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    book,
		Message: "Book retrieved successfully",
	})
}

func (api *BookAPI) searchBooks(c echo.Context) error {
	query := c.QueryParam("q")
	title := c.QueryParam("title")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	if query == "" && title == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Search query (q) or title parameter is required",
		})
	}

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var books []models.Book
	var err error

	if title != "" {
		books, err = api.bookRepo.SearchByTitle(title, limit, offset)
	} else {
		books, err = api.bookRepo.SearchBooks(query, limit, offset)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to search books",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data: map[string]any{
			"books":  books,
			"query":  query,
			"title":  title,
			"limit":  limit,
			"offset": offset,
		},
		Message: "Books search completed successfully",
	})
}

func (api *BookAPI) getAvailableBooks(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	books, err := api.bookRepo.GetAvailable(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to retrieve available books",
		})
	}

	count, err := api.bookRepo.CountAvailable()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get available book count",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data: map[string]any{
			"books":  books,
			"total":  count,
			"limit":  limit,
			"offset": offset,
		},
		Message: "Available books retrieved successfully",
	})
}

func (api *BookAPI) updateBook(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Book ID is required",
		})
	}

	book, err := api.bookRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Message: "Book not found",
		})
	}

	var req struct {
		Title             *string  `json:"title"`
		Author            *string  `json:"author"`
		ISBN              *string  `json:"isbn"`
		Publisher         *string  `json:"publisher"`
		PublicationYear   *int     `json:"publication_year"`
		Genre             *string  `json:"genre"`
		Description       *string  `json:"description"`
		Pages             *int     `json:"pages"`
		Language          *string  `json:"language"`
		Price             *float64 `json:"price"`
		Quantity          *int     `json:"quantity"`
		AvailableQuantity *int     `json:"available_quantity"`
		Location          *string  `json:"location"`
		Status            *string  `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request body",
		})
	}

	if req.ISBN != nil && *req.ISBN != "" && *req.ISBN != *book.ISBN {
		exists, err := api.bookRepo.ISBNExists(*req.ISBN)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to check ISBN existence",
			})
		}
		if exists {
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "Book with this ISBN already exists",
			})
		}
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.ISBN != nil {
		book.ISBN = req.ISBN
	}
	if req.Publisher != nil {
		book.Publisher = req.Publisher
	}
	if req.PublicationYear != nil {
		book.PublicationYear = req.PublicationYear
	}
	if req.Genre != nil {
		book.Genre = req.Genre
	}
	if req.Description != nil {
		book.Description = req.Description
	}
	if req.Pages != nil {
		book.Pages = req.Pages
	}
	if req.Language != nil {
		book.Language = *req.Language
	}
	if req.Price != nil {
		book.Price = req.Price
	}
	if req.Quantity != nil {
		book.Quantity = *req.Quantity
	}
	if req.AvailableQuantity != nil {
		book.AvailableQuantity = *req.AvailableQuantity
	}
	if req.Location != nil {
		book.Location = req.Location
	}
	if req.Status != nil {
		book.Status = *req.Status
	}

	if err := api.bookRepo.Update(book); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to update book",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    book,
		Message: "Book updated successfully",
	})
}

func (api *BookAPI) deleteBook(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Book ID is required",
		})
	}

	_, err := api.bookRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Message: "Book not found",
		})
	}

	if err := api.bookRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to delete book",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    map[string]string{"id": id},
		Message: "Book deleted successfully",
	})
}

func (api *BookAPI) updateQuantity(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Book ID is required",
		})
	}

	var req struct {
		Quantity          int `json:"quantity"`
		AvailableQuantity int `json:"available_quantity"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request body",
		})
	}

	if req.Quantity < 0 || req.AvailableQuantity < 0 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Quantities cannot be negative",
		})
	}

	if req.AvailableQuantity > req.Quantity {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Available quantity cannot exceed total quantity",
		})
	}

	_, err := api.bookRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Message: "Book not found",
		})
	}

	if err := api.bookRepo.UpdateQuantity(id, req.Quantity, req.AvailableQuantity); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to update book quantity",
		})
	}

	book, err := api.bookRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to retrieve updated book",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    book,
		Message: "Book quantity updated successfully",
	})
}
