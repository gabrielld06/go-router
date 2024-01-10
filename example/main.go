package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	router "github.com/gabrielld06/go-router/cmd"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	r := router.New(&http.Server{
		Addr: ":3000",
	})

	user := User{
		"user",
		"user@mail.com",
		22,
	}

	r.Get("/user", func(r *http.Request) (router.ApiResponse, error) {
		return router.ApiResponse{
			Status: http.StatusOK,
			Data:   user,
		}, nil
	})

	r.Post("/user", func(r *http.Request) (router.ApiResponse, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return router.ApiResponse{}, errors.New("invalid body")
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			return router.ApiResponse{}, errors.New("invalid body")
		}

		return router.ApiResponse{
			Status: http.StatusOK,
			Data:   user,
		}, nil
	})

	r.UseGlobalMiddleware(func(r *http.Request) error {
		if authorization := r.Header.Get("Authorization"); authorization != "token 123" {
			return errors.New("forbidden")
		}

		return nil
	})

	r.UseMiddleware("/user", func(r *http.Request) error {
		if r.Method == "POST" && r.Header.Get("Content-Type") != "application/json" {
			return errors.New("invalid content-type")
		}

		return nil
	})

	r.UseErrorHandler("invalid body", func(w http.ResponseWriter) {
		http.Error(w, "Invalid Body", http.StatusUnprocessableEntity)
	})

	r.UseErrorHandler("forbidden", func(w http.ResponseWriter) {
		http.Error(w, "forbidden", http.StatusForbidden)
	})

	r.UseErrorHandler("invalid content-type", func(w http.ResponseWriter) {
		http.Error(w, "Content-Type should be application/json", http.StatusBadRequest)
	})

	log.Fatal(r.Start())
}
