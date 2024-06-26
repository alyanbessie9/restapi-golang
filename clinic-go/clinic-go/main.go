package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

var db *sql.DB

func main() {
	// Database connection
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/clinic_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Clinic API!")
	})

	// Users CRUD
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.POST("/users", createUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// PatientAppointments CRUD
	e.GET("/appointments", getAppointments)
	e.GET("/appointments/:id", getAppointment)
	e.POST("/appointments", createAppointment)
	e.PUT("/appointments/:id", updateAppointment)
	e.DELETE("/appointments/:id", deleteAppointment)

	// Drugs CRUD
	e.GET("/drugs", getDrugs)
	e.GET("/drugs/:id", getDrug)
	e.POST("/drugs", createDrug)
	e.PUT("/drugs/:id", updateDrug)
	e.DELETE("/drugs/:id", deleteDrug)

	// Patients CRUD
	e.POST("/login", login)
	e.GET("/patients", getPatients)
	e.GET("/patients/:id", getPatient)
	e.POST("/patients", createPatient)
	e.PUT("/patients/:id", updatePatient)
	e.DELETE("/patients/:id", deletePatient)

	// Doctors CRUD
	e.GET("/doctors", getDoctors)
	e.GET("/doctors/:id", getDoctor)
	e.POST("/doctors", createDoctor)
	e.PUT("/doctors/:id", updateDoctor)
	e.DELETE("/doctors/:id", deleteDoctor)

	// Transactions CRUD
	e.GET("/transactions", getAllTransactions)
	e.GET("/transactions/:id", getTransactionByID)
	e.POST("/transactions", createTransaction)
	e.PUT("/transactions/:id", updateTransaction)
	e.DELETE("/transactions/:id", deleteTransaction)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// User struct represents a user in the system
type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PatientAppointment struct represents an appointment made by a patient
type PatientAppointment struct {
	ID              uint   `json:"id"`
	PatientID       uint   `json:"patient_id"`
	UserID          int    `json:"user_id"`
	AppointmentDate string `json:"appointment_date"`
	Notes           string `json:"notes"`
	Prescription    string `json:"prescription"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// Drug struct represents a drug in the clinic
type Drug struct {
	ID                uint    `json:"id"`
	DrugName          string  `json:"drug_name"`
	DrugType          string  `json:"drug_type"`
	Description       string  `json:"description,omitempty"`
	Composition       string  `json:"composition,omitempty"`
	Packaging         string  `json:"packaging,omitempty"`
	Dosage            string  `json:"dosage,omitempty"`
	Contraindications string  `json:"contraindications,omitempty"`
	SideEffects       string  `json:"side_effects,omitempty"`
	Price             float64 `json:"price"`    // New field for price
	Currency          string  `json:"currency"` // New field for currency
	ExpirationDate    string  `json:"expiration_date,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// Patient struct represents a patient in the clinic
type Patient struct {
	ID          uint   `json:"id"`
	Nik         string `json:"nik"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Doctor represents a doctor entity
type Doctor struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Specialization   string `json:"specialization"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	ProfilePhotoPath string `json:"profile_photo_path"`
}

// Transaction represents a doctor entity
type Transaction struct {
	ID           uint    `json:"id"`
	PatientID    uint    `json:"patient_id"`
	DrugID       uint    `json:"drug_id"`
	Quantity     float64 `json:"quantity"`
	TotalPrice   float64 `json:"total_price"`
	Currency     string  `json:"currency"`
	Prescription string  `json:"prescription"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// Handler function to get all users
func getUsers(c echo.Context) error {
	rows, err := db.Query("SELECT id, name, email, created_at, updated_at FROM users")
	if err != nil {
		log.Println("Error querying users:", err)
		return c.String(http.StatusInternalServerError, "Failed to get users")
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("Error scanning user row:", err)
			continue
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

// Handler function to get a specific user by ID
func getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	var user User
	err = db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "User not found")
		}
		log.Println("Error getting user:", err)
		return c.String(http.StatusInternalServerError, "Failed to get user")
	}

	return c.JSON(http.StatusOK, user)
}

// Handler function to create a new user
func createUser(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		log.Println("Error inserting user:", err)
		return c.String(http.StatusInternalServerError, "Failed to insert user")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return c.String(http.StatusInternalServerError, "Failed to get last insert ID")
	}

	user.ID = uint(id)
	return c.JSON(http.StatusCreated, user)
}

// Handler function to update an existing user
func updateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	var user User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
	if err != nil {
		log.Println("Error updating user:", err)
		return c.String(http.StatusInternalServerError, "Failed to update user")
	}

	user.ID = uint(id)
	return c.JSON(http.StatusOK, user)
}

// Handler function to delete a user by ID
func deleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting user:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.String(http.StatusOK, fmt.Sprintf("User with ID %d deleted", id))
}

// Handler function to get all appointments
func getAppointments(c echo.Context) error {
	rows, err := db.Query("SELECT id, patient_id, user_id, appointment_date, notes, prescription, status, created_at, updated_at FROM patient_appointments")
	if err != nil {
		log.Println("Error querying appointments:", err)
		return c.String(http.StatusInternalServerError, "Failed to get appointments")
	}
	defer rows.Close()

	appointments := make([]PatientAppointment, 0)
	for rows.Next() {
		var appointment PatientAppointment
		err := rows.Scan(&appointment.ID, &appointment.PatientID, &appointment.UserID, &appointment.AppointmentDate,
			&appointment.Notes, &appointment.Prescription, &appointment.Status, &appointment.CreatedAt, &appointment.UpdatedAt)
		if err != nil {
			log.Println("Error scanning appointment row:", err)
			continue
		}
		appointments = append(appointments, appointment)
	}

	return c.JSON(http.StatusOK, appointments)
}

// Handler function to get a specific appointment by ID
func getAppointment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid appointment ID")
	}

	var appointment PatientAppointment
	err = db.QueryRow("SELECT id, patient_id, user_id, appointment_date, notes, prescription, status, created_at, updated_at FROM patient_appointments WHERE id = ?", id).Scan(
		&appointment.ID, &appointment.PatientID, &appointment.UserID, &appointment.AppointmentDate,
		&appointment.Notes, &appointment.Prescription, &appointment.Status, &appointment.CreatedAt, &appointment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Appointment not found")
		}
		log.Println("Error getting appointment:", err)
		return c.String(http.StatusInternalServerError, "Failed to get appointment")
	}

	return c.JSON(http.StatusOK, appointment)
}

// Handler function to create a new appointment
func createAppointment(c echo.Context) error {
	var appointment PatientAppointment
	if err := c.Bind(&appointment); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	// Execute the SQL query to insert a new appointment
	result, err := db.Exec("INSERT INTO patient_appointments (patient_id, user_id, appointment_date, notes, prescription, status) VALUES (?, ?, ?, ?, ?, ?)",
		appointment.PatientID, appointment.UserID, appointment.AppointmentDate,
		appointment.Notes, appointment.Prescription, appointment.Status)
	if err != nil {
		log.Println("Error inserting appointment:", err)
		return c.String(http.StatusInternalServerError, "Failed to insert appointment")
	}

	// Retrieve the last insert ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return c.String(http.StatusInternalServerError, "Failed to get last insert ID")
	}

	// Set the ID of the appointment struct and return as JSON response
	appointment.ID = uint(id)
	return c.JSON(http.StatusCreated, appointment)
}

// Handler function to update an existing appointment
func updateAppointment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid appointment ID")
	}

	var appointment PatientAppointment
	if err := c.Bind(&appointment); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	_, err = db.Exec("UPDATE patient_appointments SET patient_id = ?, doctor_id = ?, appointment_date = ?, notes = ?, prescription = ?, status = ? WHERE id = ?",
		appointment.PatientID, appointment.UserID, appointment.AppointmentDate,
		appointment.Notes, appointment.Prescription, appointment.Status, id)
	if err != nil {
		log.Println("Error updating appointment:", err)
		return c.String(http.StatusInternalServerError, "Failed to update appointment")
	}

	appointment.ID = uint(id)
	return c.JSON(http.StatusOK, appointment)
}

// Handler function to delete an appointment by ID
func deleteAppointment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid appointment ID")
	}

	_, err = db.Exec("DELETE FROM patient_appointments WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting appointment:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete appointment")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Appointment with ID %d deleted", id))
}

// Handler function to get all drugs
func getDrugs(c echo.Context) error {
	rows, err := db.Query("SELECT id, drug_name, drug_type, description, composition, packaging, dosage, contraindications, side_effects, price, currency, expiration_date, created_at, updated_at FROM drugs")
	if err != nil {
		log.Println("Error querying drugs:", err)
		return c.String(http.StatusInternalServerError, "Failed to get drugs")
	}
	defer rows.Close()

	drugs := make([]Drug, 0)
	for rows.Next() {
		var drug Drug
		err := rows.Scan(&drug.ID, &drug.DrugName, &drug.DrugType, &drug.Description, &drug.Composition, &drug.Packaging, &drug.Dosage, &drug.Contraindications, &drug.SideEffects, &drug.Price, &drug.Currency, &drug.ExpirationDate, &drug.CreatedAt, &drug.UpdatedAt)
		if err != nil {
			log.Println("Error scanning drug row:", err)
			continue
		}
		drugs = append(drugs, drug)
	}

	return c.JSON(http.StatusOK, drugs)
}

// Handler function to get a specific drug by ID
func getDrug(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid drug ID")
	}

	var drug Drug
	err = db.QueryRow("SELECT id, drug_name, drug_type, description, composition, packaging, dosage, contraindications, side_effects, price, currency, expiration_date, created_at, updated_at FROM drugs WHERE id = ?", id).Scan(
		&drug.ID, &drug.DrugName, &drug.DrugType, &drug.Description, &drug.Composition, &drug.Packaging, &drug.Dosage, &drug.Contraindications, &drug.SideEffects, &drug.Price, &drug.Currency, &drug.ExpirationDate, &drug.CreatedAt, &drug.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Drug not found")
		}
		log.Println("Error getting drug:", err)
		return c.String(http.StatusInternalServerError, "Failed to get drug")
	}

	return c.JSON(http.StatusOK, drug)
}

// Handler function to create a new drug
func createDrug(c echo.Context) error {
	var drug Drug
	if err := c.Bind(&drug); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	result, err := db.Exec("INSERT INTO drugs (drug_name, drug_type, description, composition, packaging, dosage, contraindications, side_effects, price, currency, expiration_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		drug.DrugName, drug.DrugType, drug.Description, drug.Composition, drug.Packaging, drug.Dosage, drug.Contraindications, drug.SideEffects, drug.Price, drug.Currency, drug.ExpirationDate)
	if err != nil {
		log.Println("Error inserting drug:", err)
		return c.String(http.StatusInternalServerError, "Failed to insert drug")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return c.String(http.StatusInternalServerError, "Failed to get last insert ID")
	}

	drug.ID = uint(id)
	return c.JSON(http.StatusCreated, drug)
}

// Handler function to update an existing drug
func updateDrug(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid drug ID")
	}

	var drug Drug
	if err := c.Bind(&drug); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	_, err = db.Exec("UPDATE drugs SET drug_name = ?, drug_type = ?, description = ?, composition = ?, packaging = ?, dosage = ?, contraindications = ?, side_effects = ?, price = ?, currency = ?, expiration_date = ? WHERE id = ?",
		drug.DrugName, drug.DrugType, drug.Description, drug.Composition, drug.Packaging, drug.Dosage, drug.Contraindications, drug.SideEffects, drug.Price, drug.Currency, drug.ExpirationDate, id)
	if err != nil {
		log.Println("Error updating drug:", err)
		return c.String(http.StatusInternalServerError, "Failed to update drug")
	}

	drug.ID = uint(id)
	return c.JSON(http.StatusOK, drug)
}

// Handler function to delete a drug by ID
func deleteDrug(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid drug ID")
	}

	_, err = db.Exec("DELETE FROM drugs WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting drug:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete drug")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Drug with ID %d deleted", id))
}

// Handler function to login
func login(c echo.Context) error {
	var credentials struct {
		Nik      string `json:"nik"`
		Password string `json:"password"`
	}

	if err := c.Bind(&credentials); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	var patient Patient
	err := db.QueryRow("SELECT id, nik, name, gender, date_of_birth, address, password, created_at, updated_at FROM patients WHERE nik = ?", credentials.Nik).Scan(
		&patient.ID, &patient.Nik, &patient.Name, &patient.Gender, &patient.DateOfBirth,
		&patient.Address, &patient.Password, &patient.CreatedAt, &patient.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusUnauthorized, "Invalid NIK or password")
		}
		log.Println("Error getting patient:", err)
		return c.String(http.StatusInternalServerError, "Failed to get patient")
	}

	// Check if the password matches
	if patient.Password != credentials.Password {
		return c.String(http.StatusUnauthorized, "Invalid NIK or password")
	}

	return c.JSON(http.StatusOK, patient)
}

// Handler function to get all patients
func getPatients(c echo.Context) error {
	rows, err := db.Query("SELECT id, nik, name, gender, date_of_birth, address, password, created_at, updated_at FROM patients")
	if err != nil {
		log.Println("Error querying patients:", err)
		return c.String(http.StatusInternalServerError, "Failed to get patients")
	}
	defer rows.Close()

	patients := make([]Patient, 0)
	for rows.Next() {
		var patient Patient
		err := rows.Scan(&patient.ID, &patient.Nik, &patient.Name, &patient.Gender, &patient.DateOfBirth,
			&patient.Address, &patient.Password, &patient.CreatedAt, &patient.UpdatedAt)
		if err != nil {
			log.Println("Error scanning patient row:", err)
			continue
		}
		patients = append(patients, patient)
	}

	return c.JSON(http.StatusOK, patients)
}

// Handler function to get a specific patient by ID
func getPatient(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var patient Patient
	err = db.QueryRow("SELECT id, nik, name, gender, date_of_birth, address, password, created_at, updated_at FROM patients WHERE id = ?", id).Scan(
		&patient.ID, &patient.Nik, &patient.Name, &patient.Gender, &patient.DateOfBirth,
		&patient.Address, &patient.Password, &patient.CreatedAt, &patient.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Patient not found")
		}
		log.Println("Error getting patient:", err)
		return c.String(http.StatusInternalServerError, "Failed to get patient")
	}

	return c.JSON(http.StatusOK, patient)
}

// Handler function to create a new patient
func createPatient(c echo.Context) error {
	var patient Patient
	if err := c.Bind(&patient); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	result, err := db.Exec("INSERT INTO patients (nik, name, gender, date_of_birth, address, password) VALUES (?, ?, ?, ?, ?, ?)",
		patient.Nik, patient.Name, patient.Gender, patient.DateOfBirth, patient.Address, patient.Password)
	if err != nil {
		log.Println("Error inserting patient:", err)
		return c.String(http.StatusInternalServerError, "Failed to insert patient")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return c.String(http.StatusInternalServerError, "Failed to get last insert ID")
	}

	patient.ID = uint(id)
	return c.JSON(http.StatusCreated, patient)
}

// Handler function to update an existing patient
func updatePatient(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var patient Patient
	if err := c.Bind(&patient); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	_, err = db.Exec("UPDATE patients SET nik = ?, name = ?, gender = ?, date_of_birth = ?, address = ?, password = ? WHERE id = ?",
		patient.Nik, patient.Name, patient.Gender, patient.DateOfBirth, patient.Address, patient.Password, id)
	if err != nil {
		log.Println("Error updating patient:", err)
		return c.String(http.StatusInternalServerError, "Failed to update patient")
	}

	patient.ID = uint(id)
	return c.JSON(http.StatusOK, patient)
}

// Handler function to delete a patient by ID
func deletePatient(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	_, err = db.Exec("DELETE FROM patients WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting patient:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete patient")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Patient with ID %d deleted", id))
}

// Handler function to get all doctors
func getDoctors(c echo.Context) error {
	rows, err := db.Query("SELECT id, user_id, specialization, created_at, updated_at, profile_photo_path FROM doctors")
	if err != nil {
		log.Println("Error querying doctors:", err)
		return c.String(http.StatusInternalServerError, "Failed to get doctors")
	}
	defer rows.Close()

	doctors := make([]Doctor, 0)
	for rows.Next() {
		var doctor Doctor
		err := rows.Scan(&doctor.ID, &doctor.UserID, &doctor.Specialization, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.ProfilePhotoPath)
		if err != nil {
			log.Println("Error scanning doctor row:", err)
			continue
		}
		doctors = append(doctors, doctor)
	}

	// Print doctors slice for debugging
	fmt.Println("Doctors:", doctors)

	return c.JSON(http.StatusOK, doctors)
}

// Handler function to get a specific doctor by ID
func getDoctor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid doctor ID")
	}

	var doctor Doctor
	err = db.QueryRow("SELECT id, user_id, specialization, created_at, updated_at, profile_photo_path FROM doctors WHERE id = ?", id).
		Scan(&doctor.ID, &doctor.UserID, &doctor.Specialization, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.ProfilePhotoPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Doctor not found")
		}
		log.Println("Error getting doctor:", err)
		return c.String(http.StatusInternalServerError, "Failed to get doctor")
	}

	return c.JSON(http.StatusOK, doctor)
}

// Handler function to create a new doctor
func createDoctor(c echo.Context) error {
	var doctor Doctor
	if err := c.Bind(&doctor); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	result, err := db.Exec("INSERT INTO doctors (user_id, specialization, created_at, updated_at, profile_photo_path) VALUES (?, ?, ?, ?, ?)",
		doctor.UserID, doctor.Specialization, time.Now(), time.Now(), doctor.ProfilePhotoPath)
	if err != nil {
		log.Println("Error inserting doctor:", err)
		return c.String(http.StatusInternalServerError, "Failed to insert doctor")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return c.String(http.StatusInternalServerError, "Failed to get last insert ID")
	}

	doctor.ID = int(id)
	return c.JSON(http.StatusCreated, doctor)
}

// Handler function to update an existing doctor
func updateDoctor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid doctor ID")
	}

	var doctor Doctor
	if err := c.Bind(&doctor); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	_, err = db.Exec("UPDATE doctors SET user_id = ?, specialization = ?, profile_photo_path = ?, updated_at = ? WHERE id = ?",
		doctor.UserID, doctor.Specialization, doctor.ProfilePhotoPath, time.Now(), id)
	if err != nil {
		log.Println("Error updating doctor:", err)
		return c.String(http.StatusInternalServerError, "Failed to update doctor")
	}

	doctor.ID = id
	return c.JSON(http.StatusOK, doctor)
}

// Handler function to delete a doctor by ID
func deleteDoctor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid doctor ID")
	}

	_, err = db.Exec("DELETE FROM doctors WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting doctor:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete doctor")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Doctor with ID %d deleted", id))
}

func getAllTransactions(c echo.Context) error {
	rows, err := db.Query("SELECT id, patient_id, drug_id, quantity, total_price, currency, prescription, created_at, updated_at FROM transactions")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get transactions")
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.PatientID, &t.DrugID, &t.Quantity, &t.TotalPrice, &t.Currency, &t.Prescription, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to scan transactions")
		}
		transactions = append(transactions, t)
	}

	return c.JSON(http.StatusOK, transactions)
}

func getTransactionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid transaction ID")
	}

	var t Transaction
	err = db.QueryRow("SELECT id, patient_id, drug_id, quantity, total_price, currency, prescription, created_at, updated_at FROM transactions WHERE id = ?", id).Scan(
		&t.ID, &t.PatientID, &t.DrugID, &t.Quantity, &t.TotalPrice, &t.Currency, &t.Prescription, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Transaction not found")
		}
		return c.String(http.StatusInternalServerError, "Failed to get transaction")
	}

	return c.JSON(http.StatusOK, t)
}

func createTransaction(c echo.Context) error {
	var t Transaction
	if err := c.Bind(&t); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	result, err := db.Exec("INSERT INTO transactions (patient_id, drug_id, quantity, total_price, currency, prescription) VALUES (?, ?, ?, ?, ?, ?)",
		t.PatientID, t.DrugID, t.Quantity, t.TotalPrice, t.Currency, t.Prescription)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to insert transaction")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get last insert ID")
	}
	t.ID = uint(id)

	return c.JSON(http.StatusCreated, t)
}

func updateTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid transaction ID")
	}

	var t Transaction
	if err := c.Bind(&t); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	_, err = db.Exec("UPDATE transactions SET patient_id = ?, drug_id = ?, quantity = ?, total_price = ?, currency = ?, prescription = ? WHERE id = ?",
		t.PatientID, t.DrugID, t.Quantity, t.TotalPrice, t.Currency, t.Prescription, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update transaction")
	}

	t.ID = uint(id)
	return c.JSON(http.StatusOK, t)
}

func deleteTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid transaction ID")
	}

	_, err = db.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete transaction")
	}

	return c.NoContent(http.StatusNoContent)
}
