package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type TableMeta struct {
    columnName              sql.NullString     `json:"COLUMN_NAME"`
    ordinalPosition         sql.NullInt64      `json:"ORDINAL_POSITION"`
    dataType                sql.NullString     `json:"DATA_TYPE"`
    characterMaximumLength  sql.NullInt64      `json:"CHARACTER_MAXIMUM_LENGTH"`
    numericPrecision        sql.NullInt64      `json:"NUMERIC_PRECISION"`
    numericScale            sql.NullInt64      `json:"NUMERIC_SCALE"`
    datetimePrecision       sql.NullInt64      `json:"DATETIME_PRECISION"`
    columnType              sql.NullString     `json:"COLUMN_TYPE"`
    columnKey               sql.NullString     `json:"COLUMN_KEY"`
}


// Refer to: https://github.com/go-sql-driver/mysql/wiki/Examples
//           https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
func main() {
    db, err := sql.Open("mysql", "root@tcp(172.16.6.200:4000)/dev_tb6290_4")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    var tableMeta TableMeta

    var queryStr string = "select COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION, COLUMN_TYPE, COLUMN_KEY  from information_schema.columns where table_schema = 'dev_tb6290_4' and table_name = 'SummaryCategoryTable' order by ORDINAL_POSITION"
    rows, err := db.Query(queryStr)
    defer rows.Close()
    fmt.Printf("Starting before fetching the data ")
    
    for rows.Next() {
        if err := rows.Scan(&tableMeta.columnName, &tableMeta.ordinalPosition, &tableMeta.dataType, &tableMeta.characterMaximumLength, &tableMeta.numericPrecision, &tableMeta.numericScale, &tableMeta.datetimePrecision, &tableMeta.columnType, &tableMeta.columnKey); err != nil {
            fmt.Printf("The square number of 13 is: %d", tableMeta.columnName)
        }
        if (tableMeta.columnName.Valid == true){
            fmt.Printf(" columnName:  %s \n", tableMeta.columnName.String)
        }else{
            fmt.Printf(" column name: It:s unvalid string \n")
        }

        if (tableMeta.datetimePrecision.Valid == true){
            fmt.Printf(" The date time is  %s \n", tableMeta.datetimePrecision.Int64)
        }else{
            fmt.Printf(" datetimePrecison \n")
        }
    }

    if rows.Err() != nil {
        fmt.Printf("The error is ", rows.Err())
    }
}
