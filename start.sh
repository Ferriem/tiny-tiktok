#!/bin/bash

go mod tidy

# Start the first process

go run service/social_service/cmd/main.go
go run service/user_service/cmd/main.go
go run service/video_service/cmd/main.go

go run api_router/cmd/main.go

