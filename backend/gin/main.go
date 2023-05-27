package main

import (
	d "backend/database/dbInfo"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// dALocation would be used to record the current CSV file's data attribute location.
var dALocation = map[string]int{
	"Started Time":        0,
	"Duration":            0,
	"Suite":               0,
	"Board":               0,
	"Model":               0,
	"Build Version":       0,
	"Hostname":            0,
	"Test":                0,
	"Status":              0,
	"Failure Reason":      0,
	"Firmware RO Version": 0,
	"Firmware RW Version": 0,
}

var DB *sql.DB
var Router *gin.Engine

const csvFolder = "stainless"

func init() {
	// Init Router
	Router = gin.Default()

	// Set up DB connection
	out, _ := ioutil.ReadFile("dbConfig.json")
	var dbLogin d.DbConfig
	if err := json.Unmarshal(out, &dbLogin); err != nil {
		log.Fatalf("Failed to unmarshal: %s", out)
	}

	sqlDB := d.NewSql(dbLogin)
	db, err := sqlDB.NewSqlDBConnect()
	if err != nil {
		log.Fatal("Failed to new sql DB")
	}
	DB = db
}

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	Router.Use(cors.New(corsConfig))

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	d.Log("Database is connected")
	defer DB.Close()

	// Create if necessary
	createTables := func() {
		t := d.TableJson{
			JsonFile: "dbConfig_table.json",
		}
		mapping := t.GenerateTableMap()
		for tableName, sqlString := range mapping {
			if err := d.CreateTable(DB, string(tableName), sqlString); err != nil {
				log.Fatal("Failed to create tables: ", err)
			}
		}
	}
	createTables()
	d.Log("All needed tables were created")

	Router.GET("/localTest", handleLocalTest)
	Router.GET("/stainlessSearch", handleStainlessTest)
	Router.POST("/uploadCSV", handleUploadCsv)
	Router.GET("/passingRate", handlePassingRate)
	Router.GET("/getTicketList", handelTicketList)
	Router.GET("/ownerPassingRate", handleOwnerCasePassingRate)

	Router.Static("/logDB", "../logDB")
	Router.Run(":8082")
}

func handleLocalTest(ctx *gin.Context) {
	params := d.SearchInput{
		TableName: "Result",
		Board:     ctx.Query("board"),
		Reason:    ctx.Query("reason"),
		Name:      ctx.Query("testName"),
		Result:    ctx.Query("result"),
		StartDate: ctx.Query("startDate"),
		EndDate:   ctx.Query("endDate"),
		OrderBy:   ctx.Query("orderBy"),
	}

	dataResult, err := d.SearchOnLocalTest(DB, params)
	if err != nil {
		log.Panic("Failed to search data in table `Result`: ", err)
	}
	ctx.IndentedJSON(http.StatusOK, dataResult)
}

func handleStainlessTest(ctx *gin.Context) {
	params := d.SearchInput{
		TableName: "Stainless_Result",
		Board:     ctx.Query("board"),
		Reason:    ctx.Query("reason"),
		Name:      ctx.Query("testName"),
		Result:    ctx.Query("result"),
		StartDate: ctx.Query("startDate"),
		EndDate:   ctx.Query("endDate"),
		OrderBy:   ctx.Query("orderBy"),
	}

	stainlessResult, err := d.SearchStainlessResult(DB, params)
	if err != nil {
		log.Panicf("Failed to search data in table %q: %v", params.TableName, err)
	}
	ctx.IndentedJSON(http.StatusOK, stainlessResult)
}

func handleUploadCsv(ctx *gin.Context) {
	file, _ := ctx.FormFile("stainlessData")
	d.Log(file.Filename)
	bytes, err := getCsvFiles()
	if err != nil {
		log.Fatalf("Failed to get csv files in %v", csvFolder)
	}
	if err := d.ValidCsv(file.Filename, bytes); err != nil {
		ctx.String(http.StatusOK, "Error happend: ", err)
	} else {
		ctx.SaveUploadedFile(file, "stainless/"+file.Filename)
		d.Log("Saving files in server sucessfully. Start inserting data to DB")
		if err := runGoToSaveStainlessData("stainless/"+file.Filename, DB); err != nil {
			// Todo: error happens deletes the csv file below
			ctx.String(http.StatusOK, "Error happend while inserting the data into DB. Data is ready to rollback: ", err)
		} else {
			ctx.String(http.StatusOK, "file: %s", file.Filename)
		}
	}
}

func handlePassingRate(ctx *gin.Context) {
	testName := ctx.Query("testName")
	passingRate, err := d.CatchPassAndFail(DB, "Stainless_Result", testName)
	if err != nil {
		log.Panic("Failed to catch passing rate from backend: ", err)
	}
	ctx.IndentedJSON(http.StatusOK, passingRate)
}

func handelTicketList(ctx *gin.Context) {
	ticketList, err := d.SearchTicketList(DB)
	if err != nil {
		log.Panic("Failed to catch passing rate from backend: ", err)
	}
	ctx.IndentedJSON(http.StatusOK, ticketList)
}

