##  Go RBAC example

This is a simple example of Role Based Access Control (RBAC) with Go and MongoDB.

It consists of a REST API that exposes endpoints to manage users, roles and permissions.

The API is built using the pure golang standard library with these features:

- Project based on Port and Adapters architecture
- Data model based on mongodb external's reference
- Design pattern as Builder pattern
- Authentication with OTP and Paseto (jwt alternative) with invalidate token
- Authorization with RBAC (role based access controll) pattern
- Goroutines and more.


### usage

1. install golang 1.22+
2. clone repository 
3. setup .env file
4. and run!

```
git clone https://github.com/jmjp/go-rbac.git
cd go-rbac
go run cmd/main.go
```