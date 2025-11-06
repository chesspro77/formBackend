package main

import (
	"backend/db"
	"backend/handler"
	"fmt"
	"net/http"
)

func main() {
	db.ConnectDB()
	http.HandleFunc("/api/box", handler.BoxDepositHandler)
	http.HandleFunc("/download-excel", handler.DownloadBoxDepositExcel)
	fs := http.FileServer(http.Dir("dist"))
	go http.Handle("/", fs)
	fmt.Println("server started at localhost:80")
	http.ListenAndServe("0.0.0.0:80", nil)
}
