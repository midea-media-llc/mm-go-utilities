package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	IGNORE_FIELDS = []string{"state", "sizeCache", "unknownFields"}
	isDevelopment = true
)

func SetIsDevelopment(isDev bool) {
	isDevelopment = isDev
}

func Execute(db IGormDB, controller string, action string, claims IClaims, request interface{}, result interface{}) error {
	queryText := FindQueryWithinParamAndUser(controller, action, ToSqlScript(request, "Model", IGNORE_FIELDS...), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("Execute", controller, action, queryText)

	rows, queryError := db.Raw(queryText).Rows()
	if queryError != nil {
		consoleError("Execute", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	if err := rows.Scan(result); err != nil {
		consoleError("Execute", controller, action, queryText, queryError)
		return HandleSqlError(err)
	}

	return nil
}

func ExecuteId(db IGormDB, controller string, action string, claims IClaims, id interface{}, result interface{}) error {
	queryText := FindQueryWithinUser(controller, action, claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("ExecuteId", controller, action, queryText)

	rows, queryError := db.Raw(queryText).Rows()
	if queryError != nil {
		consoleError("ExecuteId", controller, action, queryText, queryError)
		return HandleSqlError(queryError)
	}

	if err := rows.Scan(result); err != nil {
		consoleError("ExecuteId", controller, action, queryText, queryError)
		return HandleSqlError(err)
	}

	return nil
}

func ExecuteMultipleResult(db IGormDB, controller string, action string, claims IClaims, request interface{}, results ...interface{}) error {
	queryText := FindQueryWithinParamAndUser(controller, action, ToSqlScript(request, "Model", IGNORE_FIELDS...), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("ExecuteMultipleResult", controller, action, queryText)

	rows, err := db.Raw(queryText).Rows()
	if err != nil {
		consoleError("ExecuteMultipleResult", controller, action, queryText, err)
		return HandleSqlError(err)
	}

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	errScan := scanResults(rows, results)
	if errScan != nil {
		consoleError("ExecuteMultipleResult", controller, action, queryText, err)
	}

	return errScan
}

func ExecuteIdMultipleResult(db IGormDB, controller string, action string, claims IClaims, request interface{}, results ...interface{}) error {
	queryText := FindQueryWithinUser(controller, action, claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("ExecuteIdMultipleResult", controller, action, queryText)

	rows, err := db.Raw(queryText, request).Rows()
	if err != nil {
		consoleError("ExecuteIdMultipleResult", controller, action, queryText, err)
		return HandleSqlError(err)
	}

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	errScan := scanResults(rows, results)
	if errScan != nil {
		consoleError("ExecuteIdMultipleResult", controller, action, queryText, err)
	}

	return errScan
}

func FilterPagination(db IGormDB, controller string, action string, claims IClaims, filters interface{}, paging interface{}, total interface{}, items interface{}) error {
	builder := strings.Builder{}
	builder.WriteString(ToSqlScript(filters, "Filter", IGNORE_FIELDS...))
	builder.WriteString("\n")
	builder.WriteString(ToSqlScript(paging, "Pagination", IGNORE_FIELDS...))
	queryText := FindQueryWithinParamAndUser(controller, action, builder.String(), claims, replaceClaims)
	if queryText == "" {
		return errors.New("action_not_found")
	}

	consoleQuery("FilterPagination", controller, action, queryText)

	rows, err := db.Raw(queryText).Rows()
	if err != nil {
		consoleError("FilterPagination", controller, action, queryText, err)
		return HandleSqlError(err)
	}

	defer rows.Close()
	if !rows.Next() {
		return nil
	}

	if err1 := rows.Scan(rows, total); err1 != nil {
		consoleError("FilterPagination", controller, action, queryText, err1)
		return HandleSqlError(err1)
	}

	if rows.NextResultSet() && rows.Next() {
		if err2 := rows.Scan(rows, items); err2 != nil {
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

func scanResults(rows *sql.Rows, results []interface{}) error {
	for i, e := range results {
		if i != 0 && !rows.Next() {
			rows.NextResultSet()
			continue
		}

		if err1 := rows.Scan(rows, e); err1 != nil {
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
