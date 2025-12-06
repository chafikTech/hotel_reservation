package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/chafikTeck/hotel_reservation/db"
	"github.com/chafikTeck/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))

	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "chafik@gmail.com",
		Firstname: "Mohamed",
		Lasttname: "chafik",
		Password:  "294nf!@nA1",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Error("expecting user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expecting EncryptedPassword not to be include in the json response")
	}

	if user.Firstname != params.Firstname {
		t.Errorf("expected Firstname %s but got %s", params.Firstname, user.Firstname)
	}
	if user.Lastname != params.Lasttname {
		t.Errorf("expected Lasttname %s but got %s", params.Lasttname, user.Lastname)
	}
	if user.Email != params.Email {
		t.Errorf("expecte Email %s but got %s", params.Email, user.Email)
	}

}
