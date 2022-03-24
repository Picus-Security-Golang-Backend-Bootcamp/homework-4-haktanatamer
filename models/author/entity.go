package author

import (
	"fmt"

	"package.local/hktn/helper"
)

type AuthorList struct {
	AuthorList []Author
}

type Author struct {
	AuthorId int `gorm:"primary_key"`
	Name     string
}

//InitializeAuthorList() yazarlar AuthorList atılıyor
func InitializeAuthorList(authorCsv []helper.AuthorCsv) AuthorList {

	var authorList AuthorList
	var authors []Author

	for _, author := range authorCsv {
		var a Author

		a.AuthorId = author.ID
		a.Name = author.Name

		authors = append(authors, a)
	}
	authorList.AuthorList = authors

	return authorList
}

//ToListString() anlaşılır şekilde gerekli alanlar geri dönüyor
func (a *Author) ToListString() string {
	return fmt.Sprintf("Id : %d, Name : %s", a.AuthorId, a.Name)
}
