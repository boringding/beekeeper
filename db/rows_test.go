package db

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type S1 struct {
	Col1 int    `col:"true"`
	Col2 string `col:"true"`
	Col3 string `col:"true"`
	Col4 string `col:"true"`
	Col5 int64  `col:"true"`
	Col6 string `col:"true"`
}

type S2 struct {
	Col1 int64  `col:"true" id:"true"`
	Col2 string `col:"true"`
	Colx int
	Coly S1      `col:"true"`
	Col3 int     `col:"true"`
	Col4 int     `col:"true"`
	Col5 float64 `col:"true"`
}

func Test_Rows2Slice(t *testing.T) {
	dataSrcName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "", "127.0.0.1", 3306, "test_db")

	db, err := sql.Open("mysql", dataSrcName)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	rows, err := db.Query("select col1,col2,col3,col4,col5,col6 from tb_a")
	if err != nil {
		fmt.Println(err)
	}

	s1s := make([]S1, 0, 10)

	err = Rows2Slice(rows, &s1s)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s1s)
	rows.Close()

	rows, err = db.Query("select col1,col2,col3,col4,col5 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	s2s := make(map[int64]S2, 0)

	err = Rows2Map(rows, &s2s)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s2s)
	rows.Close()
}
