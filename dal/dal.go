package dal

import (
	"ServerGin/app"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// PSQL represents data for connection to Data base
type PSQL struct {
	Host     string
	DataBase *sql.DB
}

// NewPSQL constructs object of PSQL
func NewPSQL(host string, port int) (*PSQL, error) {
	var err error
	var db *sql.DB

	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		log.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, username, password, dbName)
	_ = connString
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected to db!\n")

	res := &PSQL{
		Host:     host,
		DataBase: db}

	return res, nil
}

func (t *PSQL) Read(current *app.Event) (*app.AllEvents, error) {
	events := app.GetEvents()

	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return nil, err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	i1, err := strconv.Atoi(current.Page)
	if err != nil {
		log.Fatal("page type error. ", err.Error())
	}
	i2, err := strconv.Atoi(current.Count)
	if err != nil {
		log.Fatal("count type error. ", err.Error())
	}

	conditions := getStrings(current)

	tsql := fmt.Sprintf("SELECT name, post, datestart, dateend FROM info %s limit %d offset %d;",
		conditions, i2, (i1-1)*i2)

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error")
		return nil, err
	}
	defer stmt.Close()

	rows, err := t.DataBase.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var name, post, dateStart, dateEnd string

		// Get values from row.
		err := rows.Scan(&name, &post, &dateStart, &dateEnd)
		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
			return nil, err
		}

		ev := *app.NewEvent()
		ev.Name = name
		ev.Post = post
		ev.DateStart = dateStart
		ev.DateEnd = dateEnd
		*events = append(*events, ev)
	}
	return events, nil
}

func (t *PSQL) Read2(current *app.Event) (int, error) {

	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return 0, err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	conditions := getStrings(current)

	tsql1 := fmt.Sprintf("SELECT count(*) FROM info %s;", conditions)

	log.Print("Ok1 in Reda2")

	stmt1, err := t.DataBase.Prepare(tsql1)
	if err != nil {
		log.Println("Prepare error")
		return 0, err
	}
	defer stmt1.Close()

	log.Print("Ok2 in Reda2")

	rows1, err := t.DataBase.QueryContext(ctx, tsql1)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return 0, err
	}

	defer rows1.Close()

	var count int
	for rows1.Next() {
		err = rows1.Scan(&count)
	}

	log.Print("Count: ", count)

	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return 0, err
	}

	return count, nil

}

func getStrings(current *app.Event) string {
	conditions := ""
	check_condition := false

	if current.Name != "" {
		conditions += " name = " + current.Name
		check_condition = true
	}
	if current.Post != "" {
		if check_condition {
			conditions += " and "
		}
		conditions += " post = " + current.Post
		check_condition = true
	}
	if current.DateStart != "" {
		if check_condition {
			conditions += " and "
		}
		conditions += " datestart >= '" + current.DateStart + "'::date"
		check_condition = true
	}
	if current.DateEnd != "" {
		if check_condition {
			conditions += " and "
		}
		conditions += " dateend <= '" + current.DateEnd + "'::date"
		check_condition = true
	}

	if check_condition {
		conditions = " where " + conditions + " "
	}
	return conditions
}

func (t *PSQL) GetCsv() error {
	//var err error
	//file, err := ioutil.ReadFile("get_csv.sql") //file
	var err error
	ctx := context.Background()

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	//conditions := "where name = 9" //getStrings(current)
	//tsql := fmt.Sprintf("COPY info TO '/tmp/products.csv' WITH (FORMAT CSV, HEADER);", conditions) //(select * from info %s)

	_, err = t.DataBase.Exec("COPY info TO '/tmp/products.csv' WITH (FORMAT CSV, HEADER);") //Exec(string(file))
	if err != nil {
		return err
	}
	return nil

}
