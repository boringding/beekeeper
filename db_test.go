package beekeeper

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
	//Coly S1      `col:"true"`
	Col3 int     `col:"true"`
	Col4 int     `col:"true"`
	Col5 float64 `col:"true"`
}

type S3 struct {
	Col1 uint64 `col:"true"`
	Col2 string `col:"true"`
	Col3 bool   `col:"true"`
}

func Test_QueryRow(t *testing.T) {
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

	var a float64
	err = QueryRow(db, &a, "select col5 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("a:", a)

	var b string
	err = QueryRow(db, &b, "select col3 from tb_a where col1=?", 22)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("b:", b)

	var s1 S1
	err = QueryRow(db, &s1, "select col1,col2,col3,col4,col5,col6 from tb_a where col1=? and col3=?", 22, "2016-05-11")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("s1:", s1)

	var s2 S2
	err = QueryRow(db, &s2, "select col1,col2,col3,col4,col5 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("s2:", s2)
}

func Test_QueryRows(t *testing.T) {
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

	a := make([]float64, 0, 5)
	err = QueryRows(db, &a, "select col5 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("a:", a)

	b := make([]string, 0, 5)
	err = QueryRows(db, &b, "select col3 from tb_a where col1=?", 22)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("b:", b)

	s1 := make([]S1, 0, 5)
	err = QueryRows(db, &s1, "select col1,col2,col3,col4,col5,col6 from tb_a")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("s1:", s1)

	s2 := make([]S2, 0, 5)
	err = QueryRows(db, &s2, "select col1,col2,col3,col4,col5 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("s2:", s2)

	s3 := make(map[int64]S2, 0)
	err = QueryRows(db, &s3, "select col1,col2,col3,col4,col5 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("s3:", s3)

	s4 := make([]S3, 0, 5)
	err = QueryRows(db, &s4, "select col1,col2,col3 from tb_b")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("s4:", s4)
}
