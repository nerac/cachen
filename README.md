# Cachen

[![Build Status](https://travis-ci.org/nerac/cachen.svg?branch=master)](https://travis-ci.org/nerac/cachen)
[![Coverage Status](https://coveralls.io/repos/github/nerac/cachen/badge.svg?branch=master)](https://coveralls.io/github/nerac/cachen?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nerac/cachen)](https://goreportcard.com/report/github.com/nerac/cachen)
[![GoDoc](https://godoc.org/github.com/nerac/cachen?status.svg)](https://godoc.org/github.com/nerac/cachen)

Golang library that simplify the way you manage the http cache.

## Usage

Install the library with:

`$ go get github.com/nerac/cachen`

Import library into your code with as in the next example:

    package main

    import (
        "fmt"
        "net/http"
        "github.com/nerac/cachen"
    )
    func homeHandler(w http.ResponseWriter, r *http.Request) {
        cache := cachen.New()

        etag := "my etag"

        cache.ReusableRequest(true).
        RevalidateEachTime(false).
        IntermediatesAllowed(true).
        MaxAge(5 * SECONDS).
        StaleAllowed(true).
        Etag(etag).
        Bind(w,r)

        fmt.Fprintf(w, "This content is cached 5 seconds by all parties!")
    }
    func main() {

        http.HandleFunc("/", homeHandler)
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
        }
    }