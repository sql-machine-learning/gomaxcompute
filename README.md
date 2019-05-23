# MaxCompute Go Driver

[MaxCompute](https://www.alibabacloud.com/product/maxcompute), also known as ODPS, is a distributed storage service and SQL engine provided by [Alibaba Cloud](www.alibabacloud.com/). This repository contains a Go [SQLdriver](https://github.com/golang/go/wiki/SQLDrivers) of MaxCompute.  If you are going to write a Go program that calls the standard library `database/sql` to access MaxCompute databases, you could use this driver.

This project is in its early stage. Your issues and pull requests are very welcome!


## What This Is and Isn't

This project is a driver that helps Go's standard database API talking to MaxCompute server. It has the following features:
 
- In pure Go. Not a wrapper of any C/C++ library.
- Connect to MaxCompute through its [HTTP interface](http://repo.aliyun.com/api-doc/).
- Improve I/O throughput using MaxCompute's [tunnel service](https://www.alibabacloud.com/help/doc-detail/27833.htm).

Alibaba Cloud open sourced some client SDKs of MaxCompute:

- Java: https://github.com/aliyun/aliyun-odps-java-sdk
- Python: https://github.com/aliyun/aliyun-odps-python-sdk

This project is not an SDK.

Alibaba Cloud also provides ODBC/JDBC drivers:

- https://github.com/aliyun/aliyun-odps-jdbc

This project is a Go's `database/sql` driver.


## How to Use

Please make sure you have Go 1.6 or high release. 

You can clone the source code by running the following command.

```go
go get -u sqlflow.org/gomaxcompute
```

Here is a simple example:

```go
package main

import (
    "database/sql"
    "sqlflow.org/gomaxcompute"
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

Please be aware that to connect to a MaxCompute database, the user needs to provide 

1. the access ID
1. the access key
1. the endpoint pointing to the MaxCompute service
1. a project, which is something similar to a database in MySQL.

## Acknowledgement

Our respect and thanks to Ruohang Feng, who wrote a Go SDK for MaxCompute when he worked in Alibaba, for his warm help that enabled this project.
