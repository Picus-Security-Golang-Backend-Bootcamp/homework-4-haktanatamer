package helper

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// küçük harf yapat
func StringLower(data string) string {
	return strings.ToLower(data)
}

// büyük harf yapar
func StringUpper(data string) string {
	return strings.ToUpper(data)
}

// yollanan trim değerini kaldırır
func StringTrim(data, trim string) string {
	return strings.Trim(data, trim)
}

// yollanan word değeri line alanında var mı
func StringContains(line, word string) bool {
	return strings.Contains(StringLower(line), StringLower(word))
}

// konsola kitap basar
func ConsoleBook(counter int, book string) {
	fmt.Printf("%d . %s \n", counter, book)
}

// konsola çizgi basar
func AddSeparator() {
	fmt.Println("--------------------------")
}

// verilen aralıka int üretir
func RandomIntegerCreator(min, max int) int {
	DelaySystem()
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// üretilen değerler farklılaşması için sistem bekletilir
func DelaySystem() {
	time.Sleep(500 * time.Nanosecond)
}

// verilen aralıka noktalı sayı üretir
func RandomFloatCreator(min, max float64) float64 {
	DelaySystem()
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

// Stok kod üretir kitap adı ilk 2 karakter - 6 haneli sayı - kitap adı son 2 karakter olacak sekilde
func CreateSku(name string) string {
	first2 := name[0:2]
	last2 := name[len(name)-2:]
	var b bytes.Buffer
	b.WriteString(StringUpper(StringTrim(first2, " ")))
	b.WriteString("-")
	b.WriteString(strconv.Itoa(RandomIntegerCreator(100000, 999999)))
	b.WriteString("-")
	b.WriteString(StringUpper(StringTrim(last2, " ")))
	return b.String()
}

// ISBN kod üretir ISBN - 13 haneli sayı olacak sekilde
func CreateIsbn() string {
	rand.Seed(time.Now().Unix())
	var b bytes.Buffer
	b.WriteString("ISBN-")
	b.WriteString(strconv.Itoa(RandomIntegerCreator(1000000000000, 9999999999999)))
	return b.String()
}

// rastgele bool değer döner
func RandomBoolCreator() bool {
	value := RandomIntegerCreator(0, 85)
	if value%2 == 0 {
		return true
	}
	return false
}

// noktalı sayıyı 2 basamak yuvarlar
func RoundFloat(f float64) float64 {
	return math.Round(f*100) / 100
}

type BookCsv struct {
	Name     string `json:"NAME"`
	AuthorId int    `json:"AUTHOR_ID"`
}

type AuthorCsv struct {
	Name string `json:"NAME"`
	ID   int    `json:"AUTHOR_ID"`
}
