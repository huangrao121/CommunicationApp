package user

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huangrao121/CommunicationApp/BackendService/config/pkg"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/types"
)

type UserHandler struct {
	userStore *UserStore
}

func NewUserHandler(userStore *UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user types.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := h.userStore.db.Create(&user)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		return
	}

	slog.Info("User created", "user", user)

	token, err := pkg.GenerateJWKToken(&user, nil, os.Getenv("PK_PATH"), time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	acl := getACL(user.ID.String())

	mqttToken, err := pkg.GenerateJWKToken(&user, acl, os.Getenv("MQTT_PK_PATH"), 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("app_jwt", token, 24*3600, "/", "localhost", false, true)
	c.SetCookie("mqtt_jwt", mqttToken, 15*60, "/", "localhost", false, true)
	SignupResp := types.LoginResp{
		AppJwt:  token,
		MqttJwt: mqttToken,
	}
	c.JSON(http.StatusCreated, SignupResp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var user types.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := h.userStore.db.First(&user, "email = ?", user.Email, "password = ?", user.Password)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		return
	}
	if r.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := pkg.GenerateJWKToken(&user, nil, os.Getenv("PK_PATH"), time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	acl := getACL(user.ID.String())

	mqttToken, err := pkg.GenerateJWKToken(&user, acl, os.Getenv("MQTT_PK_PATH"), 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	loginResp := types.LoginResp{
		AppJwt:  token,
		MqttJwt: mqttToken,
	}
	c.SetCookie("app_jwt", token, 24*3600, "/", "localhost", false, true)
	c.SetCookie("mqtt_jwt", mqttToken, 15*60, "/", "localhost", false, true)
	c.JSON(http.StatusOK, loginResp)
}

// 刷新短时间的mqtt token，用于mqtt连接
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID, okId := c.Get("user_id")
	userName, okName := c.Get("user_name")
	email, okEmail := c.Get("user_email")
	if !okId || !okName || !okEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	acl := getACL(userID.(string))
	refreshToken, err := pkg.GenerateJWKToken(&types.Users{ID: userID.(uuid.UUID), Username: userName.(string), Email: email.(string)}, acl, os.Getenv("PK_PATH"), time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("mqtt_jwt", refreshToken, 15*60, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"mqtt_token": refreshToken})
}

// TODO: add group acl
func getACL(userID string) *[]types.ACL {
	acl := &[]types.ACL{
		{
			Permission: "allow",
			Action:     "subscribe",
			Topic:      "users/" + userID + "/inbox",
		},
		{
			Permission: "allow",
			Action:     "subscribe",
			Topic:      "users/" + userID + "/presence",
		},
		{
			Permission: "allow",
			Action:     "subscribe",
			Topic:      "users/" + userID + "/cmd",
		},
		{
			Permission: "allow",
			Action:     "publish",
			Topic:      "chats/p2p/" + userID,
		},
		{
			Permission: "allow",
			Action:     "subscribe",
			Topic:      "chats/group/" + userID,
		},
	}
	return acl
}
