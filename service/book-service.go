package service

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/tangoctrung/golang_api_v2/dto"
	"github.com/tangoctrung/golang_api_v2/entity"
	"github.com/tangoctrung/golang_api_v2/repository"
)

type BookService interface {
	InsertBook(b dto.BookCreateDTO) entity.Book
	UpdateBook(b dto.BookUpdateDTO) entity.Book
	DeleteBook(b entity.Book)
	FindBookByID(id uint64) entity.Book
	GetAllBooks() []entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (s *bookService) InsertBook(book dto.BookCreateDTO) entity.Book {
	bookInsert := entity.Book{}
	err := smapping.FillStruct(&bookInsert, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	bookAdd := s.bookRepository.InsertBook(bookInsert)
	return bookAdd
}

func (s *bookService) UpdateBook(book dto.BookUpdateDTO) entity.Book {
	bookUpdate := entity.Book{}
	err := smapping.FillStruct(&bookUpdate, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	bookEdit := s.bookRepository.UpdateBook(bookUpdate)
	return bookEdit
}

func (s *bookService) DeleteBook(book entity.Book) {
	s.bookRepository.DeleteBook(book)
}

func (s *bookService) FindBookByID(id uint64) entity.Book {
	return s.bookRepository.FindBookByID(id)
}

func (s *bookService) GetAllBooks() []entity.Book {
	return s.bookRepository.GetAllBooks()
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
