package main

import (
	"log"
	"time"

	"github.com/ejuju/go-cicd"
)

func main() {
	start := time.Now()
	defer func() { log.Printf("Time elapsed since start: %s\n", time.Since(start)) }()

	err := cicd.NewRunner().Run(
		cicd.NewStep("Run Go code checks",
			cicd.Exec("go mod tidy"),
			cicd.Exec("go mod verify"),
			cicd.Exec("go vet ./..."),
			cicd.Exec("go build -o /dev/null"),
		),
		cicd.NewStep("Run unit tests",
			cicd.SetEnv("CGO_ENABLED", "1"),
			cicd.Exec("go test ./... -cover -race -timeout=60s"),
		),
	)
	if err != nil {
		panic(err)
	}
}
