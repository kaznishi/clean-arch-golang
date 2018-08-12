package dao

import "github.com/kaznishi/clean-arch-golang/domain/model"

type AccountDAO struct {}

func (accountDAO *AccountDAO) Get(id int) *model.Account {
	hoge := &model.Account{}
	return hoge
}

func (accountDAO *AccountDAO) GetByAccountNumber(accountNumber *model.AccountNumber) *model.Account {
	hoge := &model.Account{}
	return hoge
}