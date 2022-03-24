package book

import (
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

//FindAll() tüm kitapları getir
func (r *BookRepository) FindAll() []Book {
	var books []Book
	r.db.Preload("Author").Find(&books)
	return books
}

//Migration() db tablo oluştur
func (r *BookRepository) Migration() {
	r.db.AutoMigrate(&Book{})
}

//InsertInitialData() db olmayan kitaplar db basılıyor
func (r *BookRepository) InsertInitialData(books []Book) {
	for _, b := range books {
		r.db.Where(Book{Name: b.Name}).Attrs(Book{AuthorId: b.AuthorId, Name: b.Name}).FirstOrCreate(&b)
	}
}

//FindAllWithRawSQL() tüm kitapları getir sql ile
func (r *BookRepository) FindAllWithRawSQL() []Book {
	var books []Book
	r.db.Raw("SELECT NumberOfPages, Year, Quantity,Price,Name,Sku, Isbn,IsDeleted FROM Book ").Scan(&books)

	return books
}

//BookSearch() kitap adı, kitap sku ve yazar adında gelen parametreyi like ile arar
func (r *BookRepository) BookSearch(word string) []Book {
	var books []Book
	r.db.Joins("Book").Joins("Author").Where("Book.Name LIKE ?", "%"+word+"%").Or("Sku LIKE ?", "%"+word+"%").Or("Author.Name LIKE ?", "%"+word+"%").Find(&books)
	return books
}

//GetBookByID() id göre kitap getir
func (r *BookRepository) GetBookByID(id int) Book {
	var book Book
	r.db.Where("ID = ?", id).Find(&book)
	return book
}

//UpdateQuantityById() id göre stok günceller
func (r *BookRepository) UpdateQuantityById(id, quantity int) {
	r.db.Exec("update book set Quantity = ? where ID = ?", quantity, id)
}

//UpdateIsDeletedById() id göre kitap silinme durumunu günceller
func (r *BookRepository) UpdateIsDeletedById(id int) {
	r.db.Exec("update book set isDeleted = ? where ID = ?", true, id)
}
