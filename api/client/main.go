package main

import (
	"fmt"
	"genesis/pkg/api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := api.NewCreateUserClient(conn)

	g := gin.Default()
	g.GET("/user/auth", func(ctx *gin.Context) {

		email := ctx.Query("email")
		pass := ctx.Query("pass")
		req := &api.UserRequest{Email: email, Pass: pass}

		if response, err := client.Auth(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(response.Response)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	g.GET("/user/create", func(ctx *gin.Context) {
		a := ctx.Param("email")
		b := ctx.Param("pass")

		req := &api.UserRequest{Email: a, Pass: b}

		if response, err := client.Create(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(response.Response)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	g.GET("/:a", func(ctx *gin.Context) {
		a := ctx.Param("a")
		req := &api.URL{X: a}
		//
		if response, err := client.GetBTC(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf(response.Response)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	if err := g.Run(":8000"); err != nil {
		log.Fatal("fail")
	}
}
