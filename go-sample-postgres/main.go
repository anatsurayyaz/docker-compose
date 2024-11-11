package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Connect to database
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Create Database Connection failed: ", err)
		return
	}

	// https://stackoverflow.com/a/32345308
	// Real connect to db
	fmt.Println("Connect to: ", connStr)
	if err := db.Ping(); err != nil {
		fmt.Printf("Connection to database failed (DB_HOST: %s): %s\n", dbHost, err)
	} else {
		fmt.Println("Successfully connected to the database: ", db)
	}
	defer db.Close() // Ensure the connection is closed when the program exits

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route untuk halaman utama
	e.GET("/", func(c echo.Context) error {
		fmt.Println("Someone hit me!")
		return c.HTML(http.StatusOK, "Do your best and let the results speak for themselves")
	})

	// Route baru untuk /connect
	e.GET("/connect", func(c echo.Context) error {
		// Coba koneksi ke database
		if err := db.Ping(); err != nil {
			// Jika koneksi gagal
			fmt.Println("Failed connect to db: ", err)
			return c.HTML(http.StatusInternalServerError, "Failed connect to db")
		}

		// Insert log akses ke dalam tabel {yourusername}_access_log
		accessLogTable := "anatsurayyaz_access_log" // Ganti dengan nama tabel sesuai username Anda
		insertQuery := fmt.Sprintf("INSERT INTO %s (timestamp) VALUES (NOW())", accessLogTable)
		_, err := db.Exec(insertQuery)
		if err != nil {
			fmt.Println("Failed to insert log: ", err)
			return c.HTML(http.StatusInternalServerError, "Failed to insert log into db")
		}

		// Jika sukses, tampilkan pesan
		fmt.Println("Success connect to db")
		return c.HTML(http.StatusOK, "Success connect to db")
	})

	// Set port untuk menjalankan aplikasi
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "80"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
