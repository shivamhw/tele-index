package db

import "database/sql"

type SqliteOpts struct {
	Path string
}

type SqlDB struct {
	Db *sql.DB
}

func NewSqlDb(opts *SqliteOpts) (*SqlDB, error) {
	db, err := sql.Open("sqlite3", opts.Path + "?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	return &SqlDB{
		Db: db,
	}, nil
}

// CreateUser(models.UserModel) (models.UserModel, error)
// GetUser(id string) (models.UserModel, error)
// GetCollection(id string) (models.CollectionModel, error)
// GetCollectionsForUser(id string) ([]models.CollectionModel, error)
// GetItem(id string)(models.ItemModel, error)
// CreateItem(models.ItemModel) (models.ItemModel, error)
// AddItemToCollection(itemId string, colId string) (models.CollectionModel, error)
