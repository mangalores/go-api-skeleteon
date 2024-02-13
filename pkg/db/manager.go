package db

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	Handle(q QueryObject) QueryObject
	Supports(t interface{}) bool
}

type RepositoryNotFoundErr struct {
	t interface{}
}

func (e RepositoryNotFoundErr) Error() string {
	return fmt.Sprintf("could not find repository for type %t", e.t)
}

func NewRepositoryNotFound(t interface{}) RepositoryNotFoundErr {
	return RepositoryNotFoundErr{t}
}

type DefaultRepositoryNotSetErr struct {
}

func (e DefaultRepositoryNotSetErr) Error() string {
	return "default repository was not set"
}

func NewDefaultRepositoryNotSetErr() DefaultRepositoryNotSetErr {
	return DefaultRepositoryNotSetErr{}
}

type QueryManager struct {
	defaultRepository Repository
	repositories      []Repository
	db                *gorm.DB
}

func NewQueryManager(db *gorm.DB) *QueryManager {
	return &QueryManager{
		defaultRepository: NewQueryHandler(db),
		db:                db,
	}
}

func (m *QueryManager) DB() *gorm.DB {
	return m.db
}

func (m *QueryManager) Register(r Repository) {
	m.repositories = append(m.repositories, r)
}

func (m *QueryManager) Get(t interface{}) (Repository, error) {
	for _, r := range m.repositories {
		if r.Supports(t) {
			return r, nil
		}
	}

	if m.defaultRepository != nil {
		return m.defaultRepository, nil
	}

	return nil, NewRepositoryNotFound(t)
}
func (m *QueryManager) Default() (Repository, error) {
	if m.defaultRepository == nil {
		return nil, NewDefaultRepositoryNotSetErr()
	}

	return m.defaultRepository, nil
}
