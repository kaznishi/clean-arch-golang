package repository

import "github.com/kaznishi/clean-arch-golang/domain/model"

type AccountRepository interface {
	Get(id int) *model.Account
	GetByAccountNumber(an *model.AccountNumber) *model.Account
}