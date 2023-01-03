package repository

import (
	"URL-Shortener/model"
)

type RedirectRepo interface {
	Find(code string) (*model.Redirect, error)
	Store(redirect *model.Redirect) error
	All() ([]model.Redirect, error)
}
