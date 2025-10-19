package db

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shivamhw/tele-index/internal/models"
	"github.com/shivamhw/tele-index/pkg/log"
)

type SqliteOpts struct {
	Path string
}

type SqlDB struct {
	Db *sql.DB
}

func NewSqlDb(opts *SqliteOpts) (*SqlDB, error) {
	db, err := sql.Open("sqlite3", opts.Path+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	return &SqlDB{
		Db: db,
	}, nil
}

func (s *SqlDB) CreateUser(u models.UserModel) (string, error) {
	res, err := s.Db.ExecContext(context.Background(), "INSERT INTO users (id, username, teleId, password) VALUES (?, ?, ?, ?)", u.Id, u.UserName, u.TeleId, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return "", ErrUserAlreadyExists
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.teleId") {
			return "", ErrTeleIdAlreadyExists
		}
		return "", err
	}
	log.InfoF("success created user", "res", res)
	return u.Id, nil
}


func (s *SqlDB) GetUserById(id string) (u models.UserModel, err error) {
	err = s.Db.QueryRowContext(context.Background(), "SELECT * from users where id = ?", id).Scan(&u.Id, &u.UserName, &u.TeleId, &u.CreateOn, &u.Password)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result ") {
		return u, ErrUserNotFound
	}
	return
}

func (s *SqlDB) GetUserByUsername(username string) (u models.UserModel, err error) {
	err = s.Db.QueryRowContext(context.Background(), "SELECT * from users where username = ?", username).Scan(&u.Id, &u.UserName, &u.TeleId, &u.CreateOn, &u.Password)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result ") {
		return u, ErrUserNotFound
	}
	return
}

func (s *SqlDB) CreateCollection(c models.CollectionModel) (id string, err error) {
	res, err := s.Db.ExecContext(context.Background(), "INSERT INTO collections (id, name, description, owner) VALUES (?,?,?,?)", c.Id, c.Name, c.Description, c.Owner)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return id, ErrInvalidUserId
		}
		return
	}
	log.InfoF("successfull created collection", "res", res)
	id = c.Id
	return
}

func (s *SqlDB) GetCollectionById(id string) (c models.CollectionModel, err error) {
	err = s.Db.QueryRowContext(context.Background(), "SELECT * from collections where id = ?", id).Scan(&c.Id, &c.Name, &c.Description, &c.Owner, &c.CreatedOn)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result") {
		return c, ErrCollectionNotFound
	}
	return
}


func (s *SqlDB) GetCollectionsForUser(owner_id string) (cols []models.CollectionModel, err error) {
	res, err := s.Db.QueryContext(context.Background(), "SELECT *  FROM collections where owner = ?", owner_id)
	if err != nil {
		return
	}
	for res.Next() {
		t := models.CollectionModel{}
		res.Scan(&t.Id, &t.Name, &t.Description, &t.Owner, &t.CreatedOn)
		cols = append(cols, t)
	}
	return
}

func (s *SqlDB) CreateItem(it models.ItemModel) (id string, err error) {
	res, err := s.Db.ExecContext(context.Background(), "INSERT INTO items (id, msg_id, from_id, chat_id, filename) VALUES (?,?,?,?,?)", it.Id, it.MsgId, it.From, it.ChatID, it.FileName)
	if err != nil {
		return
	}
	log.InfoF("successfully created item", "res", res)
	id = it.Id
	return
}

func (s *SqlDB) GetItem(id string) (it models.ItemModel, err error) {
	err = s.Db.QueryRowContext(context.Background(), "SELECT * FROM items WHERE id = ?", id).Scan(&it.Id, &it.MsgId, &it.From, &it.ChatID, &it.FileName, &it.CreatedOn)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result") {
		return it, ErrItemNotFound
	}
	return
}

// trust user is sending correct ids, validate at service level
func (s *SqlDB) AddItemToCollection(itemId string, colId string) error {
	res, err := s.Db.ExecContext(context.Background(), "INSERT INTO col_items (item_id, col_id) VALUES (?,?)", itemId, colId)
	if err != nil {
		return err
	}
	log.InfoF("successfully added item to col", "item", itemId, "col", colId, "res", res)
	return nil
}

// check col id at service level
func (s *SqlDB) GetItemsInCollection(colId string) (it []string, err error) {
	res, err := s.Db.QueryContext(context.Background(), "SELECT item_id FROM col_items WHERE col_id = ?", colId)
	if err != nil {
		return
	}
	for res.Next() {
		var t string
		res.Scan(&t)
		it = append(it, t)
	}
	return
}

func (s *SqlDB) AddUserToCollection(userId string, colId string) error {
	res, err := s.Db.ExecContext(context.Background(), "INSERT INTO col_users (user_id, col_id) VALUES (?,?)", userId, colId)
	if err != nil {
		return err
	}
	log.InfoF("successfully added user to col", "user", userId, "col", colId, "res", res)
	return nil
}

func (s *SqlDB) AddTelegramToUser(userId string, teleId int64) error {
	res, err := s.Db.ExecContext(context.Background(), "UPDATE users SET teleId = ? WHERE id = ?", teleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.teleId") {
			return ErrTeleIdAlreadyExists
		}
		return err
	}
	log.InfoF("successfully added telegram to user", "user", userId, "teleId", teleId, "res", res)
	return nil
}