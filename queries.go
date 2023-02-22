package utils

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

var REMOVE_PATHS = []string{"cmd/main", "cmd\\main"}

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

// FindQuery reads an XML file containing controller and action data,
// and returns the query string associated with a given controller and action.
// It takes the following arguments:
// - controller: the name of the XML file (without the extension) to be read
// - action: the name of the action to find in the XML file
// The function returns the query string associated with the given controller and action,
// or an empty string if the XML file cannot be read or the controller and action cannot be found.
func FindQuery(controller string, action string) string {
	// Load the XML file into the XmlControllers struct.
	controllers := XmlControllers{}
	err := loadXml(&controllers, controller)

	// Return an empty string if there was an error loading the XML file.
	if err != nil {
		return ""
	}

	// Return an empty string if there are no controllers in the XmlControllers struct.
	if len(controllers.Controllers) < 1 {
		return ""
	}

	// Iterate through each controller in the XmlControllers struct.
	for _, contr := range controllers.Controllers {
		// If the controller name matches the one we're looking for, or if it's the "Base" controller...
		if contr.Name == controller || contr.Name == "Base" {
			// Iterate through each action in the controller.
			for _, act := range contr.Actions {
				// If the action name matches the one we're looking for, return the query string.
				if act.Name == action {
					return act.Text
				}
			}
		}
	}

	// Return an empty string if the controller and action couldn't be found.
	return ""
}

// FindQueryWithinParam reads an XML file containing controller and action data,
// replaces a placeholder string in the query string with a given parameter value,
// and returns the resulting query string.
// It takes the following arguments:
// - controller: the name of the XML file (without the extension) to be read
// - action: the name of the action to find in the XML file
// - queryParam: the value to replace the "[QUERY_PARAMS]" placeholder string with
// The function returns the resulting query string,
// or an empty string if the XML file cannot be read or the controller and action cannot be found.
func FindQueryWithinParam(controller string, action string, queryParam string) string {
	// Get the query string associated with the given controller and action.
	result := FindQuery(controller, action)

	// Replace the "[QUERY_PARAMS]" placeholder string with the given queryParam.
	result = strings.ReplaceAll(result, "[QUERY_PARAMS]", queryParam)

	// Return the resulting query string.
	return result
}

// FindQueryWithinParamAndUser reads an XML file containing controller and action data,
// replaces a placeholder string in the query string with a given parameter value,
// replaces another placeholder string in the resulting string with a given user object,
// and returns the resulting query string.
// It takes the following arguments:
//   - controller: the name of the XML file (without the extension) to be read
//   - action: the name of the action to find in the XML file
//   - queryParam: the value to replace the "[QUERY_PARAMS]" placeholder string with
//   - user: an object representing the user to replace the "[USER]" placeholder string with
//   - replaceUserFunc: a function that takes a string and a user object,
//     and returns the string with the "[USER]" placeholder string replaced with the user object
//
// The function returns the resulting query string,
// or an empty string if the XML file cannot be read or the controller and action cannot be found.
func FindQueryWithinParamAndUser[T any](controller string, action string, queryParam string, user T, replaceUserFunc func(string, T) string) string {
	// Get the query string associated with the given controller and action,
	// and replace the "[QUERY_PARAMS]" placeholder string with the given queryParam.
	result := FindQueryWithinParam(controller, action, queryParam)

	// Replace the "[USER]" placeholder string with the given user object.
	result = replaceUserFunc(result, user)

	// Return the resulting query string.
	return result
}

// FindQueryWithinUser reads an XML file containing controller and action data,
// replaces another placeholder string in the query string with a given user object,
// and returns the resulting query string.
// It takes the following arguments:
//   - controller: the name of the XML file (without the extension) to be read
//   - action: the name of the action to find in the XML file
//   - user: an object representing the user to replace the "[USER]" placeholder string with
//   - replaceUserFunc: a function that takes a string and a user object,
//     and returns the string with the "[USER]" placeholder string replaced with the user object
//
// The function returns the resulting query string,
// or an empty string if the XML file cannot be read or the controller and action cannot be found.
func FindQueryWithinUser[T any](controller string, action string, user T, replaceUserFunc func(string, T) string) string {
	// Get the query string associated with the given controller and action.
	result := FindQuery(controller, action)

	// Replace the "[USER]" placeholder string with the given user object.
	result = replaceUserFunc(result, user)

	// Return the resulting query string.
	return result
}

// loadXml reads an XML file at a given path and unmarshals its contents into an interface.
// It takes the following arguments:
// - result: a pointer to an interface that will hold the unmarshalled data
// - controller: the name of the XML file (without the extension) to be read
// The function returns an error if it fails to read or unmarshal the XML file.
func loadXml(result interface{}, controller string) error {
	// Get the current working directory.
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// Remove certain paths from the working directory.
	for _, e := range REMOVE_PATHS {
		path = strings.ReplaceAll(path, e, "")
	}

	// Construct the path to the XML file.
	filePath := fmt.Sprintf("%s/xml/%s.xml", path, controller)

	// Read the contents of the XML file.
	xmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the XML into the provided result interface.
	xml.Unmarshal(xmlBytes, result)

	// Return nil if everything succeeded.
	return nil
}
