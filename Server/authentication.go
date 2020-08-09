package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	log2 "log"
	"math/rand"
	"time"
)

type user struct {
	Username string `bson:"username,omitempty"`
	Password string `bson:"password,omitempty"`
}

var contextCancel context.CancelFunc
var userCollection *mongo.Collection

func initDBClient(mongoURI string) {
	var ctx context.Context
	mdbClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log2.Panic(err)
	}
	ctx, contextCancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = mdbClient.Connect(ctx)
	if err != nil {
		log2.Fatal(err)
	}
	err = mdbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log2.Fatal(err)
	}
	userCollection = mdbClient.Database("mancala").Collection("users")
}

func disconnectDBClient() {
	contextCancel()
	userCollection = nil
}

func generateKey() string {
	result := ""
	for i := 0; i < 4; i++ {
		result += fmt.Sprintf("%x", rand.Int63())
	}
	return result
}

func checkAuthentication(username, signature string) (bool, error) {
	if !userExists(username) {
		return false, errors.New("corrupted cookie")
	}
	hmc := hmac.New(sha256.New, CookieKey)
	hmc.Write([]byte(username))
	valid := hmc.Sum(nil)

	decodedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	if !hmac.Equal(valid, decodedSignature) {
		println(valid)
		println(decodedSignature)
		return false, errors.New("corrupted cookie")
	}
	return true, nil
}

func sign(message string) string {
	hmc := hmac.New(sha256.New, CookieKey)
	hmc.Write([]byte(message))
	sum := hmc.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}

func checkPassword(username, password string) bool {
	filter := bson.D{{"username", username},
		{"password", password}}
	cursor, err := userCollection.Find(context.TODO(), filter)
	if err != nil {
		log2.Print("Error while checking if user exists")
		return false //TODO process error better
	}
	defer cursor.Close(context.TODO())
	if cursor.Next(context.TODO()) {
		return true
	}
	return false
}

func userExists(username string) bool {
	filter := bson.D{{"username", username}}
	cursor, err := userCollection.Find(context.TODO(), filter)
	if err != nil {
		log2.Print("Error while checking if user exists")
		return true //TODO process error better
	}
	defer cursor.Close(context.TODO())
	if cursor.Next(context.TODO()) {
		return true
	}
	return false
}
func addUser(username, password string) {
	newUser := user{
		Username: username,
		Password: password,
	}
	insertResult, err := userCollection.InsertOne(nil, newUser)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)
}
