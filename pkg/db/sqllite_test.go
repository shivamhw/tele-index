package db

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/shivamhw/tele-index/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	dbPath = "./testdata/data.db"
	user = models.UserModel{}
	col = models.CollectionModel{}
	item = models.ItemModel{}
)

func Init() *SqlDB {
	opts := &SqliteOpts{
		Path: dbPath,
	}
	db, err := NewSqlDb(opts)
	if err != nil {
		panic(err)
	}
	user = models.UserModel{
		Id:       uuid.NewString(),
		UserName: "shivamhw" + strconv.Itoa(rand.Int()),
		TeleId:   int64(rand.Int()),
		Password: "test12312",
	}
	_, err = db.CreateUser(user)
	if err != nil {
		panic(err)
	}
	col = models.CollectionModel{
		Id:          uuid.NewString(),
		Name:        "demo collection",
		Description: "thi sis demo",
		Owner:       user.Id,
	}
	_, err = db.CreateCollection(col)
	if err != nil {
		panic(err)
	}
	item = models.ItemModel{
		Id:        uuid.NewString(),
		MsgId:     123123,
		From:      231,
		ChatID:    123,
		FileName:  "testfile.jpg",
	}
	_, err = db.CreateItem(item)
	if err!= nil{
		panic(err)
	}
	return db
}

func TestCreateUser(t *testing.T) {
	db := Init()
	u1 := models.UserModel{
		Id:       uuid.NewString(),
		UserName: "shivamhw1" + strconv.Itoa(rand.Int()),
		TeleId:   rand.Int63(),
		Password: "test12312",
	}
	u, err := db.CreateUser(u1)
	if err != nil {
		t.Fatal(err)
	}
	u1.TeleId = 1
	_, err = db.CreateUser(u1)
	if err == nil || !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatal("this should not happening")
	}
	u1.UserName = "test324"
	u1.TeleId = user.TeleId
	_, err = db.CreateUser(u1)
	if err == nil || !errors.Is(err, ErrTeleIdAlreadyExists) {
		t.Fatal("this should not happening, telegram id")
	}
	fmt.Println(u)
}

func TestGetUserByUsername(t *testing.T) {
	db := Init()
	username := user.UserName
	u, err := db.GetUserByUsername(username)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
	_, err = db.GetUserByUsername(username + "23")
	if err == nil || !errors.Is(err, ErrUserNotFound) {
		t.Fatal(err)
	}
}

func TestCreateCollection(t *testing.T) {
	db := Init()
	c1 := models.CollectionModel{
		Id:          uuid.NewString(),
		Name:        "demo collection",
		Description: "thi sis demo",
		Owner:       user.Id,
	}
	c, err := db.CreateCollection(c1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
	c1.Id = "test324"
	c1.Owner = "serfdsf"
	c, err = db.CreateCollection(c1)
	if err == nil && !errors.Is(err, ErrInvalidUserId) {
		t.Fatal("this should failed")
	}
	fmt.Println(c)
}

func TestGetCollection(t *testing.T) {
	db := Init()
	c, err := db.GetCollectionById(col.Id)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
	_, err = db.GetCollectionById("test324")
	if err == nil || !errors.Is(err, ErrCollectionNotFound) {
		t.Fatal(err)
	}
}

func TestGetCollectionsForUser(t *testing.T) {
	db := Init()
	c, err := db.GetCollectionsForUser(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c))
	fmt.Println(c)
	c, err = db.GetCollectionsForUser("reand")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c))
}

func TestCreateItem(t *testing.T) {
	db := Init()
	it := models.ItemModel{
		Id:        uuid.NewString(),
		MsgId:     rand.Int63(),
		From:      rand.Int63(),
		ChatID:    rand.Int63(),
		FileName:  "testfile.jpg",
	}
	i, err := db.CreateItem(it)
	assert.NoError(t, err)
	fmt.Println(i)
}

func TestGetItem(t *testing.T) {
	db := Init()
	i, err := db.GetItem(item.Id)
	assert.NoError(t, err)
	fmt.Println(i)
	i, err = db.GetItem("resr1234")
	assert.ErrorIs(t, err, ErrItemNotFound)
	assert.Empty(t, i)
}

func TestAddItemToCollection(t *testing.T) {
	db := Init()
	it := models.ItemModel{
		Id:        uuid.NewString(),
		MsgId:     rand.Int63(),
		From:      rand.Int63(),
		ChatID:    rand.Int63(),
		FileName:  "testfile.jpg",
	}
	_, err := db.CreateItem(it)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddItemToCollection(it.Id, col.Id)
	assert.NoError(t, err)
	err = db.AddItemToCollection(it.Id, col.Id)
	assert.Error(t, err)
}

func TestGetItemsInCollection(t *testing.T) {
	db := Init()
	it := models.ItemModel{
		Id:        uuid.NewString(),
		MsgId:     rand.Int63(),
		From:      rand.Int63(),
		ChatID:    rand.Int63(),
		FileName:  "testfile.jpg",
	}
	_, err := db.CreateItem(it)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddItemToCollection(it.Id, col.Id)
	if err != nil {
		t.Fatal(err)
	}
	i, err := db.GetItemsInCollection(col.Id)
	assert.NoError(t, err)
	assert.Contains(t, i, it.Id)
}

func TestAddUserToCollection(t *testing.T) {
	db := Init()
	u := models.UserModel{
		Id:       uuid.NewString(),
		UserName: "shivamhw" + strconv.Itoa(rand.Int()),
		TeleId:   int64(rand.Int()),
		Password: "test12312",
	}
	_, err := db.CreateUser(u)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddUserToCollection(u.Id, col.Id)
	assert.NoError(t, err)
	err = db.AddItemToCollection(u.Id, col.Id)
	assert.Error(t, err)
}