package author

import (
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{
		db: db,
	}
}

//FindAll() tüm yazarları getir
func (r *AuthorRepository) FindAll() []Author {
	var authors []Author
	r.db.Find(&authors)

	return authors
}

//Migration() db tablo oluştur
func (r *AuthorRepository) Migration() {
	r.db.AutoMigrate(&Author{})
}

//InsertInitialData() db olmayan yazarlar db basılıyor
func (r *AuthorRepository) InsertInitialData(authors []Author) {
	for _, a := range authors {
		r.db.Where(Author{Name: a.Name}).Attrs(Author{AuthorId: a.AuthorId, Name: a.Name}).FirstOrCreate(&a)
	}
}

//AddAuthor() yeni yazar ekle dbde olmayan
func (r *AuthorRepository) AddAuthor(name string) {
	var author Author
	r.db.Where(Author{Name: name}).Attrs(Author{Name: name}).FirstOrCreate(&author)
}

//GetByID() id göre yazar getir
func (r *AuthorRepository) GetByID(id int) Author {
	var author Author
	r.db.Where("AuthorId = ?", id).Find(&author)
	return author
}
