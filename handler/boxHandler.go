package handler

import (
	"backend/db"
	"backend/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func BoxDepositHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// ✅ Handle preflight OPTIONS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodPost:
		CreateBoxRequest(w, r)
	case http.MethodGet:
		GetBoxDeposit(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateBoxRequest(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	boxDeposit := db.Database.Collection("boxDeposit")
	var req models.BoxDeposit
	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := boxDeposit.InsertOne(ctx, bson.M{
		"firstName":      req.FirstName,
		"lastName":       req.LastName,
		"email":          req.Email,
		"phone":          req.Phone,
		"contactMethod":  req.ContactMethod,
		"boxSize":        req.BoxSize,
		"duration":       req.Duration,
		"referral":       req.Referral,
		"additionalInfo": req.AdditionalInfo,
		"createdAt":      time.Now(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Box request created successfully",
		"id":      result.InsertedID,
	})
}

func GetBoxDeposit(w http.ResponseWriter, r *http.Request) {
	collection := db.Database.Collection("boxDeposit")

	// ✅ Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to query database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []models.BoxDeposit
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, "Failed to decode results: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Return JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func DownloadBoxDepositExcel(w http.ResponseWriter, r *http.Request) {
	collection := db.Database.Collection("boxDeposit")

	// ✅ Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to query database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var deposits []models.BoxDeposit
	if err := cursor.All(ctx, &deposits); err != nil {
		http.Error(w, "Failed to decode results: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Buat file Excel
	f := excelize.NewFile()
	sheet := "BoxDeposit"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1") // Hapus sheet default

	// Header Excel
	headers := []string{"First Name", "Last Name", "Email", "Phone", "Contact Method", "Box Size", "Duration", "Referral", "Additional Info"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data Excel
	for rowIdx, d := range deposits {
		values := []interface{}{d.FirstName, d.LastName, d.Email, d.Phone, d.ContactMethod, d.BoxSize, d.Duration, d.Referral, d.AdditionalInfo}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	// Response headers untuk download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="box_deposit.xlsx"`)

	// Kirim file ke client
	if err := f.Write(w); err != nil {
		http.Error(w, "Failed to generate Excel file: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
