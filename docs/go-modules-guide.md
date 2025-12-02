# Go Modules Guide

## Introduction

Go modules define project dependencies and enable reproducible builds.

## Key Concepts

* `go.mod` — module path + dependency versions.
* `go.sum` — verification checksums.
* Semantic import versioning for major v2+.

## Common Commands

* `go mod init <module>`
* `go get <package>`
* `go mod tidy`
* `go list -m all`

## Best Practices

* Keep module paths stable.
* Avoid renaming modules.
* Use tags for versioning.
