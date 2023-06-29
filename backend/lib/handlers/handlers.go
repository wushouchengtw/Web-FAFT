package handlers

import (
	"backend/lib/csv"
	"backend/lib/remoteRepo/dut"
	"backend/lib/remoteRepo/result"
	"backend/lib/remoteRepo/test"
	"backend/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func HanlderUploadCsv(db *sqlx.DB) gin.HandlerFunc {
	// To-do: change the name
	return func(ctx *gin.Context) {
		file, _ := ctx.FormFile("stainlessData")
		if err := csv.ValidCsv(file.Filename); err != nil {
			ctx.String(http.StatusOK, "Error happend: ", err)
		} else {
			ctx.SaveUploadedFile(file, csv.CsvFolder+file.Filename)
			log.Println("Saving files in server sucessfully. Start inserting data to DB")

			dutRepo := dut.NewDUTRepoInMySQL(db)
			testRepo := test.NewTestRepoInMySQL(db)
			resultRepo := result.NewResultRepoMySQL(db)

			if err := csv.SaveRemoteDataByCsv(csv.CsvFolder+"/"+file.Filename, dutRepo, testRepo, resultRepo); err != nil {
				if err := csv.RemoveCsvFile(file.Filename); err != nil {
					ctx.String(http.StatusOK, "Failed to remove the csv file: ", err)
				}
				ctx.String(http.StatusOK, "Error happend while inserting the data into DB. Data is ready to rollback: ", err)
			} else {
				ctx.String(http.StatusOK, "file: %s", file.Filename)
			}
		}
	}
}

func HanlderTesthaus(db *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime, err := time.Parse(timeLayout, ctx.Query("startDate"))
		if err != nil {
			log.Panicf("unexpected startTime [%s]", ctx.Query("startDate"))
		}
		endTime, err := time.Parse(timeLayout, ctx.Query("endDate"))
		if err != nil {
			log.Panicf("unexpected endTime [%s]", ctx.Query("endDate"))
		}

		var status bool
		switch strings.ToLower(ctx.Query("result")) {
		case "pass":
			status = true
		case "fail":
			status = false
		default:
			log.Panicf("unexpected status [%s]", ctx.Query("result"))
		}

		params := utils.QueryParameter{
			StartDate: startTime,
			EndDate:   endTime,
			Board:     ctx.Query("board"),
			Reason:    ctx.Query("reason"),
			Name:      ctx.Query("testName"),
			Status:    status,
		}
		resultRepo := result.NewResultRepoMySQL(db)
		output, err := resultRepo.SearchTestHaus(params)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, output)
		}
		ctx.IndentedJSON(http.StatusOK, output)
	}
}
