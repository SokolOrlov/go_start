package library

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
)

var ErrBookElreadyExsists = errors.New("Book already exsists in library.")

// //////////////////////////////////////////////--Book
type Book struct {
	id     *uint32
	Author string
	Name   string
}

func (b Book) GetName() string {
	return b.Name
}
func (b Book) String() string {

	return fmt.Sprintf("{Id: %d, Name: \"%s\", Author: %s}\n", *b.id, b.Name, b.Author)
}

// //////////////////////////////////////////////--Storage
type Storage interface {
	AddBook(Book) error
	FindBook(name string) *Book
	getBooks() []Book
}

type SliceStorage struct {
	books []Book
}

func (s *SliceStorage) AddBook(b Book) error {
	for _, v := range s.books {
		if strings.EqualFold(b.Name, v.Name) {
			return ErrBookElreadyExsists
		}
	}

	bid := hash(b.Name)
	b.id = &bid

	s.books = append(s.books, b)

	return nil
}

func (s *SliceStorage) FindBook(name string) *Book {
	for _, v := range s.books {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			return &v
		}
	}

	return nil
}

func (s SliceStorage) getBooks() []Book {
	return s.books
}

type MapStorage struct {
	books map[uint32]Book
}

func (s *MapStorage) AddBook(b Book) error {

	if s.books == nil {
		s.books = make(map[uint32]Book)
	}

	for _, v := range s.books {
		if strings.EqualFold(b.Name, v.Name) {
			return ErrBookElreadyExsists
		}
	}

	bid := hash(b.Name)
	b.id = &bid
	s.books[bid] = b

	return nil
}

func (s *MapStorage) FindBook(name string) *Book {
	for _, v := range s.books {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			return &v
		}
	}
	return nil
}

func (s MapStorage) getBooks() []Book {
	sl := make([]Book, 0, len(s.books))
	for _, v := range s.books {
		sl = append(sl, v)
	}
	return sl
}

// //////////////////////////////////////////////--Library
type Library struct {
	Storage
}

func (l Library) String() string {
	bks := l.getBooks()

	b := make([]string, 0, len(bks))

	for _, v := range bks {
		b = append(b, v.String())
	}

	return strings.Join(b, "\n")
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

/////////////////////////////////////////////////////////////

func NewLibrary(s Storage) *Library {

	return &Library{Storage: s}
}
