package models

import "github.com/shivamhw/tele-index/pkg/util"

type ItemAddRequest struct {
	Id       int64  `json:"id" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
	ChatId   int64  `json:"chat_id" binding:"required"`
	From     int64  `json:"from" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
}

type ItemSearchRequest struct {
	Type  string `json:"type"`
	Q     string `json:"q" binding:"required"`
	Limit int    `json:"limit"`
	From  int    `json:"from"`
}

func (i ItemAddRequest) GetItem() Item {
	return Item{
		Id:       i.Id,
		FileName: i.FileName,
		ChatId:   i.ChatId,
		From:     i.From,
		Size:     i.Size,
		Tokens:   util.GetTokensFromString(i.FileName),
	}
}

type UserCreateRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateCollectionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddItemToCollectionRequest struct {
	MsgId    int64  `json:"msg_id" binding:"required"`
	ChatId   int64  `json:"chat_id" binding:"required"`
	From     int64  `json:"from" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
}

type UserLoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
