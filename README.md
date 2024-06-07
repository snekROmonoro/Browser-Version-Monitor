# Browser-Version-Monitor

Monitor browser versions release and alert

## About

Made to monitor changes to browser versions. Uses PostgresSQL.

## Setup

First step is to make sure your go.mod is tidied:

```sh
go mod tidy
```

## Database Setup

The database is self-hosted PostgresSQL.

*cba explain how you host that so just google about it sob*

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

To look into the database:

```sh
go run github.com/steebchen/prisma-client-go studio

prisma studio
```
