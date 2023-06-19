package handlers

import (
	"backend/lib/csv"
	"backend/lib/remoteRepo/dut"
	"backend/lib/remoteRepo/result"
	"backend/lib/remoteRepo/test"
	"backend/utils/query"
	"log"
	"net/http"

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
		params := query.QueryParameter{
			Board:     ctx.Query("board"),
			Reason:    ctx.Query("reason"),
			Name:      ctx.Query("testName"),
			Result:    ctx.Query("result"),
			StartDate: ctx.Query("startDate"),
			EndDate:   ctx.Query("endDate"),
			OrderBy:   ctx.Query("orderBy"),
		}
		resultRepo := result.NewResultRepoMySQL(db)
		// To-do
		resultRepo.SearchTestHaus(params)
		// 	stainlessResult, err := resultRepo.SearchTestHaus(params)
		// 	if err != nil {
		// 		log.Panicf("Failed to search data in table %q: %v", params.TableName, err)
		// 	}
		// 	ctx.IndentedJSON(http.StatusOK, stainlessResult)
		// }
	}
}
