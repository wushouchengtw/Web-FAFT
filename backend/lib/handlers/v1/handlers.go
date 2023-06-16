package v1

import (
	"backend/lib"
	remoterepo "backend/lib/remoteRepo"
	"backend/lib/remoteRepo/dut"
	"backend/lib/remoteRepo/result"
	"backend/lib/remoteRepo/test"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UploadCsv(db *sqlx.DB) gin.HandlerFunc {
	// To-do: change the name
	return func(ctx *gin.Context) {
		file, _ := ctx.FormFile("stainlessData")
		if err := lib.ValidCsv(file.Filename); err != nil {
			ctx.String(http.StatusOK, "Error happend: ", err)
		} else {
			ctx.SaveUploadedFile(file, lib.CsvFolder+file.Filename)
			log.Println("Saving files in server sucessfully. Start inserting data to DB")

			dutRepo := dut.NewDUTRepoInMySQL(db)
			testRepo := test.NewTestRepoInMySQL(db)
			resultRepo := result.NewResultRepoMySQL(db)

			if err := remoterepo.SaveRemoteDataByCsv(lib.CsvFolder+"/"+file.Filename, dutRepo, testRepo, resultRepo); err != nil {
				if err := lib.RemoveCsvFile(file.Filename); err != nil {
					ctx.String(http.StatusOK, "Failed to remove the csv file: ", err)
				}
				ctx.String(http.StatusOK, "Error happend while inserting the data into DB. Data is ready to rollback: ", err)
			} else {
				ctx.String(http.StatusOK, "file: %s", file.Filename)
			}
		}
	}
}
