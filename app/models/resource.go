package models

import (
	"time"
	"github.com/go-gorp/gorp"
)

type Resource struct {
	Id       		int64  		`db:"id, primarykey, autoincrement"`
	CreatedAt  		time.Time 	`db:"created_at"`
	UpdatedAt  		time.Time 	`db:"updated_at"`
	LastRetrieval	time.Time 	`db:"last_retrieval"`
	LongURL  		string 		`db:"long_url, size:512"`
}

func (i *Resource) PreInsert(s gorp.SqlExecutor) error {
	i.CreatedAt = time.Now()
	i.UpdatedAt = i.CreatedAt
	return nil
}

func (i *Resource) PreUpdate(s gorp.SqlExecutor) error {
	i.UpdatedAt = time.Now()
	return nil
}
