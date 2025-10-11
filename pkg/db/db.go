package db

import "github.com/shivamhw/tele-index/internal/models"


type DB interface {
	CreateUser(models.UserModel) (string, error)
	GetUser(id string) (models.UserModel, error)
	CreateCollection(models.CollectionModel) (string, error)
	GetCollection(id string) (models.CollectionModel, error)
	GetCollectionsForUser(id string) ([]models.CollectionModel, error)
	GetItem(id string)(models.ItemModel, error)
	GetItemsInCollection(colId string)([]string, error)
	CreateItem(models.ItemModel) (string, error)
	AddItemToCollection(itemId string, colId string) (error)
	AddUserToCollection(userId string, colId string) (error)
}