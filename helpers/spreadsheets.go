package helpers

import (
	"fmt"
	"reflect"

	"github.com/xuri/excelize/v2"
)

func GenerateExcelFile(report interface{}, sheetName, localFilePath string) (string, error) {
	f := excelize.NewFile()
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating Excel file: %w", err)
	}

	v := reflect.ValueOf(report)
	if v.Kind() != reflect.Slice {
		return "", fmt.Errorf("provided report is not a slice")
	}

	// Extract headers from the struct fields
	if v.Len() > 0 {
		headers := getHeaders(v.Index(0).Interface())
		for i, header := range headers {
			cell := fmt.Sprintf("%s%d", fmt.Sprint('A'+i), 1) // A1, B1, C1, etc.
			f.SetCellValue(sheetName, cell, header)
		}

		// Fill in the data
		for i := 0; i < v.Len(); i++ {
			row := v.Index(i).Interface()
			for j, header := range headers {
				cell := fmt.Sprintf("%s%d", fmt.Sprint('A'+j), i+2) // Start from row 2
				value := reflect.ValueOf(row).FieldByName(header)
				f.SetCellValue(sheetName, cell, value.Interface())
			}
		}
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs(localFilePath); err != nil {
		return "", fmt.Errorf("error saving Excel file: %w", err)
	}

	return localFilePath, nil
}

func getHeaders(data interface{}) []string {
	val := reflect.TypeOf(data)
	var headers []string

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			headers = append(headers, val.Field(i).Name)
		}
	}
	return headers
}
