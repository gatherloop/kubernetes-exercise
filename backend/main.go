package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Port             string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
}

func getConfig() Config {
	return Config{
		Port:             os.Getenv("PORT"),
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabaseUsername: os.Getenv("DB_USERNAME"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
	}
}

func connect(config Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DatabaseUsername, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName))
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Student struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	config := getConfig()

	db, err := connect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM students")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()

		var students []Student

		for rows.Next() {
			var student Student
			err := rows.Scan(&student.Id, &student.Name, &student.Age, &student.Address, &student.Phone)

			if err != nil {
				fmt.Println(err)
				return
			}

			students = append(students, student)
		}

		jsonData, err := json.Marshal(students)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
		w.Write(jsonData)
	})

	http.ListenAndServe(":"+config.Port, nil)
}
