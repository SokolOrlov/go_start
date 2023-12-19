//go:build

package main

import (
	"fmt"
	"qwe/library"
)

func main() {
	b1 := library.Book{Name: "book1", Author: "author1"}
	b2 := library.Book{Name: "book2", Author: "author2"}
	b3 := library.Book{Name: "book3", Author: "author3"}
	b4 := library.Book{Name: "book4", Author: "author4"}
	b5 := library.Book{Name: "book5", Author: "author5"}

	l := library.NewLibrary(&library.MapStorage{})

	l.AddBook(&b1)
	l.AddBook(&b2)
	l.AddBook(&b3)
	l.AddBook(&b4)
	l.AddBook(&b5)

	fmt.Println(l)

	if f := l.FindBook("book3"); f != nil {
		fmt.Println(f)
	}

	if err := l.AddBook(&b3); err != nil {
		fmt.Println(err)
	}
}
