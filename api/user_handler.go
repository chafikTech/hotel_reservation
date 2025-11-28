package api

import (
	"context"

	"github.com/chafikTeck/hotel_reservation/db"
	"github.com/chafikTeck/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		Firstname: "Mohamed",
		Lasttname: "Chafik",
	}
	return c.JSON(u)
}
