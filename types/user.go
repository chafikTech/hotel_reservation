package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 8
)

type CreateUserParams struct {
	Firstname string `json:"firstname"`
	Lasttname string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.Firstname) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstname lenght should be at least %d characters", minFirstNameLen))
	}
	if len(params.Lasttname) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("Lastname lenght should be at least %d characters", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password lenght should be at least %d characters", minPasswordLen))
	}

	if !IsValidEmail(params.Email) {
		errors = append(errors, fmt.Sprintf("Email is not valid"))
	}
	return errors
}

func IsValidEmail(email string) bool {
	re := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(re).MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Firstname         string             `bson:"firstname" json:"firstname"`
	Lastname          string             `bson:"lastname" json:"lastname"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedpassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}

	return &User{
		Firstname:         params.Firstname,
		Lastname:          params.Lasttname,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
