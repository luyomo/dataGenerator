package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

// Refer to: https://github.com/go-sql-driver/mysql/wiki/Examples
//           https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
func main() {
    db, err := sql.Open("mysql", "root@tcp(172.16.6.200:4000)/dev_tb6290_4")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    var columnName, dataType, columnType, columnKey sql.NullString
    var ordinalPosition, characterMaximumLength, numericPrecision, numericScale, datetimePrecision sql.NullInt64

    var queryStr string = "select COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION, COLUMN_TYPE, COLUMN_KEY  from information_schema.columns where table_schema = 'dev_tb6290_4' and table_name = 'SummaryCategoryTable' order by ORDINAL_POSITION"
    rows, err := db.Query(queryStr)
    defer rows.Close()
    fmt.Printf("Starting before fetching the data ")
    
    for rows.Next() {
        if err := rows.Scan(&columnName, &ordinalPosition, &dataType, &characterMaximumLength, &numericPrecision, &numericScale, &datetimePrecision, &columnType, &columnKey); err != nil {
            fmt.Printf("The square number of 13 is: %d", columnName)
        }
        if (columnName.Valid == true){
            fmt.Printf(" The error is %s \n", columnName.String)
        }else{
            fmt.Printf(" It:s unvalid string ")
        }

        if (datetimePrecision.Valid == true){
            fmt.Printf(" The error is %s \n", datetimePrecision.Int64)
        }else{
            fmt.Printf(" It:s unvalid string \n")
        }
    }

    if rows.Err() != nil {
        fmt.Printf("The error is ", rows.Err())
    }


//    // Prepare statement for reading data
//    stmtOut, err := db.Prepare("select COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION, COLUMN_TYPE, COLUMN_KEY  from information_schema.columns where table_schema = 'dev_tb6290_4' and table_name = 'SummaryCategoryTable' order by ORDINAL_POSITION")
//    if err != nil {
//        panic(err.Error()) // proper error handling instead of panic in your app
//    }
//    defer stmtOut.Close()
//
//
//    var rows
//    rows, err = stmtOut.Query()
//
//    // Query the square-number of 13
//    for rows.Next() {
//        err = rows.Scan(&columnName, &ordinalPosition, &dataType, &characterMaximumLength, &numericPrecision, &numericScale, &datetimePrecision, &columnType, &columnKey) 
//        if err != nil {
//            panic(err.Error()) // proper error handling instead of panic in your app
//        }
//        fmt.Printf("The square number of 13 is: %d", columnName)
//    }

    // Prepare statement for inserting data
    //stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
    //if err != nil {
    //    panic(err.Error()) // proper error handling instead of panic in your app
    //}
    //defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

    //// Prepare statement for reading data
    //stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
    //if err != nil {
    //    panic(err.Error()) // proper error handling instead of panic in your app
    //}
    //defer stmtOut.Close()
    //
    //// Insert square numbers for 0-24 in the database
    //for i := 0; i < 25; i++ {
    //    _, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
    //    if err != nil {
    //        panic(err.Error()) // proper error handling instead of panic in your app
    //	}
    //}
    //
    //var squareNum int // we "scan" the result in here
    //
    //// Query the square-number of 13
    //err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
    //if err != nil {
    //    panic(err.Error()) // proper error handling instead of panic in your app
    //}
    //fmt.Printf("The square number of 13 is: %d", squareNum)
    //
    //// Query another number.. 1 maybe?
    //err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
    //if err != nil {
    //    panic(err.Error()) // proper error handling instead of panic in your app
    //}
    //fmt.Printf("The square number of 1 is: %d", squareNum)
}
