package mongo

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/smileinnovation/imannotate/api/user"
	"golang.org/x/crypto/bcrypt"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

var SigningKey = []byte("AllYourBase")

func init() {
	db := getMongo()
	defer db.Session.Close()

	idx := mgo.Index{
		Key:      []string{"email", "username"},
		Unique:   true,
		DropDups: true,
	}
	db.C("user").EnsureIndex(idx)

	pass, _ := bcrypt.GenerateFromPassword([]byte("toto123"), bcrypt.DefaultCost)
	db.C("user").Upsert(bson.M{"username": "Bob"}, &user.User{
		Username: "Bob",
		Password: string(pass),
	})
}

type MongoAuth struct{}
type CustomClaim struct {
	ID string
	jwt.StandardClaims
}

func (ma *MongoAuth) Login(u *user.User) error {
	db := getMongo()
	defer db.Session.Close()

	r := bson.M{"username": u.Username}
	real := &user.User{}
	if err := db.C("user").Find(r).One(real); err != nil {
		return errors.New("User from database isn't found or incorrect, " + err.Error())
	}

	// Password should be encrypted with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(real.Password), []byte(u.Password)); err != nil {
		return err
	}

	claims := CustomClaim{
		u.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(time.Hour * 24)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if ss, err := token.SignedString(SigningKey); err != nil {
		return err
	} else {
		u.Token = "Bearer " + ss
	}

	return nil
}
func (ma *MongoAuth) Logout(u *user.User) error {
	return nil
}
func (ma *MongoAuth) Allowed(req *http.Request) error {

	bearer := req.Header.Get("Authorization")
	bearer = strings.Replace(bearer, "Bearer ", "", 1)

	_, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	log.Println(err)
	return err
}

func (ma *MongoAuth) GetCurrentUsername(req *http.Request) (string, error) {
	bearer := req.Header.Get("Authorization")
	bearer = strings.Replace(bearer, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(bearer, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	claim := token.Claims.(*CustomClaim)
	return claim.ID, nil
}