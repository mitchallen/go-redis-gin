/**
 * Author: Mitch Allen
 * File: server.go
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchallen/go-redis-gin/demo"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK!",
		})
	})

	/*

		curl http://localhost:8080/lock/resource/alpha

	*/

	r.GET("/lock/resource/:resource", func(c *gin.Context) {
		resource := c.Param("resource")
		key := demo.MakeLockKey(resource)
		val, err := client.Get(c, key).Result()
		if len(val) == 0 { // or val == ""
			// if an empty sting was retuned, the key was not found
			fmt.Println("--- key not found ---")
		}
		if err != nil {
			if err == redis.Nil {
				// If the error was redis.Nil, the key was not found
				fmt.Printf("--- GET returned redis.Nil, err: %v ---\n", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "resource lock for user not found"})
				return
			}
			// Otherwise an unexpected error occurred
			fmt.Printf("ERROR [Get]: %v \n", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(200)
		bytes := []byte(val)
		c.Writer.Write(bytes)
	})

	/*
		curl -X POST http://localhost:8080/lock \
		-H 'content-type: application/json' \
		-d '{"resource": "alpha", "userId": "admin"}'
	*/

	r.POST("/lock", func(c *gin.Context) {

		type LockPostBody struct {
			Resource string `json:"resource"`
			UserID   string `json:"userId"`
		}

		var requestBody LockPostBody

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resource := requestBody.Resource
		userId := requestBody.UserID

		duration := time.Second * 10

		lock := demo.Lock{
			Resource: resource,
			UserID:   userId,
			Duration: duration.String(),
		}

		json, err := json.Marshal(lock)

		if err != nil {
			// error marshalling record
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		key := demo.MakeLockKey(resource)

		err = client.Set(c, key, json, duration).Err()
		if err != nil {
			// error setting key
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, lock)
	})

	// listen and serve on 0.0.0.0:8080
	// on windows "localhost:8080"
	// can be overriden with the PORT env var
	r.Run()
}
