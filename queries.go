package utils

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

const (
	REMOVE_PATH = "cmd/main"
	APPEND_PATH = "/"
)

type XmlNameNode struct {
	Name string `xml:"name,attr"`
}

type XmlAction struct {
	XmlNameNode
	XMLName xml.Name `xml:"action"`
	Text    string   `xml:"text"`
}

type XmlController struct {
	XmlNameNode
	XMLName xml.Name    `xml:"controller"`
	Actions []XmlAction `xml:"action"`
}

type XmlControllers struct {
	XmlNameNode
	XMLName     xml.Name        `xml:"controllers"`
	Controllers []XmlController `xml:"controller"`
}

func FindQueryWithinParam(controller string, action string, queryParam string) string {
	result := FindQuery(controller, action)
	return strings.ReplaceAll(result, "[QUERY_PARAMS]", queryParam)
}

func FindQuery(controller string, action string) string {
	controllers := XmlControllers{}
	err := loadXml(&controllers, controller)

	if err != nil {
		return ""
	}

	if len(controllers.Controllers) < 1 {
		return ""
	}

	for _, contr := range controllers.Controllers {
		if contr.Name == controller || contr.Name == "Base" {
			for _, act := range contr.Actions {
				if act.Name == action {
					return act.Text
				}
			}
		}
	}

	return ""
}

func FindQueryWithinParamAndUser[T any](controller string, action string, queryParam string, user T, replaceUserFunc func(string, T) string) string {
	result := FindQueryWithinParam(controller, action, queryParam)
	return replaceUserFunc(result, user)
}

func FindQueryWithinUser[T any](controller string, action string, user T, replaceUserFunc func(string, T) string) string {
	result := FindQuery(controller, action)
	return replaceUserFunc(result, user)
}

func loadXml(result interface{}, controller string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, REMOVE_PATH, "") + APPEND_PATH

	filePath := fmt.Sprintf("%sxml/%s.xml", path, controller)
	xmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	xml.Unmarshal(xmlBytes, result)
	return nil
}
