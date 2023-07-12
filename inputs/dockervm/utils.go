package dockervm

import (
	"os"
	"strings"
)

//	func formatFloat(value float64, precision int, carry *float64) float64 {
//		if carry != nil {
//			return round(value/(*carry), precision)
//		} else {
//			return round(value, precision)
//		}
//	}
//
//	func round(value float64, precision int) float64 {
//		pow := math.Pow(10, float64(precision))
//		return math.Round(value*pow) / pow
//	}
func getFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getFileLines(path string, strip bool, lineSep string) []string {
	data, err := getFileContent(path)
	if err != nil {
		return []string{}
	}
	if strip {
		data = strings.TrimSpace(data)
	}
	if lineSep == "" {
		return strings.Split(data, "\n")
	} else {
		return strings.Split(strings.TrimRight(data, lineSep), lineSep)
	}
}