func handleOwnerCasePassingRate(ctx *gin.Context) {
	owner := ctx.Query("owner")
	caseList, err := d.OnwerCaseMaintain(DB, owner)
	if err != nil {
		log.Panic("Failed to catch owner case list from DB: ", err)
	}
	ctx.IndentedJSON(http.StatusOK, caseList)
}

func runGoToSaveStainlessData(fileName string, db *sql.DB) error {
	rollback := true
	var maxID int

	defer func(rollback *bool, maxID *int) {
		if *rollback {
			if err := removeCSV(fileName); err != nil {
				log.Panic("Failed to remove csv file: ", err)
			}
			log.Println("Rollback: ", *rollback, " MaxID: ", *maxID)
			queryDeleteCmd := fmt.Sprintf("DELETE FROM %s WHERE id > ?;", "Stainless_Result")
			if err := d.RunSqlStmt(db, queryDeleteCmd, *maxID); err != nil {
				log.Panic("Failed to rollback data: ", err)
			}
		}
	}(&rollback, &maxID)

	errMaxId := db.QueryRow("SELECT Max(id) FROM Stainless_Result").Scan(&maxID)
	if errMaxId != nil && !strings.Contains(errMaxId.Error(), "Scan error on column index 0") {
		log.Panic("Failed to find max id in Stainless_Result: ", errMaxId)
	}
	d.Logf("Max id: %d", maxID)

	rows := readCSV(fileName)
	d.Log("Start to check collums in csv")
	if len(rows[0]) != 12 {
		log.Panicf(fmt.Sprintf("The numbers of collumns should be 12, but it has %d", len(rows)))
	}
	for attributeString, _ := range dALocation {
		for i := 0; i < 12; i++ {
			if strings.Contains(rows[0][i], attributeString) {
				dALocation[attributeString] = i
				break
			}
			if i == 11 {
				log.Panicf(fmt.Sprintf("Attribute [%v] wasn't found on csv files", attributeString))
			}
		}
	}
	rows = rows[1:]

	batchSize := 12000
	batchCount := 0
	dataInsert := ""
	batchData := []string{}
	for i, data := range rows {
		rowData := d.StainlessResult{
			Time:              data[dALocation["Started Time"]],
			Duration:          data[dALocation["Duration"]],
			Suite:             data[dALocation["Suite"]],
			Board:             data[dALocation["Board"]],
			Model:             data[dALocation["Model"]],
			BuildVersion:      data[dALocation["Build Version"]],
			Host:              data[dALocation["Hostname"]],
			TestName:          data[dALocation["Test"]][14:],
			Status:            data[dALocation["Status"]],
			Reason:            data[dALocation["Failure Reason"]],
			FirmwareROVersion: data[dALocation["Firmware RO Version"]],
			FirmwareRWVersion: data[dALocation["Firmware RW Version"]],
		}

		if batchCount >= batchSize || i == len(rows)-1 {
			d.Log("Batch insert")
			value := strings.Join(batchData, ",")
			preSQL := "insert into Stainless_Result (time,duration,suite,board,model,buildVersion,host,testName,status,reason,firmwareROVersion,firmwareRWVersion) values " + value + ";"
			errInsert := d.RunSqlStmt(db, preSQL)
			if errInsert != nil {
				log.Panicf(fmt.Sprintf("Failed to insert [ID: %v] %s into DB: %v", i, data, errInsert))
			}
			batchData, batchCount = nil, 0
		} else {
			dataInsert = fmt.Sprintf("(%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q)", rowData.Time, rowData.Duration, rowData.Suite, rowData.Board, rowData.Model, rowData.BuildVersion, rowData.Host, rowData.TestName, rowData.Status, rowData.Reason, rowData.FirmwareROVersion, rowData.FirmwareRWVersion)
			batchData = append(batchData, dataInsert)
		}
		batchCount++
	}
	rollback = false
	d.Log("Insert Successfully")
	return nil
}

// readCSV would read a csv file and put it in a slice
func readCSV(csvFile string) [][]string {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Panicf("Failed to open %q: %v", csvFile, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true

	rows, err := r.ReadAll()
	if err != nil {
		log.Panic("Cannot read CSV data:", err)
	}
	return rows
}

func removeCSV(csvFile string) error {
	cmdString := fmt.Sprintf("rm %s", csvFile)
	cmd := exec.Command("bash", "-c", cmdString)

	d.Log("Removing csv")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("while running %s: %v", cmdString, err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("an error happened when remove csv file in server: %v", err)
	}
	return nil
}

func getCsvFiles() ([]byte, error) {
	cmdString := fmt.Sprintf("ls %s | grep .csv", csvFolder)
	cmd := exec.Command("bash", "-c", cmdString)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	d.Log("Searching the logs in Server...")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("while running ls command: %v", err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("while reading the command output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		// "exit status 1" means "not found" here.
		if !strings.Contains(err.Error(), "exit status 1") {
			return nil, fmt.Errorf("unexpected error happened while running ls command: %v", err)
		}
	}
	return bytes, nil
}
