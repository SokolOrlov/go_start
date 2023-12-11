package library

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {

	b1 := Book{Name: "book1", Author: "author1"}
	b2 := Book{Name: "book2", Author: "author2"}
	b3 := Book{Name: "book3", Author: "author3"}

	l := NewLibrary(&MapStorage{})

	l.AddBook(b1)
	l.AddBook(b2)
	l.AddBook(b3)

	// want := []Book{b1, b2, b3}
	want := []Book{{Name: "book1", Author: "author1"}, {Name: "book2", Author: "author2"}, {Name: "book3", Author: "author3"}}
	got := l.getBooks()

	if len(got) != len(want) {
		t.Errorf("count: got %d, wanted %d", len(got), len(want))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
