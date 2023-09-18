package services

import (
	"context"
	"errors"

	"github.com/Hussein-miracle/golang-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (userServiceImpl *UserServiceImpl) CreateUser(user *models.User) error {
	_, err := userServiceImpl.usercollection.InsertOne(userServiceImpl.ctx, user)
	return err
}

func (userServiceImpl *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User

	query := bson.D{bson.E{Key: "name", Value: name}}

	err := userServiceImpl.usercollection.FindOne(userServiceImpl.ctx, query).Decode(&user)

	return user, err
}

func (userServiceImpl *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := userServiceImpl.usercollection.Find(userServiceImpl.ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	for cursor.Next(userServiceImpl.ctx) {
		var user *models.User
		err := cursor.Decode(&user)

		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(userServiceImpl.ctx)

	if len(users) == 0 {
		return nil, errors.New("No documents found")
	}
	return users, nil
}

func (userServiceImpl *UserServiceImpl) UpdateUser(user *models.User) error {

	filterQuery := bson.D{bson.E{Key: "user_name", Value: user.Name}}

	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{
		Key: "user_name", Value: user.Name,
	}, bson.E{
		Key: "user_age", Value: user.Age,
	}, bson.E{
		Key:   "user_address",
		Value: user.Address,
	}}}}

	result, _ := userServiceImpl.usercollection.UpdateOne(userServiceImpl.ctx, filterQuery, update)

	if result.MatchedCount != 1 {
		return errors.New("No matched found for update")
	}

	return nil
}

func (userServiceImpl *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{bson.E{Key: "user_name", Value: name}}

	result, _ := userServiceImpl.usercollection.DeleteOne(userServiceImpl.ctx, filter)

	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}

	return nil
}
