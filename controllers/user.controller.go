package controllers

import (
	"net/http"

	"github.com/Hussein-miracle/golang-mongo/models"
	"github.com/Hussein-miracle/golang-mongo/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func New(userservice services.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

func (uc UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	// _,err := uc.UserService.CreateUser()
	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := uc.UserService.CreateUser(&user)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "success"})
}

func (uc UserController) GetUser(ctx *gin.Context) {
	name := ctx.Param("name")

	user, err := uc.UserService.GetUser(&name)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)

}

func (uc UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}

func (uc UserController) UpdateUser(ctx *gin.Context) {
	var user *models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := uc.UserService.UpdateUser(user)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc UserController) DeleteUser(ctx *gin.Context) {
	username := ctx.Param("name")
	err := uc.UserService.DeleteUser(&username)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroutes := rg.Group("/user")

	userroutes.POST("/create", uc.CreateUser)
	userroutes.GET("/get/:name", uc.GetUser)
	userroutes.GET("/getall", uc.GetAll)
	userroutes.PATCH("/update", uc.UpdateUser)
	userroutes.DELETE("/delete", uc.DeleteUser)
}
