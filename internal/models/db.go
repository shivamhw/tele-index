package models

import (
	"time"
)

type UserModel struct {
	Id          string
	UserName    string
	TeleId      int64
	CreateOn    time.Time
	Password    string
}

type CollectionModel struct {
	Id          string
	Name        string
	Description string
	Owner       string
	CreatedOn   time.Time
}

type Collection_ItemModel struct {
	ColId	string
	ItemId	string
}

type Collection_UserModel struct {
	ColId	string
	UserId	string
	Permission int
}

type ItemModel struct {
	Id        string
	MsgId     int64
	From      int64
	ChatID    int64
	FileName  string
	CreatedOn time.Time
}
