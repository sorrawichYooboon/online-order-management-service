package main

import (
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/http"
)

func main() {
	e := http.NewEchoServer()

	e.Logger.Fatal(e.Start(":8080"))
}
