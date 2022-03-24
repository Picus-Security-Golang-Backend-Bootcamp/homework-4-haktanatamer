package main

import (
	"errors"
	"fmt"
	"os"

	"package.local/hktn/helper"
	"package.local/hktn/infrastructure"
	"package.local/hktn/models/author"
	"package.local/hktn/models/book"
	reqlog "package.local/hktn/models/reqLog"
)

var (
	bookList   book.BookList
	authorList author.AuthorList
	reset      string
	searchWord string
	counter    int

	bookRepository   *book.BookRepository
	authorRepository *author.AuthorRepository
	reqRepository    *reqlog.ReqRepository
)

var (
	ErrNotEnoughStock     = errors.New("Yeterli stok bulunmamaktadır")
	ErrZeroValue          = errors.New("0 dan büyük bir değer giriniz")
	ErrBookAlreadyDeleted = errors.New("Kitap silinemez")
	ErrBookNotFound       = errors.New("Kitap bulunamadı")
	ErrAuthorNotFound     = errors.New("Yazar bulunamadı")
)

const (
	searchChoice   int    = 1
	listAllChoice  int    = 2
	buyChoice      int    = 3
	deleteChoice   int    = 4
	exitChoice     int    = 5
	addBook        int    = 6
	addAuthor      int    = 7
	resetChoice    string = "r"
	wrongChoice    string = "Yanlış seçim. Tekrar deneyiniz"
	notFoundChoice string = "Kitap bulunamadı. Tekrar deneyiniz"
	isAdd          bool   = true //true yaparsak kitap ve yazar ekleme özelliği açılıyor
)

//db bilgileri
const (
	username = "root"
	password = "1575"
	hostname = "127.0.0.1:3306"
	nameDb   = "golang"
)

//getConnectionString() connectionString oluşturuluyor
func getConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, nameDb)
}

func init() {

	dbTransactions()

}

//dbTransactions() db bağlantı açılıyor migration işlemi varsa yapılıyor
func dbTransactions() {

	db := infrastructure.NewMySQLDB(getConnectionString())

	bookRepository = book.NewBookRepository(db)
	authorRepository = author.NewAuthorRepository(db)
	reqRepository = reqlog.NewReqRepository(db)

}

func main() {

	printUsage()
}

//GetRepos() repolari geri doner
func GetRepos() (*book.BookRepository, *author.AuthorRepository, *reqlog.ReqRepository) {

	return bookRepository, authorRepository, reqRepository
}

//printUsage() uygulama işlemlerinin console basılıyor
func printUsage() {

	fmt.Println("Kitaplık uygulamasında kullanabileceğiniz komutlar :")
	fmt.Printf(" search => arama işlemi için %d\n", searchChoice)
	fmt.Printf(" list => listeleme işlemi için %d\n", listAllChoice)
	fmt.Printf(" buy => kitap satın almak için %d\n", buyChoice)
	fmt.Printf(" delete => kitap silmek için %d\n", deleteChoice)
	fmt.Printf(" exit => uygulamadan çıkmak için %d\n", exitChoice)
	if isAdd {
		fmt.Printf(" add new book => yeni kitap eklemek için %d\n", addBook)
		fmt.Printf(" add new author => yeni yazar eklemek için %d\n", addAuthor)
	}
	fmt.Println("Tuşlarına basınız")

	getChoiceFromUser()
}

//getChoiceFromUser() kullanıcının yapmak istediği işlemi seçiyor
func getChoiceFromUser() {
	var choice int

	fmt.Scan(&choice)

	switch choice {
	case searchChoice:
		search()
	case listAllChoice:
		listAll()
	case buyChoice:
		buyBook()
	case deleteChoice:
		deleteBook()
	case exitChoice:
		terminate()
	case addBook:
		addNewBook()
	case addAuthor:
		addNewAuthor()
	default:
		helper.AddSeparator()
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		printUsage()
	}
}

