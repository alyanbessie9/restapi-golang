package main

import (
	"database/sql" // Package database/sql digunakan untuk berinteraksi dengan database SQL
	"flag"         // Package flag digunakan untuk membaca argumen dari baris perintah
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"         // Framework web untuk membuat REST API
	_ "github.com/go-sql-driver/mysql" // Driver MySQL untuk package database/sql
)

// Struct untuk menyimpan data person
type person struct {
	ID       string `json:"id"`
	FullName string `json:"FullName"`
	Age      int    `json:"age"`
}

var db *sql.DB // Variabel untuk menyimpan koneksi ke database

// Struct untuk menyimpan data person yang telah diurutkan
type sortedPersons []person

// Fungsi yang akan digunakan untuk mengurutkan data person berdasarkan ID
func (s sortedPersons) Len() int           { return len(s) }
func (s sortedPersons) Less(i, j int) bool { return s[i].ID < s[j].ID }
func (s sortedPersons) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// initDB menginisialisasi koneksi ke database MySQL
func initDB() error {
	var err error
	db, err = sql.Open("mysql", "alyanbessie:15210003@tcp(localhost:3306)/dbrestapi")
	if err != nil {
		return err
	}

	// Pastikan koneksi ke database berhasil
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Handler untuk mendapatkan seluruh data persons dari database, tersusun dengan rapi berdasarkan ID
func getSortedPersons(c *gin.Context) {
	rows, err := db.Query("SELECT id, FullName, age FROM persons")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var persons []person
	for rows.Next() {
		var p person
		if err := rows.Scan(&p.ID, &p.FullName, &p.Age); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		persons = append(persons, p)
	}

	// Mengurutkan data persons berdasarkan ID
	sort.Sort(sortedPersons(persons))

	c.JSON(http.StatusOK, persons)
}

// createPersons menambahkan data person baru ke dalam database
func createPersons(c *gin.Context) {
	var newPerson person
	if err := c.BindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO persons (id, FullName, age) VALUES (?, ?, ?)", newPerson.ID, newPerson.FullName, newPerson.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPerson)
}

// Handler untuk mendapatkan data person berdasarkan ID dari database
func getPersonById(c *gin.Context) {
	id := c.Param("id")

	var p person
	err := db.QueryRow("SELECT id, FullName, age FROM persons WHERE id = ?", id).Scan(&p.ID, &p.FullName, &p.Age)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// Handler untuk menghapus data person dari database berdasarkan ID
func deletePerson(c *gin.Context) {
	id := c.Param("id")

	// Menjalankan query DELETE untuk menghapus data person berdasarkan ID
	result, err := db.Exec("DELETE FROM persons WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Memeriksa apakah data person berhasil dihapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
		return
	}

	// Mengembalikan pesan sukses setelah data person berhasil dihapus
	c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}

func main() {
	// Inisialisasi koneksi ke database MySQL
	if err := initDB(); err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Membaca argumen dari baris perintah untuk menambahkan data person baru
	addID := flag.String("id", "", "ID of the person")
	addFullName := flag.String("name", "", "Full name of the person")
	addAge := flag.Int("age", 0, "Age of the person")
	flag.Parse()

	// Jika argumen untuk ID, nama lengkap, dan usia disediakan, tambahkan data person baru ke dalam database
	if *addID != "" && *addFullName != "" && *addAge > 0 {
		newPerson := person{
			ID:       *addID,
			FullName: *addFullName,
			Age:      *addAge,
		}
		_, err := db.Exec("INSERT INTO persons (id, FullName, age) VALUES (?, ?, ?)", newPerson.ID, newPerson.FullName, newPerson.Age)
		if err != nil {
			fmt.Println("Error adding person:", err)
			return
		}
	}

	// Membuat router menggunakan framework Gin
	router := gin.Default()

	// Menentukan endpoint untuk mendapatkan seluruh data persons, tersusun dengan rapi berdasarkan ID
	router.GET("/persons", getSortedPersons)
	// Menentukan endpoint untuk menambahkan data person baru
	router.POST("/persons", createPersons)
	// Menentukan endpoint untuk mendapatkan data person berdasarkan ID
	router.GET("/persons/:id", getPersonById)
	// Menentukan endpoint untuk menghapus data person berdasarkan ID
	router.DELETE("/persons/:id", deletePerson)

	// Menjalankan server HTTP pada localhost:8080
	router.Run("localhost:8080")
}
