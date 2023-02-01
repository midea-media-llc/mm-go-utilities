package utils

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/midea-media-llc/mm-go-utilities/logs"
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
		logs.Errorf("[FindQuery]: %s: ", err.Error())
		return ""
	}

	if len(controllers.Controllers) < 1 {
		logs.Errorf("[FindQuery]: Controller Is Empty")
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

	logs.Errorf("[FindQuery]: Action %s not found", action)
	return ""
}

func loadXml(result interface{}, controller string) error {
	path, err := os.Getwd()
	if err != nil {
		logs.Errorf("[Load Xml]: %s", err.Error())
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
