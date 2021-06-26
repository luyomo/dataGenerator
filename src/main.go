package main

import (
    "database/sql"
    "fmt"
    "time"
    "strings"
    "math/rand"
    "math"
    "strconv"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Refer to: https://github.com/go-sql-driver/mysql/wiki/Examples
//           https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
func main() {
    rand.Seed(time.Now().UnixNano())

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

    var listTableMeta []TableMeta
    
    for rows.Next() {
        if err := rows.Scan(&tableMeta.columnName, &tableMeta.ordinalPosition, &tableMeta.dataType, &tableMeta.characterMaximumLength, &tableMeta.numericPrecision, &tableMeta.numericScale, &tableMeta.datetimePrecision, &tableMeta.columnType, &tableMeta.columnKey); err != nil {
            fmt.Printf("The square number of 13 is: %d", tableMeta.columnName)
        }
        listTableMeta = append(listTableMeta, tableMeta)

        //if (tableMeta.dataType.Valid == true) {
        //    fmt.Printf("columnType is :  %s \n", tableMeta.dataType.String)
        //    switch tableMeta.dataType.String {
        //        case "varchar" : generateString(10)
        //    }
        //}else{
        //    fmt.Printf(" column name: It:s unvalid string \n")
        //}



        //if (tableMeta.columnName.Valid == true){
        //    fmt.Printf(" columnName:  %s \n", tableMeta.columnName.String)
        //}else{
        //    fmt.Printf(" column name: It:s unvalid string \n")
        //}

        //if (tableMeta.datetimePrecision.Valid == true){
        //    fmt.Printf(" The date time is  %s \n", tableMeta.datetimePrecision.Int64)
        //}else{
        //    fmt.Printf(" datetimePrecison \n")
        //}
    }

    var colValue string
    for _, tableMeta := range listTableMeta {
        if (tableMeta.dataType.Valid == true) {
            switch tableMeta.dataType.String {
                case "varchar" : 
                    colValue = generateString(int(tableMeta.characterMaximumLength.Int64))
                    fmt.Printf("               The string value is <%s> \n", colValue)
                case "int":
                    colValue = generateInt()
                    fmt.Printf("               The int value is <%s> \n", colValue)
                case "char":
                    colValue = generateChar(int(tableMeta.characterMaximumLength.Int64))
                    fmt.Printf("               The char value is <%s> \n", colValue)
                case "decimal":
                    colValue = generateDecimal(int(tableMeta.numericPrecision.Int64), int(tableMeta.numericScale.Int64))
                    fmt.Printf("               The decimal value is <%s> \n", colValue)
                default:
                    fmt.Printf("columnType is :  %s \n", tableMeta.dataType.String)
            }
        }
    }


    if rows.Err() != nil {
        fmt.Printf("The error is ", rows.Err())
    }

    var data []string
    data = append(data, "\"value 01\"")
    data = append(data, "value 02")
    data = append(data, "value 03")
    data = append(data, "value 04")
    fmt.Printf("The error is %s ", strings.Join(data, ","))
}

func generateData(tableMeta TableMeta, numRows int){
}

// Generate the random string
func generateString(maxSize int) string {
    b := make([]rune, rand.Intn(maxSize+1))
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

// Generate random int
func generateInt() string {
    return strconv.Itoa(rand.Intn(2147483647))
}
// Generate the random string
func generateChar(size int) string {
    b := make([]rune, size)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

// Generate random decimal
func generateDecimal(precision int, scale int) string {
    power10 := math.Pow10(scale)
    randInt := rand.Intn(int(math.Pow10(precision - scale))) - 1

    return fmt.Sprintf("%d.%d", randInt,  int(math.Round(rand.Float64()*power10)))
}
