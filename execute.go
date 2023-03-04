package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

var (
	IGNORE_FIELDS = []string{"state", "sizeCache", "unknownFields"}
	isDevelopment = true
)

func SetIsDevelopment(isDev bool) {
	isDevelopment = isDev
}

func Execute(db *gorm.DB, controller string, action string, claims IClaims, request interface{}, result interface{}) error {
	queryText := FindQueryWithinParamAndUser(controller, action, ToSqlScript(request, "Model", IGNORE_FIELDS...), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("Execute", controller, action, queryText)
	query := db.Raw(queryText).Scan(result)
	if query.Error != nil {
		consoleError("Execute", controller, action, queryText, query.Error)
		return HandleSqlError(query.Error)
	}
	return nil
}

func ExecuteId(db *gorm.DB, controller string, action string, claims IClaims, id interface{}, result interface{}) error {
	queryText := FindQueryWithinUser(controller, action, claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("ExecuteId", controller, action, queryText)
	query := db.Raw(queryText, id).Scan(result)
	if query.Error != nil {
		consoleError("ExecuteId", controller, action, queryText, query.Error)
		return HandleSqlError(query.Error)
	}
	return nil
}

func ExecuteMultipleResult(db *gorm.DB, controller string, action string, claims IClaims, request interface{}, results ...interface{}) error {
	queryText := FindQueryWithinParamAndUser(controller, action, ToSqlScript(request, "Model", IGNORE_FIELDS...), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("ExecuteMultipleResult", controller, action, queryText)
	query := db.Raw(queryText)
	rows, err := query.Rows()
	if err != nil {
		consoleError("ExecuteMultipleResult", controller, action, queryText, err)
		return HandleSqlError(err)
	}

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	errScan := scanResults(rows, query, results)
	if errScan != nil {
		consoleError("ExecuteMultipleResult", controller, action, queryText, err)
	}

	return errScan
}

func ExecuteIdMultipleResult(db *gorm.DB, controller string, action string, claims IClaims, request interface{}, results ...interface{}) error {
	queryText := FindQueryWithinUser(controller, action, claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("ExecuteIdMultipleResult", controller, action, queryText)
	query := db.Raw(queryText, request)
	rows, err := query.Rows()
	if err != nil {
		consoleError("ExecuteIdMultipleResult", controller, action, queryText, err)
		return HandleSqlError(query.Error)
	}

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	errScan := scanResults(rows, query, results)
	if errScan != nil {
		consoleError("ExecuteIdMultipleResult", controller, action, queryText, err)
	}

	return errScan
}

func FilterPagination(db *gorm.DB, controller string, action string, claims IClaims, filters interface{}, paging interface{}, total interface{}, items interface{}) error {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(ToSqlScript(filters, "Filter", IGNORE_FIELDS...))
	queryBuilder.WriteString("\n")
	queryBuilder.WriteString(ToSqlScript(paging, "Pagination", IGNORE_FIELDS...))
	queryText := FindQueryWithinParamAndUser(controller, action, queryBuilder.String(), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("FilterPagination", controller, action, queryText)
	query := db.Raw(queryText)
	rows, err := query.Rows()
	if err != nil {
		consoleError("FilterPagination", controller, action, queryText, err)
		return HandleSqlError(err)
	}

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	if err1 := query.ScanRows(rows, total); err1 != nil {
		consoleError("FilterPagination", controller, action, queryText, err1)
		return HandleSqlError(err1)
	}

	if rows.NextResultSet() && rows.Next() {
		if err2 := query.ScanRows(rows, items); err2 != nil {
			consoleError("FilterPagination", controller, action, queryText, err2)
			return HandleSqlError(err2)
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

func scanResults(rows *sql.Rows, query *gorm.DB, results []interface{}) error {
	for i, e := range results {
		if i != 0 && !rows.Next() {
			rows.NextResultSet()
			continue
		}

		if err1 := query.ScanRows(rows, e); err1 != nil {
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
