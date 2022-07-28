package repository

import (
	"github.com/tangoctrung/golang_api_v2/entity"
	"gorm.io/gorm"
)

type BookRepository interface {
	InsertBook(book entity.Book) entity.Book
	UpdateBook(book entity.Book) entity.Book
	DeleteBook(book entity.Book)
	FindBookByID(id uint64) entity.Book
	GetAllBooks() []entity.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookConnection{
		connection: db,
	}
}

func (db *bookConnection) InsertBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) UpdateBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) DeleteBook(book entity.Book) {
	db.connection.Delete(&book, book.ID)
}

func (db *bookConnection) FindBookByID(bookID uint64) entity.Book {
	var book entity.Book
	db.connection.Preload("User").Find(&book, bookID)
	return book
}

func (db *bookConnection) GetAllBooks() []entity.Book {
	var books []entity.Book
	db.connection.Preload("User").Find(&books)
	return books
}
