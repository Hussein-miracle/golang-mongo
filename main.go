package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Hussein-miracle/golang-mongo/controllers"
	"github.com/Hussein-miracle/golang-mongo/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	server         *gin.Engine
	userservice    services.UserService
	usercontroller controllers.UserController
	ctx            context.Context
	usercollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.Please add")
	}

	ctx = context.TODO()

	mongoconnection := options.Client().ApplyURI(uri)
	mongoclient, err = mongo.Connect(ctx, mongoconnection)

	if err != nil {
		log.Fatal(err)
	}

	// defer func() {
	// 	if err := mongoclient.Disconnect(ctx); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// Send a ping to confirm a successful connection
	if err := mongoclient.Database("mongo-users").RunCommand(ctx, bson.D{bson.E{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	fmt.Println("connection established")

	usercollection = mongoclient.Database("mongo-users").Collection("users")
	fmt.Println("collection retrieved")

	userservice = services.NewUserService(usercollection, ctx)
	fmt.Println("services retrieved")
	usercontroller = controllers.New(userservice)
	fmt.Println("controllers retrieved")
	server = gin.Default()
}

// "/v1/"
func main() {
	defer mongoclient.Disconnect(ctx)
	basePath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basePath)
	server.Run("localhost:3030")
}
