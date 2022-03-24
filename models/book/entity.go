package book

import (
	"fmt"

	"package.local/hktn/helper"
	"package.local/hktn/models/author"
)

const (
	max_page     int     = 1000
	min_page     int     = 25
	max_year     int     = 2022
	min_year     int     = 1970
	max_price    float64 = 10.05
	min_price    float64 = 155.83
	max_quantity int     = 20
	min_quantity int     = 0
)

type BookList struct {
	BookList []Book
}

type Book struct {
	ID                            uint `gorm:"primaryKey"`
	NumberOfPages, Year, Quantity int
	Price                         float64
	Name                          string
	AuthorId                      int
	Sku, Isbn                     string
	IsDeleted                     bool
	Author                        author.Author `gorm:"foreignKey:AuthorId;references:AuthorId"`
}

//InitializeBookList() kitaplar BookList atılıyor
func InitializeBookList(bookCsv []helper.BookCsv) BookList {

	var bookList BookList
	var books []Book

	for _, book := range bookCsv {
		var b Book = fillBookValues(book.Name, book.AuthorId)

		books = append(books, b)
	}
	bookList.BookList = books

	return bookList
}

//fillBookValues() kitap alanları türetiliyor
func fillBookValues(name string, authorId int) Book {

	var bf Book
	bf.Name = name
	bf.AuthorId = authorId
	bf.NumberOfPages = helper.RandomIntegerCreator(min_page, max_page)
	bf.Year = helper.RandomIntegerCreator(min_year, max_year)
	bf.Price = helper.RoundFloat(helper.RandomFloatCreator(min_price, max_price))
	bf.Quantity = helper.RandomIntegerCreator(min_quantity, max_quantity)
	bf.Sku = helper.CreateSku(name)
	bf.Isbn = helper.CreateIsbn()
	bf.IsDeleted = false

	return bf
}

//ToListString() anlaşılır şekilde gerekli alanlar geri dönüyor
func (b *Book) ToListString() string {
	return fmt.Sprintf("Id : %d, Name : %s, Author : %s, Quantity : %d, Price : %g, Sku : %s, ISBN : %s, isDeleted : %t", b.ID, b.Name, b.Author.Name, b.Quantity, b.Price, b.Sku, b.Isbn, b.IsDeleted)
}
