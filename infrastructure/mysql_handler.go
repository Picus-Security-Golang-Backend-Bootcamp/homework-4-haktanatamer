package infrastructure

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gorm.io/gorm/schema"
)

func NewMySQLDB(conString string) *gorm.DB {

	db, err := gorm.Open(mysql.Open(conString), &gorm.Config{
		PrepareStmt: true, // sonraki sorgular için cache
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true,                              // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database : %s", err.Error()))
	}

	return db
}
