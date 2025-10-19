package db

import "github.com/shivamhw/tele-index/internal/models"


type DB interface {
	CreateUser(models.UserModel) (string, error)
	GetUserByUsername(id string) (models.UserModel, error)
	GetUserById(id string) (models.UserModel, error)
	CreateCollection(models.CollectionModel) (string, error)
	GetCollectionById(id string) (models.CollectionModel, error)
	GetCollectionsForUser(id string) ([]models.CollectionModel, error)
	GetItem(id string)(models.ItemModel, error)
	GetItemsInCollection(colId string)([]string, error)
	CreateItem(models.ItemModel) (string, error)
	AddItemToCollection(itemId string, colId string) (error)
	AddTelegramToUser(userId string, teleId int64) (error)
	AddUserToCollection(userId string, colId string) (error)
}