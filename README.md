# MaxCompute Go Driver

This repo contains a Go driver for MaxCompute(ODPS) by AliCloud. If your Go projects involve with talking to MaxCompute server, you are welcome to try it out and give us some feedback. Note this implementation contains a database abstraction layer `database/sql` compatible with Go standard library using native HTTP API. As you can see from version number, this repo is under fast iteration, please create issue if you find anything missing.

## Features
- Light weight and fast
- Native Go implementation, compatible with Go interface `database/sql`
- Connection through [HTTP API](http://repo.aliyun.com/api-doc/)
- Supports large queries via MaxCompute tunnel server

## How to use MaxComputer Go Driver
1. Installation
Please make sure you have Go 1.6 or high release. 

```go
go get -u github.com/sql-machine-learning/gomaxcompute
```

Because Go program normally requires a connection string to talk to the select database, we build the string from `access_id`,`access_key`,`endpoint`,`project` in the form of `http://<access_id>:<access_key>@<endpoint>/api?curr_project=<project>`.

2. Example

```go
package main

import (
	"database/sql"
	"github.com/sql-machine-learning/maxcompute"
)

func assertNoError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	config := goodps.Config{
		AccessID:  "<access_id>",
		AccessKey: "<access_key>",
		Endpoint:  "<end_point>",
		Project:   "<project_name>"}
	db, e := sql.Open("maxcompute", config.FormatDSN())
	assertNoError(e)
	defer db.Close()

	const sql = `SELECT
                    cast('1'                   AS BIGINT)  AS a,
                    cast(TRUE                  AS BOOLEAN) AS b,
                    cast('hi'                  AS STRING)  AS c,
                    cast('3.14'                AS DOUBLE)  AS d,
                    cast('2017-11-11 03:12:11' AS DATETIME) AS e,
                    cast('100.01' AS DECIMAL)  AS f;`
	rows, e := db.Query(sql)
	assertNoError(e)
	defer rows.Close()

	for rows.Next() {
		// do your stuff
	}
}
```

## Contributions
Thanks Ruohang Feng, the earlist contributor while at Alibaba. SQLFlow team from Ant Financial also contributes to this codebase.
