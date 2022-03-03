# TempShare
A website for anonymously and securely sharing text data with links that are destroyed after being viewed.

## Creating a file: 
![image](https://user-images.githubusercontent.com/8042849/155024070-5d1e7c6e-baf8-470c-8f07-23801d6b2352.png)

## Viewing a file:
![image](https://user-images.githubusercontent.com/8042849/155024163-7d12e4bc-4ca3-4ca2-80a1-316405ce23e2.png)


This can be run directly via:
> go run cmd/web/!(*_test.go)

Or compiled into an executable
> go build -o server cmd/web/!(*_test.go)

> ./server

Tests can be run from the project root folder via:
> go test -v ./...
