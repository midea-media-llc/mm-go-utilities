package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// UnmarshallJSONString decode json string
func UnmarshallJSONString(jsonString string, item interface{}) error {
	if jsonString == "" {
		return nil
	}

	bytes := []byte(jsonString)

	if err := json.Unmarshal(bytes, item); err != nil {
		return err
	}

	return nil
}

func ExtractToken(bearerToken string) string {
	t := strings.Split(bearerToken, " ")
	if len := len(t); len < 2 {
		return ""
	}
	return t[1]
}

func GetTimeNowWithLocation() (time.Time, error) {
	return time.Now().In(Location), nil
}

func Now() time.Time {
	return time.Now().In(Location)
}

// func ConvertStringToDatetimeWithLocation(dateTime string, layout string) (time.Time, error) {
// 	// Parse string to datetime with location
// 	date, err := time.ParseInLocation(layout, dateTime, Location)
// 	if err != nil {
// 		return time.Time{}, common.NewBadRequestError("error_parse_datetime")
// 	}
// 	return date, nil
// }

func ConvertInt64ArrayToString(a []int64, delim string) string {
	return strings.Replace(fmt.Sprint(a), " ", delim, -1)
	// return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func SubtractTimeFromInt64(seconds int64) int64 {
	t := time.Unix(seconds, 0)
	return seconds - int64(t.Hour()*3600+t.Minute()*60+t.Second())
}

func ArrInt64ToString(input []int64, delim string) string {
	result := []string{}

	for _, val := range input {
		result = append(result, fmt.Sprintf("%d", val))
	}

	return strings.Join(result, delim)
}

func StringToArrayInt(input string, delim string) ([]int, error) {
	items := strings.Split(input, delim)
	result := make([]int, len(items))
	for i, item := range items {
		value, err := strconv.Atoi(item)
		if err != nil {
			return result, err
		}

		result[i] = value
	}
	return result, nil
}

func JoinArrayIntToString(input []int, delim string) string {
	result := []string{}

	for _, val := range input {
		result = append(result, fmt.Sprintf("%d", val))
	}

	return strings.Join(result, delim)
}

func ArrayStringToInt(arr []string) []int {
	var t2 = []int{}

	for _, i := range arr {
		if i != "" {
			j, err := strconv.Atoi(strings.TrimSpace(i))
			if err != nil {
				panic(err)
			}
			t2 = append(t2, j)
		}
	}

	return t2
}

func ArrayContains(arr []interface{}, finder interface{}) bool {
	if arr == nil {
		return false
	}

	for _, item := range arr {
		if finder == item {
			return true
		}
	}
	return false
}

func ArrayStringContains(arr []string, finder string) bool {
	if arr == nil {
		return false
	}

	for _, item := range arr {
		if finder == item {
			return true
		}
	}
	return false
}

func ArrayIntContains(arr []int, finder int) bool {
	if arr == nil {
		return false
	}

	for _, item := range arr {
		if finder == item {
			return true
		}
	}
	return false
}

func ArrayInt64Contains(arr []int64, finder int64) bool {
	if arr == nil {
		return false
	}

	for _, item := range arr {
		if finder == item {
			return true
		}
	}
	return false
}

func ReplaceAllSpecialChar(inputText string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return inputText
	}
	return reg.ReplaceAllString(inputText, "")

}

func UniqueIn64Array(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ConvertTimeDurationToString(fromDay time.Time, toDay time.Time, withDay bool, withHour bool, withMinute bool) string {
	var res string = ""

	hs := toDay.Sub(fromDay).Hours()
	days := hs / 24

	days, hf := math.Modf(days)
	hours, mf := math.Modf(hf * 24)

	mins := mf * 60

	// fmt.Println(days, "days", hours, "hours", math.Round(mins), "minutes")
	if withDay {
		res = res + fmt.Sprint(days) + " ngày"
	}
	if withHour {
		res = res + " " + fmt.Sprint(math.Round(hours)) + " giờ"
	}
	if withMinute {
		res = res + " " + fmt.Sprint(math.Round(mins)) + " phút"
	}
	return res
}

func RoundNumber(money interface{}) int64 {

	valStr := fmt.Sprintf("%.100f", money)

	moneyFloat64, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0
	}
	temp := math.Round(moneyFloat64)

	return int64(temp)
}

func HandleErrorMessage(err error) string {
	return strings.ReplaceAll(strings.Split(err.Error(), ";")[0], "mssql: ", "")
}

func HandleNewErrorMessage(err error) error {
	return errors.New(HandleErrorMessage(err))
}
