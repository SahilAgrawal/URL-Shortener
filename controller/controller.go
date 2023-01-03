package controller

import (
	"URL-Shortener/model"
	"URL-Shortener/repository"
	serviceinterface "URL-Shortener/service/ServiceInterface"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type handler struct {
	redirectServiceInterface serviceinterface.RedirectServiceInterface
}

func NewHandler(repo repository.RedirectRepo) *handler {
	return &handler{
		redirectServiceInterface: repo,
	}
}

func (h *handler) Find(ctx *gin.Context) {

	code := ctx.Param("id")
	data, err := h.redirectServiceInterface.Find(code)
	if err != nil {
		ctx.IndentedJSON(http.StatusFound, gin.H{
			"message": "Data Not Found for given id: " + code,
		})
	}
	ctx.Redirect(http.StatusMovedPermanently, data.URL)
	// ctx.IndentedJSON(http.StatusFound, data)

}

func (h *handler) Store(ctx *gin.Context) {
	var redirect *model.Redirect

	ctx.BindJSON(&redirect)
	err := h.redirectServiceInterface.Store(redirect)
	fmt.Println(redirect)
	if err != nil {
		ctx.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "URL can not able to short, may be already presented",
			"error":   err.Error(),
		})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, redirect)
}

func (h *handler) All(ctx *gin.Context) {
	data, err := h.redirectServiceInterface.All()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, bson.M{
			"msg":    err,
			"status": http.StatusInternalServerError,
		})
	}

	ctx.IndentedJSON(http.StatusOK, data)
}
