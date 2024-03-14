# Browser-Version-Monitor
Monitor browser versions release and alert

## About

Made to monitor changes to browser versions. Uses PostgresSQL.

## Setup

First step is to make sure your go.mod is up-to-date:

```sh
go mod tidy
```

## Prisma

When changing something in the database:

```sh
go run github.com/steebchen/prisma-client-go db push
```

When something was changed in the database:

```sh
go run github.com/steebchen/prisma-client-go db pull
```

To generate the database files:

```sh
go run github.com/steebchen/prisma-client-go generate
```
