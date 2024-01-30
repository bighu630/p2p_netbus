package handler

import (
	"net/http"
	"server/log"
	"server/user"

	"github.com/gin-gonic/gin"
)

type HeartBeat struct {
	userManager *user.UserManager
}

func NewHeartBeatHandler() *HeartBeat {
	return &HeartBeat{userManager: user.NewUserManager()}
}

func (h *HeartBeat) Handle(ctx *gin.Context) {
	userID := ctx.GetHeader("Authorization")

	address := ctx.Request.RemoteAddr
	log.Debug("get a user", "id", userID, "address", address)
	if u, exist := h.userManager.Users[userID]; exist {
		u.AddMachine(address)
	} else {
		h.userManager.Users[userID] = new(user.User)
		h.userManager.Users[userID].AddMachine(address)
	}
	jsonResponse(ctx, http.StatusOK, h.userManager.Users[userID].GetUserMachines())
}
