package main

import (
    "database/sql"
    "fmt"
    "time"
    "strings"
    "math/rand"
    "math"
//    "strconv"
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
    extra                   sql.NullString     `json:"EXTRA"`
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

    var queryStr string = "select COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION, COLUMN_TYPE, COLUMN_KEY, EXTRA  from information_schema.columns where table_schema = 'dev_tb6290_4' and table_name = 'TransactionHead' order by ORDINAL_POSITION"
    rows, err := db.Query(queryStr)
    defer rows.Close()
    //fmt.Printf("Starting before fetching the data ")

    var listTableMeta []TableMeta
    
    for rows.Next() {
        if err := rows.Scan(&tableMeta.columnName, &tableMeta.ordinalPosition, &tableMeta.dataType, &tableMeta.characterMaximumLength, &tableMeta.numericPrecision, &tableMeta.numericScale, &tableMeta.datetimePrecision, &tableMeta.columnType, &tableMeta.columnKey, &tableMeta.extra); err != nil {
            fmt.Printf("The square number of 13 is: %d", tableMeta.columnName)
        }
        listTableMeta = append(listTableMeta, tableMeta)
    }
    if rows.Err() != nil {
        fmt.Printf("The error is ", rows.Err())
    }

    // Original smt
    smt := "insert into %s.%s(%s) values %s"

    // Generate the columns list
    colStrings := []string{}
    arrMarks := []string{}
    for _, tableMeta := range listTableMeta {
        if tableMeta.extra.Valid == false || tableMeta.extra.String != "auto_increment" {
            colStrings = append(colStrings, tableMeta.columnName.String)
            arrMarks = append(arrMarks, "?")
        }
    }
    strMarks := "(" + strings.Join(arrMarks, ",") + ")"
    strCols := strings.Join(colStrings, ",")
    //fmt.Printf("All the columns are %s \n", smt)
    //fmt.Printf("The marks is %s \n", strMarks)


    for range [1]int{} {
        valueStrings := []string{}
        valueArgs := []interface{}{}
        for range [1]int{} {
            //fmt.Printf("\n\n")
    
            valueStrings = append(valueStrings, strMarks)
            
            for _, tableMeta := range listTableMeta {
                if (tableMeta.dataType.Valid == true) {
                    if tableMeta.extra.Valid == true && tableMeta.extra.String == "auto_increment" {
                        continue
                    }
                    switch tableMeta.dataType.String {
                        case "varchar" : 
                            valueArgs = append(valueArgs, generateString(int(tableMeta.characterMaximumLength.Int64)))
                        case "text" : 
                            valueArgs = append(valueArgs, generateString(10000))
                        case "int":
                            valueArgs = append(valueArgs, generateInt())
                        case "tinyint":
                            valueArgs = append(valueArgs, generateTinyint())
                        case "bigint":
                            valueArgs = append(valueArgs, generateBigint())
                        case "char":
                            valueArgs = append(valueArgs,generateChar(int(tableMeta.characterMaximumLength.Int64)))
                        case "decimal":
                            valueArgs = append(valueArgs, generateDecimal(int(tableMeta.numericPrecision.Int64), int(tableMeta.numericScale.Int64)))
                        case "datetime":
                            valueArgs = append(valueArgs, generateDatetime())
                        default:
                            fmt.Printf("columnType is :  %s \n", tableMeta.dataType.String)
                    }
                }
            }
        }
        insertSmt := fmt.Sprintf(smt, "dev_tb6290_4", "TransactionHead",  strCols, strings.Join(valueStrings, ","))
        //fmt.Printf("The query is %s", insertSmt)
    
        tx, err := db.Begin()
        if err != nil {
            fmt.Printf("Error on starting transaction ")
            fmt.Println(err)
        }
        defer tx.Rollback()
        stmt, err := tx.Prepare(insertSmt)
        if err != nil {
            fmt.Printf("Error on preparing query")
            fmt.Println(err)
        }
        //fmt.Println(valueArgs)
        defer stmt.Close() // danger!
        	_, err = stmt.Exec(valueArgs...)
        	if err != nil {
                    fmt.Printf("bbbbbbbaaaaaa--------")
                    fmt.Println(err)
        	}
        err = tx.Commit()
        if err != nil {
                    fmt.Printf("xxxxxxxxxxxaa--------")
        }
    }

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
func generateInt() int {
    return rand.Intn(2147483647)
}

// Generate random int
func generateBigint() int64 {
    return rand.Int63n(9223372036854775807)
}

// Generate random tinyint
func generateTinyint() int {
    return rand.Intn(127)
}

// Generate the random string
func generateChar(size int) string {
    b := make([]rune, size)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

// Generate random datetime
func generateDatetime() time.Time {
    min := time.Date(2010, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
    max := time.Date(2030, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
    delta := max - min

    sec := rand.Int63n(delta) + min
    //fmt.Printf("This is the testing    %d  \n", sec)
    return time.Unix(sec, 0)
}

// Generate random decimal
func generateDecimal(precision int, scale int)  float64{
    power10 := math.Pow10(scale)
    //fmt.Printf("The precisiin is %d, the scale is %d", precision, scale)

    var num int
    if ( precision - scale > 10 ){
        num = 10
    } else {
        num =  precision - scale
    }
    randInt := rand.Int63n(int64(math.Pow10(num))) - 1
    return (math.Round(rand.Float64()*power10)/power10) + float64(randInt)

    //return fmt.Sprintf("%d.%d", randInt,  int(math.Round(rand.Float64()*power10)))
}
