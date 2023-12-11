package library

import (
	"testing"
)

var l *Library

// var want []Book

func init() {
	b1 := Book{Name: "book1", Author: "author1"}
	b2 := Book{Name: "book2", Author: "author2"}
	b3 := Book{Name: "book3", Author: "author3"}

	l = NewLibrary(&SliceStorage{})

	l.AddBook(&b1)
	l.AddBook(&b2)
	l.AddBook(&b3)
}

func TestAdd(t *testing.T) {
	want := []Book{{Name: "book1", Author: "author1"}, {Name: "book2", Author: "author2"}, {Name: "book3", Author: "author3"}}
	got := l.getBooks()

	if len(got) != len(want) {
		t.Errorf("count: got %d, wanted %d", len(got), len(want))
	}

	for i, w := range want {
		if got[i].Name != w.Name || got[i].Author != w.Author {
			t.Errorf("got %q, wanted %q", got[i].Name, w.Name)
		}
	}
}

func TestFind(t *testing.T) {
	if b := l.FindBook("book1"); b == nil {
		t.Errorf("not find book")
	}
}

// func TestFind2(t *testing.T) {
// 	t.Run("Find", func(t *testing.T){
// 		TestFind(asd)
// 	})
// }
