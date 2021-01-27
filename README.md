# URL shortener webservice

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
TODO:
* Generate short url.
* Open short url based on given generated short url.

---

## References
* [GoLang Docs](https://golangdocs.com/)

---

## Contact
Created by [Anas Juwaidi](mailto:anas.didi95@gmail.com)
