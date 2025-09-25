package db

import "github.com/shivamhw/tele-index/internal/models"


type DB interface {
	CreateUser(models.UserModel) (models.UserModel, error)
	GetUser(id string) (models.UserModel, error)
	GetCollection(id string) (models.CollectionModel, error)
	GetCollectionsForUser(id string) ([]models.CollectionModel, error)
	GetItem(id string)(models.ItemModel, error)
	CreateItem(models.ItemModel) (models.ItemModel, error)
	AddItemToCollection(itemId string, colId string) (models.CollectionModel, error)
}