//addNewAuthor() yeni kitap ekleme işlemleri
func addNewAuthor() {
	fmt.Println("Yazar Adını Giriniz : ")

	var newAuthorName string

	if _, err := fmt.Scan(&newAuthorName); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		addNewAuthor()
	}

	authorRepository.AddAuthor(newAuthorName)

	var count int = printAllAuthors()
	fmt.Printf("%d Tane Yazar Getirildi..\n", count)
	helper.AddSeparator()
	fmt.Println("Bu işlemi postman üzerinden aşağıdaki istekle yapabilirsiniz")
	fmt.Println("http://localhost:8000/authors/create")
	fmt.Println("isteğin body kısmı => {\"name\":\"yazar_Adi_giriniz\"}")
	helper.AddSeparator()
	conclude()
}

//printAllAuthors tüm yazarları konsola yazdır
func printAllAuthors() int {
	var authors []author.Author = authorRepository.FindAll()

	for _, a := range authors {
		fmt.Println(a.ToListString())
	}

	helper.AddSeparator()
	return len(authors)
}

//addNewBook() yeni kitap ekleme işlemleri
func addNewBook() {
	fmt.Println("Kitap Adını Giriniz : ")

	var newBookName string

	if _, err := fmt.Scan(&newBookName); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		addNewBook()
	}

	helper.AddSeparator()
	printAllAuthors()

	addBookAuthor(newBookName)
}

//addBookAuthor() kitaba yazar ekleniyor
func addBookAuthor(newBookName string) {
	fmt.Println("Kitap Yazarının Id'sini Girin : ")

	var authorId int

	if _, err := fmt.Scan(&authorId); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		addBookAuthor(newBookName)
	}

	selectedAuthor, _, err := checkAuthorIdExist(authorId)

	if err != nil {
		fmt.Println(err)
		addBookAuthor(newBookName)
	}

	var bookCsv helper.BookCsv
	bookCsv.AuthorId = selectedAuthor.AuthorId
	bookCsv.Name = newBookName

	var books []helper.BookCsv
	books = append(books, bookCsv)
	bookList = book.InitializeBookList(books)
	bookRepository.InsertInitialData(bookList.BookList)
	fmt.Println("Kitap eklendi")
	helper.AddSeparator()
	fmt.Println("Bu işlemi postman üzerinden aşağıdaki istekle yapabilirsiniz")
	fmt.Println("http://localhost:8000/books/create")
	fmt.Println("isteğin body kısmı => {\"name\":\"kitap_Adi_giriniz\",\"authorId\":yazar_id}")
	fmt.Println("eğer yazar id bilmiyorsanız yazar bilgileri için => http://localhost:8000/authors")
	helper.AddSeparator()
	conclude()
}

// checkAuthorIdExist() girilen id de yazar kontrolü ve varsa atanması
func checkAuthorIdExist(id int) (author.Author, bool, error) {

	var selectedAuthor author.Author = authorRepository.GetByID(id)
	if selectedAuthor.AuthorId == 0 {

		return selectedAuthor, false, ErrAuthorNotFound
	}
	return selectedAuthor, true, nil
}

//deleteBook() kitap silme işlemleri başladı gerekli kontroller ve silme işlemine yönlendirme
func deleteBook() {

	fmt.Println("Silmek İstediğiniz Kitap Id Giriniz : ")

	var selectedBookId int

	if _, err := fmt.Scan(&selectedBookId); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		deleteBook()
	}

	selectedBook, _, err := checkIdExist(selectedBookId)

	if err != nil {
		fmt.Println(err)
		deleteBook()
	}

	bookRepository.UpdateIsDeletedById(int(selectedBook.ID))

	fmt.Println("Kitap silinmiştir")
	helper.AddSeparator()
	fmt.Println("Bu işlemi postman üzerinden aşağıdaki istekle yapabilirsiniz")
	fmt.Println("http://localhost:8000/books/remove/{kitap_id}")
	helper.AddSeparator()
	conclude()
}

//buyBook() kitap satın alma işlemleri başladı gerekli kontroller
func buyBook() {

	fmt.Println("Satın Alınacak Kitap Id Giriniz : ")

	var selectedBookId int

	if _, err := fmt.Scan(&selectedBookId); err != nil {

		fmt.Println(wrongChoice)
		helper.AddSeparator()
		buyBook()
	}

	selectedBook, _, err := checkIdExist(selectedBookId)

	if err != nil {
		fmt.Println(err)
		buyBook()
	}

	purchase(&selectedBook)
}

