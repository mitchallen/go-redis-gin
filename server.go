/**
 * Author: Mitch Allen
 * File: server.go
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Lock struct {
	Resource string `json:"resource"`
	UserID   string `json:"userId"`
	Duration string `json:"duration"`
}

const NAMESPACE = "lock"

func makeKey(namespace string, location string) string {
	return fmt.Sprintf(
		"%s:%s",
		strings.ToLower(namespace),
		strings.ToLower(location),
	)
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	r.GET("/lock/:resource/:userid", func(c *gin.Context) {
		resource := c.Param("resource")
		userId := c.Param("userid")
		duration := time.Second * 5

		lock := Lock{
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

		key := makeKey(NAMESPACE, resource)

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
