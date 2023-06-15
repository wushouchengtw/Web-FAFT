package handlers

import (
	"backend/lib"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUploadCsv(ctx *gin.Context) {
	// To-do: change the name
	file, _ := ctx.FormFile("stainlessData")
	if err := lib.ValidCsv(file.Filename); err != nil {
		ctx.String(http.StatusOK, "Error happend: ", err)
	} else {
		ctx.SaveUploadedFile(file, lib.CsvFolder+file.Filename)
		log.Println("Saving files in server sucessfully. Start inserting data to DB")
		// To-do: Save data into db
		// if err := runGoToSaveStainlessData("stainless/"+file.Filename, DB); err != nil {
		// 	// Todo: error happens deletes the csv file below
		// 	ctx.String(http.StatusOK, "Error happend while inserting the data into DB. Data is ready to rollback: ", err)
		// } else {
		// 	ctx.String(http.StatusOK, "file: %s", file.Filename)
		// }
	}
}
