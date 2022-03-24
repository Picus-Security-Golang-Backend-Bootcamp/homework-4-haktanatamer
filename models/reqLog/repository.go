package reqlog

import "gorm.io/gorm"

type ReqRepository struct {
	db *gorm.DB
}

func NewReqRepository(db *gorm.DB) *ReqRepository {
	return &ReqRepository{
		db: db,
	}
}

//ReqCreate() db log basar
func (r *ReqRepository) ReqCreate(req Requests) {
	r.db.Create(&req)
}