//purchase() kitap satın alma işlemleri devam ediyor ve satın alma işlemine yönlendirme
func purchase(selectedBook *book.Book) {

	fmt.Println("Satın Alınacak Kitap Sayısını Giriniz : ")

	var numberOfPurchases int

	if _, err := fmt.Scan(&numberOfPurchases); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		purchase(selectedBook)
	}

	if selectedBook.Quantity < numberOfPurchases {
		fmt.Println(ErrNotEnoughStock)
		helper.AddSeparator()
		purchase(selectedBook)
	}

	bookRepository.UpdateQuantityById(int(selectedBook.ID), (selectedBook.Quantity - numberOfPurchases))

	newSelectedBook, _, _ := checkIdExist(int(selectedBook.ID))

	fmt.Println("Satın alma işlemi sonrası kitap durumu")
	fmt.Printf("%+v\n", newSelectedBook)
	helper.AddSeparator()
	fmt.Println("Bu işlemi postman üzerinden aşağıdaki istekle yapabilirsiniz")
	fmt.Println("http://localhost:8000/books/buy/{kitap_id}/{satin_alinacak_kitap_sayisi}")
	helper.AddSeparator()
	conclude()
}

//checkIdExist() girilen id de kitap kontrolü ve varsa atanması
func checkIdExist(bookId int) (book.Book, bool, error) {

	var selectedBook book.Book = bookRepository.GetBookByID(bookId)
	if selectedBook.ID == 0 {

		return selectedBook, false, ErrBookNotFound
	}
	if selectedBook.IsDeleted == true {

		return selectedBook, false, ErrBookAlreadyDeleted
	}
	return selectedBook, true, nil
}

//search() kitap arama işlemleri başladı gerekli kontroller ve aramaya yönlendirme
func search() {

	fmt.Println("Aranacak Kelimeyi Giriniz : ")
	fmt.Scan(&searchWord)
	helper.AddSeparator()
	fmt.Printf("%s Kelimesi Aranıyor..\n", searchWord)
	helper.AddSeparator()
	searchByWord(searchWord)
	fmt.Println("Arama Tamamlandı..")
	helper.AddSeparator()
	conclude()
}

//searchByWord() kitap arama işlemi, ekrana basılması
func searchByWord(word string) {

	fmt.Println("Kitaplar Getiriliyor..")
	helper.AddSeparator()
	var dbBooks []book.Book = bookRepository.BookSearch(word)

	for _, dbBook := range dbBooks {
		fmt.Println(dbBook.ToListString())
	}
	helper.AddSeparator()
	fmt.Printf("%d Tane Kitaplar Getirildi..\n", len(dbBooks))
	helper.AddSeparator()
	fmt.Println("Bu işlemi postman üzerinden aşağıdaki istekle yapabilirsiniz")
	fmt.Println("http://localhost:8000/books/search/{aranacak_kelimeyi_buraya_yazın}")
	helper.AddSeparator()
}

//listAll() tüm kitapların listelenmesi
func listAll() {

	fmt.Println("Kitaplar Getiriliyor..")
	helper.AddSeparator()
	var dbBooks []book.Book = bookRepository.FindAll()

	for _, dbBook := range dbBooks {
		fmt.Println(dbBook.ToListString())
	}
	helper.AddSeparator()
	fmt.Printf("%d Tane Kitaplar Getirildi..\n", len(dbBooks))
	helper.AddSeparator()
	fmt.Println("Bu işlemi postman üzerinden de aşağıdaki istekle yapabilirsiniz")
	fmt.Println("http://localhost:8000/books")
	helper.AddSeparator()
	conclude()
}

//terminate() uygulamayı sonlandırma
func terminate() {

	fmt.Println("Uygulama Sonlandırılıyor..")
	os.Exit(3)
}

//conclude() uygulamaya devam edilecek mi kontrolleri
func conclude() {

	fmt.Println("Yeni Bir Arama Yapmak İçin R \nUygulamanyı Sonlandırmak İçin Herhangi Bir Tuşa Basınız")

	fmt.Scan(&reset)

	switch helper.StringLower(reset) {
	case resetChoice:
		main()
	default:
		terminate()
	}
}
