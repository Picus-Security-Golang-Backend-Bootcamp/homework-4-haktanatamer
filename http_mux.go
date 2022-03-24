package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"package.local/hktn/helper"
	"package.local/hktn/models/author"
	"package.local/hktn/models/book"
	reqlog "package.local/hktn/models/reqLog"
)

func init() {
	CreateServer()
}

//CreateServer() server olusturuluyor servis yollari tanımlanıyor log ve auth işlemleri
func CreateServer() {
	r := mux.NewRouter()

	CORSOptions()

	r.Use(loggingMiddleware)

	r.Use(authenticationMiddleware)

	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/search/{word}", getBookFiltered).Methods("GET")
	r.HandleFunc("/books/buy/{bookId:[0-9]+}/{count:[0-9]+}", updateBook).Methods("PUT")
	r.HandleFunc("/books/remove/{bookId:[0-9]+}", removeBook).Methods("PUT")
	r.HandleFunc("/authors", getAuthors).Methods("GET")
	r.HandleFunc("/authors/create", createAuthor).Methods("POST")
	r.HandleFunc("/books/create", createBook).Methods("POST")

	srv := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		helper.AddSeparator()
		fmt.Println("Http Sunucu Oluşturuldu")
		fmt.Println("http://localhost:8000/ ile ulaşabilirsiniz")
		fmt.Println("Tüm istekler için token zorunludur.")
		helper.AddSeparator()
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(srv, time.Second*10)

}

//authenticationMiddleware() auth kontrolleri
func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/") {
			if token != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token bulunamadi", http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

//ShutdownServer() servis sonlandırma
func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

//getBookFiltered() servis kitap sorgulama
func getBookFiltered(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	bR, _, _ := GetRepos()
	var dbBooks []book.Book = bR.BookSearch(params["word"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbBooks)

}

//getBooks() servis tüm kitaplari getir
func getBooks(w http.ResponseWriter, r *http.Request) {

	bR, _, _ := GetRepos()
	var dbBooks []book.Book = bR.FindAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbBooks)

}

//updateBook() servis kitap güncelle
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	bR, _, _ := GetRepos()
	var bId, _ = strconv.Atoi(params["bookId"])
	var cId, _ = strconv.Atoi(params["count"])

	var selectedBook book.Book = getBookById(*bR, bId)

	var d ApiResponse

	if selectedBook.ID == 0 || selectedBook.IsDeleted {
		d = ApiResponse{
			Data: ErrBookNotFound.Error(),
		}

	} else if cId < 0 {
		d = ApiResponse{
			Data: ErrZeroValue.Error(),
		}
	} else if cId > selectedBook.Quantity {

		d = ApiResponse{
			Data: ErrNotEnoughStock.Error() + " Satın alınabilecek maksimum kitap : " + strconv.Itoa(selectedBook.Quantity),
		}

	} else {

		bR.UpdateQuantityById(int(selectedBook.ID), (selectedBook.Quantity - cId))
		d = ApiResponse{
			Data: "Kitap Satışı Tamamlandı",
		}
	}
	resp, _ := json.Marshal(d)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

//getBookById id göre kitap getir
func getBookById(repo book.BookRepository, id int) book.Book {
	return repo.GetBookByID(id)
}

//removeBook() servis kitap sil
func removeBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	bR, _, _ := GetRepos()
	var bId, _ = strconv.Atoi(params["bookId"])

	var selectedBook book.Book = getBookById(*bR, bId)

	var d ApiResponse

	if selectedBook.ID == 0 || selectedBook.IsDeleted {

		d = ApiResponse{
			Data: ErrBookNotFound.Error(),
		}

	} else {

		bR.UpdateIsDeletedById(int(selectedBook.ID))
		d = ApiResponse{
			Data: "Kitap Silindi",
		}
	}
	resp, _ := json.Marshal(d)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

//CORSOptions() işlemleri
func CORSOptions() {
	handlers.AllowedOrigins([]string{"*"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT"})
}

//loggingMiddleware() servis loglama db isteklerin basıldığı yer
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		_, _, lR := GetRepos()

		rl := reqlog.Requests{
			Req:   r.Method + " " + r.URL.String(),
			Token: r.Header["Authorization"][0],
		}

		//r.Method + " " + r.URL.String()

		lR.ReqCreate(rl)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

type ApiResponse struct {
	Data interface{} `json:"result"`
}

type ApiResponseWithAllData struct {
	Result interface{} `json:"result"`
	Data   interface{} `json:"veri"`
}

//getAuthors() servis tüm yazarları getir
func getAuthors(w http.ResponseWriter, r *http.Request) {

	_, aR, _ := GetRepos()
	var dbAuthors []author.Author = aR.FindAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbAuthors)

}

//createAuthor() servis yazar oluştur
func createAuthor(w http.ResponseWriter, r *http.Request) {

	var a author.Author

	err := helper.DecodeJSONBody(w, r, &a)
	if err != nil {
		var mr *helper.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	_, aR, _ := GetRepos()

	aR.AddAuthor(a.Name)

	var d = ApiResponseWithAllData{
		Result: "Yazar Eklendi. Tüm Yazar Listesi",
		Data:   aR.FindAll(),
	}

	resp, _ := json.Marshal(d)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// createBook() servis kitap oluştur
func createBook(w http.ResponseWriter, r *http.Request) {

	var bc helper.BookCsv
	var abb ApiBookBody

	err := helper.DecodeJSONBody(w, r, &abb)
	if err != nil {
		var mr *helper.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	bR, aR, _ := GetRepos()

	a := getAuthorById(*aR, (abb.AuthorId))
	var d ApiResponseWithAllData
	if a.AuthorId == 0 {
		d = ApiResponseWithAllData{
			Result: "Yazar Bulunamadı. Yeni yazar eklemek için (http://localhost:8000/authors/create)  Tüm Yazar Listesi",
			Data:   aR.FindAll(),
		}
	} else {
		bc.Name = abb.Name
		bc.AuthorId = abb.AuthorId
		var books []helper.BookCsv
		books = append(books, bc)
		bookList = book.InitializeBookList(books)
		bR.InsertInitialData(bookList.BookList)
		d = ApiResponseWithAllData{
			Result: "Kitap Eklendi",
		}
	}

	resp, _ := json.Marshal(d)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

type ApiBookBody struct {
	Name     string `json:"name"`
	AuthorId int    `json:"authorId"`
}

//getAuthorById id göre yazar getir
func getAuthorById(repo author.AuthorRepository, id int) author.Author {
	return repo.GetByID(id)
}
