package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type ISqlRow interface {
	Close() error
	Next() bool
	NextResultSet() bool
}

type IDB[R ISqlRow] interface {
	Rows() (R, error)
	ScanRows(rows ISqlRow, dest interface{}) error
}

type IGormDB[R ISqlRow, T IDB[R]] interface {
	Rows() (R, error)
	Raw(sql string, values ...interface{}) T
	ScanRows(rows R, dest interface{}) error
}

var (
	IGNORE_FIELDS = []string{"state", "sizeCache", "unknownFields"}
	isDevelopment = true
)

func SetIsDevelopment(isDev bool) {
	isDevelopment = isDev
}

func Execute[R ISqlRow, T IDB[R]](db IGormDB[R, T], controller string, action string, claims IClaims, request interface{}, result interface{}) error {
	queryText := FindQueryWithinParamAndUser(controller, action, ToSqlScript(request, "Model", IGNORE_FIELDS...), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	rows, queryError := any(db.Raw(queryText)).(IDB[R]).Rows()
	if queryError != nil {
		consoleError("Execute", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	consoleQuery("Execute", controller, action, queryText)

	if errScan := db.ScanRows(rows, result); errScan != nil {
		consoleError("Execute", controller, action, queryText, errScan)
		return HandleSqlError(errScan)
	}

	return nil
}

func ExecuteId[R ISqlRow, T IDB[R]](db IGormDB[R, T], controller string, action string, claims IClaims, id interface{}, result interface{}) error {
	queryText := FindQueryWithinUser(controller, action, claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	rows, queryError := any(db.Raw(queryText, id)).(IDB[R]).Rows()
	if queryError != nil {
		consoleError("ExecuteId", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	consoleQuery("ExecuteId", controller, action, queryText)

	if errScan := db.ScanRows(rows, result); errScan != nil {
		consoleError("ExecuteId", controller, action, queryText, errScan)
		return HandleSqlError(errScan)
	}

	return nil
}

func ExecuteMultipleResult[R ISqlRow, T IDB[R]](db IGormDB[R, T], controller string, action string, claims IClaims, request interface{}, results ...interface{}) error {
	queryText := FindQueryWithinParamAndUser(controller, action, ToSqlScript(request, "Model", IGNORE_FIELDS...), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	rows, queryError := any(db.Raw(queryText)).(IDB[R]).Rows()
	if queryError != nil {
		consoleError("ExecuteMultipleResult", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	consoleQuery("ExecuteMultipleResult", controller, action, queryText)

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	errScan := scanResults(db, rows, results)
	if errScan != nil {
		consoleError("ExecuteMultipleResult", controller, action, queryText, errScan)
	}

	return errScan
}

func ExecuteIdMultipleResult[R ISqlRow, T IDB[R]](db IGormDB[R, T], controller string, action string, claims IClaims, request interface{}, results ...interface{}) error {
	queryText := FindQueryWithinUser(controller, action, claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	rows, queryError := any(db.Raw(queryText)).(IDB[R]).Rows()
	if queryError != nil {
		consoleError("ExecuteIdMultipleResult", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	consoleQuery("ExecuteIdMultipleResult", controller, action, queryText)

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	errScan := scanResults(db, rows, results)
	if errScan != nil {
		consoleError("ExecuteIdMultipleResult", controller, action, queryText, errScan)
	}

	return errScan
}

func FilterPagination[R ISqlRow, T IDB[R]](db IGormDB[R, T], controller string, action string, claims IClaims, filters interface{}, paging interface{}, total interface{}, items interface{}) error {
	builder := strings.Builder{}
	builder.WriteString(ToSqlScript(filters, "Filter", IGNORE_FIELDS...))
	builder.WriteString("\n")
	builder.WriteString(ToSqlScript(paging, "Pagination", IGNORE_FIELDS...))
	queryText := FindQueryWithinParamAndUser(controller, action, builder.String(), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	rows, queryError := any(db.Raw(queryText)).(IDB[R]).Rows()
	if queryError != nil {
		consoleError("FilterPagination", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	consoleQuery("FilterPagination", controller, action, queryText)

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	if errScan := db.ScanRows(rows, total); errScan != nil {
		consoleError("FilterPagination", controller, action, queryText, errScan)
		return HandleSqlError(errScan)
	}

	if rows.NextResultSet() && rows.Next() {
		if errScan := db.ScanRows(rows, items); errScan != nil {
			consoleError("FilterPagination", controller, action, queryText, errScan)
			return HandleSqlError(errScan)
		}
	}

	return nil
}

func replaceClaims(input string, claims IClaims) string {
	if claims == nil {
		return input
	}

	input = strings.ReplaceAll(input, "@@userID", fmt.Sprintf("%d", claims.GetId()))
	input = strings.ReplaceAll(input, "@@clientID", fmt.Sprintf("%d", claims.GetClientId()))
	input = strings.ReplaceAll(input, "@@unitID", fmt.Sprintf("%d", claims.GetUnitId()))
	input = strings.ReplaceAll(input, "@@username", fmt.Sprintf("'%s'", Safe(claims.GetUsername())))
	input = strings.ReplaceAll(input, "@@email", fmt.Sprintf("'%s'", Safe(claims.GetEmail())))
	input = strings.ReplaceAll(input, "@@fullName", fmt.Sprintf("'%s'", Safe(claims.GetFullname())))
	input = strings.ReplaceAll(input, "@@phone", fmt.Sprintf("'%s'", Safe(claims.GetPhone())))
	input = strings.ReplaceAll(input, "@@language", fmt.Sprintf("'%s'", Safe(claims.GetLanguage())))
	input = strings.ReplaceAll(input, "@@isAdmin", fmt.Sprintf("%d", BoolToInt(claims.GetIsAdmin())))
	input = strings.ReplaceAll(input, "@@isSystem", fmt.Sprintf("%d", BoolToInt(claims.GetIsSystem())))
	input = strings.ReplaceAll(input, "@@isBaseLanguage", fmt.Sprintf("%d", BoolToInt(claims.GetIsBaseLanguage())))
	input = strings.ReplaceAll(input, "[2]", IIF(claims.GetIsBaseLanguage(), "", "2"))
	return input
}

func scanResults[R ISqlRow, T IDB[R]](db IGormDB[R, T], rows R, results []interface{}) error {
	for i, e := range results {
		if i != 0 && !rows.Next() {
			rows.NextResultSet()
			continue
		}

		if err1 := db.ScanRows(rows, e); err1 != nil {
			return HandleSqlError(err1)
		}

		for rows.Next() {
			rows.Next()
		}

		if !rows.NextResultSet() {
			break
		}
	}

	return nil
}

func consoleError(method string, controller string, action string, queryText string, err error) {
	if isDevelopment {
		log.Printf("[%s] %s/%s: %s", method, controller, action, err.Error())
		log.Printf("Query: %s", queryText)
	}
}

func consoleQuery(method string, controller string, action string, queryText string) {
	if isDevelopment {
		fmt.Printf("%s/%s: %s", controller, action, queryText)
	}
}
