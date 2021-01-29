# URL shortener webservice

![build](https://github.com/anas-didi95/golang-urlshort-server/workflows/build/badge.svg)

---

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Environment Variables](#environment-variables)
* [Setup](#setup)
* [Features](#features)
* [References](#references)
* [Contact](#contact)

---

## General info
Back-end service which provides method to generate short url.

---

## Technologies
* Go - Version 1.15

---

## Environment Variables
Following table is a **mandatory** environment variables used in this project.

| Variable Name | Datatype | Description |
| --- | --- | --- |
| APP_HOST | String | Server host |
| APP_PORT | Number | Server port |
| MONGO_CONNECTION_STRING | String | Mongo connection string (refer [doc](https://docs.mongodb.com/manual/reference/connection-string/) for example)

---

## Setup
This project contains VSCode container. Simple install `Docker`, `Remote - Containers` extensions to open the workspace.

To launch your tests:
```
go test
```

To build your application:
```
go build main.go
```

To run your application:
```
go run main.go

# Or after build
./main
```

---

## Features
* Generate short url.

TODO:
* Open short url based on given generated short url.

---

## References
* [GoLang Docs](https://golangdocs.com/)
* [How I write HTTP services after eight years](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html)
* [Get JSON request body from a HTTP request in Go (Golang)](https://golangbyexample.com/json-request-body-golang-http/)

---

## Contact
Created by [Anas Juwaidi](mailto:anas.didi95@gmail.com)
