package models

type Item struct {
	Id       int64  `json:"id" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
	ChatId   int64  `json:"chat_id" binding:"required"`
	From     int64  `json:"from" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
	Tokens   string `json:"tokens"`
}