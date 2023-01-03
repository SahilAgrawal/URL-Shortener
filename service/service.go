package service

import (
	"URL-Shortener/model"
	"URL-Shortener/repository"
	"errors"
	"time"

	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

type RedirectService struct {
	redirectRepo repository.RedirectRepo
}

func NewRedirectService(r repository.RedirectRepo) repository.RedirectRepo {
	return &RedirectService{
		redirectRepo: r,
	}
}
func (r *RedirectService) Find(code string) (*model.Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *RedirectService) Store(redirect *model.Redirect) error {
	err := validate.Validate(redirect)
	if err != nil {
		return errors.New("invalid db.mongo.store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	return r.redirectRepo.Store(redirect)
}
func (r *RedirectService) All() ([]model.Redirect, error) {
	return r.redirectRepo.All()
}
