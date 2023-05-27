package dbInfo

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/now"
)

var (
	boardMap    = "board"
	testNameMap = "testName"
	resultMap   = "result"
	reasonMap   = "reason"
)

// sql statement map
var tableSearchMap = map[string]map[string]string{
	"Stainless_Result": {
		boardMap:    "board = ?",
		testNameMap: "testName = ?",
		resultMap:   "status = ?",
		reasonMap:   "reason LIKE ?",
	},
	"Result": {
		boardMap:    "board = ?",
		testNameMap: "name = ?",
		resultMap:   "result = ?",
		reasonMap:   "reason LIKE ?",
	},
}

func SearchOnLocalTest(db *sql.DB, params SearchInput) ([]LocalResult, error) {
	var rows *sql.Rows
	var err error
	sqlSearch := "select id,time,tester,name,board,model,version,logPath,result,reason from " + params.TableName

	sqlStmt, condition := combineSQL(sqlSearch, params, params.TableName)
	rows, err = db.Query(sqlStmt, condition...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataResult []LocalResult
	for rows.Next() {
		var data LocalResult
		err := rows.Scan(&data.Id, &data.Time, &data.Tester, &data.Name, &data.Board, &data.Model, &data.Version, &data.LogPath, &data.PassOrFail, &data.Reason)
		if err != nil {
			return nil, err
		}
		dataResult = append(dataResult, data)
	}
	return dataResult, nil
}

func SearchTicketList(db *sql.DB) ([]TicketList, error) {
	sqlSearch := "select id,CaseName,Owner,Ticket from CaseTable"
	rows, err := db.Query(sqlSearch)
	if err != nil {
		return nil, err
	}
	var dataResult []TicketList
	for rows.Next() {
		var data TicketList
		rows.Scan(&data.Id, &data.CaseName, &data.Owner, &data.Ticket)
		dataResult = append(dataResult, data)
	}
	return dataResult, nil
}

// TODO: if startDate > endDate -> Bug
func SearchStainlessResult(db *sql.DB, params SearchInput) ([]StainlessResult, error) {
	var rows *sql.Rows
	var err error
	sqlSearch := "select id,time,duration,suite,board,model,buildVersion,host,testName,status,reason,firmwareROVersion,firmwareRWVersion from " + params.TableName

	sqlStmt, condition := combineSQL(sqlSearch, params, params.TableName)
	rows, err = db.Query(sqlStmt, condition...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataResult []StainlessResult
	for rows.Next() {
		var data StainlessResult
		rows.Scan(&data.Id, &data.Time, &data.Duration, &data.Suite, &data.Board, &data.Model, &data.BuildVersion, &data.Host, &data.TestName, &data.Status, &data.Reason, &data.FirmwareROVersion, &data.FirmwareRWVersion)
		dataResult = append(dataResult, data)
	}
	return dataResult, nil
}

func combineSQL(sqlSearch string, params SearchInput, tableName string) (string, []interface{}) {
	var whereCondition []string
	var whereValue []interface{}

	addCondition := func(condition string, value interface{}) {
		whereCondition = append(whereCondition, condition)
		whereValue = append(whereValue, value)
	}

	if params.Board != "" {
		addCondition(tableSearchMap[tableName][boardMap], params.Board)
	}
	if params.Name != "" {
		addCondition(tableSearchMap[tableName][testNameMap], params.Name)
	}
	if params.Result != "" {
		addCondition(tableSearchMap[tableName][resultMap], params.Result)
	}
	if params.Reason != "" {
		addCondition(tableSearchMap[tableName][reasonMap], fmt.Sprintf("%%%s%%", params.Reason))
	}

	if len(whereCondition) > 0 {
		sqlSearch += " WHERE " + strings.Join(whereCondition, " AND ")
	}

	if params.StartDate != "" && params.EndDate != "" {
		if len(whereCondition) > 0 {
			sqlSearch += " AND "
		} else {
			sqlSearch += " WHERE "
		}
		sqlSearch = sqlSearch + fmt.Sprintf("time BETWEEN '%s' AND '%s'", params.StartDate, params.EndDate)
	}

	orderBy := fmt.Sprintf(" order by %s desc", "id")
	if params.OrderBy != "" {
		orderBy = fmt.Sprintf(" order by %s desc", params.OrderBy)
	}
	sqlSearch += orderBy
	return sqlSearch, whereValue
}

func CatchPassAndFail(db *sql.DB, tableName, testName string) ([]PassingRate, error) {
	passingRate := make([]PassingRate, 0)
	days := grabLastWeeksDays(4)

	for pastWeek, day := range days {
		temp := PassingRate{
			Week:      pastWeek,
			StartDate: day.startDay,
			EndDate:   day.endDay,
		}

		errCountTotal := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s WHERE time BETWEEN '%s' AND '%s' AND testName = ?", tableName, day.startDay, day.endDay), testName).Scan(&temp.TotalRun)
		if errCountTotal != nil {
			return nil, fmt.Errorf("failed to count total test from %q to %q", day.startDay, day.endDay)
		}

		errCountPassing := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s WHERE time BETWEEN '%s' AND '%s' AND status = ? AND testName = ?", tableName, day.startDay, day.endDay), "Pass", testName).Scan(&temp.Pass)
		if errCountPassing != nil {
			return nil, fmt.Errorf("failed to count passing test from %q to %q", day.startDay, day.endDay)
		}

		if temp.TotalRun == 0 {
			temp.PassingRate = 0
		} else {
			temp.PassingRate = float64(temp.Pass) / float64(temp.TotalRun)
		}

		passingRate = append(passingRate, temp)
	}

	return passingRate, nil
}

func OnwerCaseMaintain(db *sql.DB, owner string) ([]PassingRate, error) {
	sqlOwnerCase := "select CaseName from CaseTable where owner = ?"
	rows, err := db.Query(sqlOwnerCase, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var caseList []string
	for rows.Next() {
		var caseName string
		rows.Scan(&caseName)
		caseList = append(caseList, caseName)
	}

	var passingRate []PassingRate
	var (
		countTotal int
		countPass  int
		tableName  = "Stainless_Result"
		statusPass = "Pass"
	)
	for _, caseName := range caseList {
		index := 1
		today := fmt.Sprintf("%v", now.BeginningOfDay().AddDate(0, 0, 0).Format("2006-01-02"))
		sevenDaysBefore := fmt.Sprintf("%v", now.BeginningOfDay().AddDate(0, 0, -7).Format("2006-01-02"))
		sqlCount := fmt.Sprintf("SELECT count(*) FROM %s where time between '%s' and '%s' AND testName = ?", tableName, sevenDaysBefore, today)
		errCountTotal := db.QueryRow(sqlCount, caseName).Scan(&countTotal)
		if errCountTotal != nil {
			return nil, fmt.Errorf("failed to count total runs for %q from %q to %q", caseName, today, sevenDaysBefore)
		}

		sqlCountPass := fmt.Sprintf("SELECT count(*) FROM %s where time between '%s' and '%s' AND testName = ? AND status = ?", tableName, sevenDaysBefore, today)
		errCountPass := db.QueryRow(sqlCountPass, caseName, statusPass).Scan(&countPass)
		if errCountPass != nil {
			return nil, fmt.Errorf("failed to count pass for %q from %q to %q", caseName, sevenDaysBefore, today)
		}

		var temp PassingRate = PassingRate{
			Week:        index,
			TestName:    caseName,
			StartDate:   today,
			EndDate:     sevenDaysBefore,
			TotalRun:    countTotal,
			Pass:        countPass,
			PassingRate: 0,
		}
		if countTotal == 0 {
			temp.PassingRate = 0
		} else {
			temp.PassingRate = float64(countPass) / float64(countTotal)
		}
		passingRate = append(passingRate, temp)
		index++
	}
	return passingRate, nil
}

func grabLastWeeksDays(week int) []weekDay {
	daySlice := []weekDay{}
	startOfTheWeek := now.BeginningOfDay()
	startOfTheWeek = startOfTheWeek.AddDate(0, 0, -7*week)

	for i := 0; i < week; i++ {
		weekFirstDay := startOfTheWeek.Format("2006-01-02")
		weekLastDay := startOfTheWeek.AddDate(0, 0, 7).Format("2006-01-02")
		daySlice = append(daySlice, weekDay{startDay: weekFirstDay, endDay: weekLastDay})
		startOfTheWeek = startOfTheWeek.AddDate(0, 0, 7)
	}
	return daySlice
}
