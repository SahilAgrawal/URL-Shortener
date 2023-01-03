package serviceinterface

import (
	"URL-Shortener/model"
)

type RedirectServiceInterface interface {
	Find(code string) (*model.Redirect, error)
	Store(redirect *model.Redirect) error
	All() ([]model.Redirect, error)
}
