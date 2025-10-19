package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shivamhw/tele-index/internal/models"
	"github.com/shivamhw/tele-index/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) searchItem(c *gin.Context) {
	var data models.ItemSearchRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("searching item", "query: ", data)
	if data.Limit == 0 {
		data.Limit = 10
	}
	res, err := a.indx.Search(data.Q, data.Limit, data.From, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	logger.Info("found results for search", "nof", len(res))
	c.JSON(200, gin.H{
		"count": len(res),
		"data":  res,
	})
}

func (a *App) saveItem(c *gin.Context) {
	var data models.ItemAddRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("saving item", "item: ", data)
	if err := a.indx.Index(data.GetItem()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(201, gin.H{"msg": "success"})
}

func (a *App) getCount(c *gin.Context) {
	count, err := a.indx.GetTotalDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(200, gin.H{"count": count})
}

func (a *App) getAll(c *gin.Context) {
	a.indx.ListAll()
}

func (a *App) createUser(c *gin.Context) {
	var data models.UserCreateRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("creating user", "user: ", data)
	pwd, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed hashing password", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	data.Password = string(pwd)
	user := models.UserModel{
		Id:       uuid.New().String(),
		UserName: data.UserName,
		Password: data.Password,
	}
	id, err := a.db.CreateUser(user)
	if err != nil {
		logger.Error("failed creating user", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	logger.Info("created user successfully", "id: ", id)
	token, err := auth.SignPayload(map[string]interface{}{
		"user_id": id,
		"username": data.UserName,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}, "testToken")
	c.SetCookie("Authorization", token, 3600*72, "/", "", false, true)
	c.JSON(201, gin.H{"msg": "success"})
}

func (a *App) getUser(c *gin.Context) {
	username, ok := c.Get("username")
	if username == "" || !ok {
		logger.Error("missing user_id ")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing user_id "})
		return
	}
	logger.Info("getting user", "user_id: ", username)
	user, err := a.db.GetUserByUsername(username.(string))
	if err != nil {
		logger.Error("failed fetching user", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	c.JSON(200,user)
}

func (a *App) getUserCollections(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		logger.Error("missing user_id in context")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user not found"})
	}
	logger.Info("getting user collections", "user_id: ", userID)
	col, err := a.db.GetCollectionsForUser(userID.(string))
	if err != nil {
		logger.Error("failed fetching user collections", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("fetched user collections successfully", "count: ", len(col))
	c.JSON(200, gin.H{
		"size":       len(col),
		"collections": col,
	})
}

func (a *App) createCollection(c *gin.Context) {
	var data models.CreateCollectionRequest
	userID, ok := c.Get("user_id")
	if !ok {
		logger.Error("missing user_id in context")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user not found"})
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("creating user collection", "user_id: ", userID, "collection: ", data)
	col, err := a.db.CreateCollection(models.CollectionModel{
		Id:          uuid.New().String(),
		Name:        data.Name,
		Description: data.Description,
		Owner:       userID.(string),
	})
	if err != nil {
		logger.Error("failed creating collection", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("created collection successfully", "id: ", col)
	c.JSON(201, gin.H{"msg": "success", "collection_id": col})
}

func (a *App) addItemToCollection(c *gin.Context) {
	var data models.AddItemToCollectionRequest
	userID, ok := c.Get("user_id")
	if !ok {
		logger.Error("missing user_id in context")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user not found"})
	}
	collectionID := c.Param("id")
	if collectionID == "" {
		logger.Error("missing collection_id param")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing collection_id param"})
		return
	}
	logger.Info("checking collection exists","colid", collectionID)
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("created item")
	it := models.ItemModel{
		Id:       uuid.New().String(),
		FileName: data.FileName,
		ChatID:   data.ChatId,
		From:     data.From,
		Size:     data.Size,
		MsgId:    data.MsgId,
	}
	logger.Info("creating item", "item", it)
	id, err := a.db.CreateItem(it)
	if err != nil {
		logger.Error("failed creating item", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("adding item to collection", "item", id, userID)
	err = a.db.AddItemToCollection(it.Id, collectionID)
	if err != nil {
		logger.Error("failed adding item to collection", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("added item to collection successfully")
	c.JSON(201, gin.H{"msg": "success"})
}

func (a *App) getCollection(c *gin.Context) {
	collectionID := c.Param("id")
	if collectionID == "" {
		logger.Error("missing collection_id param")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing collection_id param"})
		return
	}
	logger.Info("getting collection", "collection_id: ", collectionID)
	col, err := a.db.GetCollectionById(collectionID)
	if err != nil {
		logger.Error("failed fetching collection", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("fetched collection successfully", "id: ", col.Id)
	its, err := a.db.GetItemsInCollection(col.Id)
	if err != nil {
		logger.Error("failed fetching items in collection", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("fetched items in collection successfully", "count: ", len(its))
	c.JSON(200, gin.H{
		"collection": col,
		"items":      its,
	})
}

func (a *App) loginUser(c *gin.Context) {
	var data models.UserLoginRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("logging in user", "user: ", data)
	user, err := a.db.GetUserByUsername(data.UserName)
	if err != nil {
		logger.Error("failed fetching user", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		logger.Error("invalid password", "err:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid credentials"})
		return
	}
	token, err := auth.SignPayload(map[string]interface{}{
		"user_id": user.Id,
		"username": user.UserName,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}, "testToken")
	if err != nil {
		logger.Error("failed signing token", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	c.SetCookie("Authorization", token, 3600*72, "/", "", false, true)
	c.JSON(200, gin.H{
		"token": token,
	})
}

func (a *App) sendOtp(c *gin.Context) {
	logger.Info("telegram login called")
	phone, ok := c.GetQuery("phone")
	if !ok || phone == "" {
		logger.Error("missing phone query param")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing phone query param"})
		return
	}
	logger.Info("sending login code to phone", "phone: ", phone)
	c.JSON(200, gin.H{
		"msg": "telegram login successful",
	})
}

func (a *App) verifyOtp(c *gin.Context) {
	logger.Info("verifying telegram otp")
	phone := c.Query("phone")
	code := c.Query("code")
	if phone == "" || code == "" {
		logger.Error("missing phone or code query param")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing phone or code query param"})
		return
	}
	logger.Info("verifying code", "phone: ", phone)
	token, err := auth.SignPayload(map[string]interface{}{
		"tele_id": phone,
		"exp":     time.Now().Add(10 * time.Hour).Unix(),
	}, "testToken")
	if err != nil {
		logger.Error("failed signing token", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	logger.Info("verified telegram otp successfully", "phone: ", phone)
	c.JSON(200, gin.H{
		"msg": "telegram verification successful",
		"tele_token": token,
	})
}

func (a *App) updateTelegram(c *gin.Context) {
	logger.Info("updating telegram details")
	userID, ok := c.Get("user_id")
	if !ok {
		logger.Error("missing user_id in context")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user not found"})
		return
	}
	token, ok := c.GetQuery("token")
	if !ok || token == "" {
		logger.Error("missing token query param")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing token query param"})
		return
	}
	details, err := auth.VerifyToken(token, "testToken")
	if err != nil {
		logger.Error("invalid token", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid token"})
		return
	}
	telegramId, ok := details["tele_id"].(string)
	if !ok || telegramId == "" {
		logger.Error("missing tele_id in token")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing tele_id in token"})
		return
	}
	logger.Info("updating telegram id", "tele_id: ", telegramId, "user_id: ", userID)
	tI, err := strconv.ParseInt(telegramId, 10, 64)
	if err != nil {
		logger.Error("invalid tele_id format", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid tele_id format"})
		return
	}
	err = a.db.AddTelegramToUser(userID.(string), tI)
	if err != nil {
		logger.Error("failed updating telegram id", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	c.JSON(200, gin.H{
		"msg": "telegram details updated successfully",
	})
}
