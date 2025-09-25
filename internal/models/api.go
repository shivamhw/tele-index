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
	Type string `json:"type"`
	Q    string `json:"q" binding:"required"`
	Size int    `json:"size"`
	From int    `json:"from"`
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
