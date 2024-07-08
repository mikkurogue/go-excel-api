package core

import (
	"errors"
	"fmt"
	"go-backend/util"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/xuri/excelize/v2"
)

func ChannelShipmentSheet(file string, resultChannel chan<- []Shipment) {
	time.Sleep(1 * time.Second)

	mapped, err := MapSheetHeadersToStructure(("shipments"), file)
	if err != nil {
		resultChannel <- nil
	}

	result, err := ReadShipmentsSheet(file, mapped)
	if err != nil {
		resultChannel <- nil
	}

	resultChannel <- result
}

func ChannelAssetSheet(file string, resultChannel chan<- []EmissionAsset) {
	time.Sleep(1 * time.Second)

	mapped, err := MapSheetHeadersToStructure(("emission assets"), file)
	if err != nil {
		resultChannel <- nil
	}

	result, err := ReadAssetsSheet(file, mapped)
	if err != nil {
		resultChannel <- nil
	}

	resultChannel <- result
}

func ProcessExcel(file string) (ExcelToJson, ParseError) {
	color.Magenta("Start processing file")
	defer util.TimeTrack(time.Now(), "ProcessExcel")

	shipmentChannel := make(chan []Shipment)
	assetChannel := make(chan []EmissionAsset)

	go ChannelShipmentSheet(file, shipmentChannel)
	go ChannelAssetSheet(file, assetChannel)

	encocded_shipments := <-shipmentChannel
	encocded_assets := <-assetChannel

	return ExcelToJson{
		Shipments: encocded_shipments,
		Assets:    encocded_assets,
	}, ParseError{}
}

func GetSheetHeaders(sheetName string, file string) []string {
	f, err := excelize.OpenFile(file)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows(sheetName)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return rows[0]
}

func ReadShipmentsSheet(file string, knownHeaders map[string]interface{}) ([]Shipment, error) {

	if knownHeaders == nil {
		return nil, errors.New("no known headers")
	}

	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, errors.New("file not found")
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Shipments")
	if err != nil {
		return nil, errors.New("no rows found")
	}

	// Map headers to their indices
	headerIndex := make(map[string]int)
	for i, header := range rows[0] {
		if _, ok := knownHeaders[header]; ok {
			headerIndex[header] = i
		}
	}

	var shipments []Shipment

	// Iterate over rows, skipping the header
	for _, row := range rows[1:] {
		shipment := Shipment{}
		v := reflect.ValueOf(&shipment).Elem()

		for json_field, col_index := range headerIndex {
			if col_index < len(row) {
				value := row[col_index]
				field := v.FieldByNameFunc(func(field_name string) bool {
					field, _ := v.Type().FieldByName(field_name)
					return field.Tag.Get("json") == json_field
				})
				if field.IsValid() && field.CanSet() {
					ptr := reflect.New(field.Type().Elem())
					ptr.Elem().SetString(value)
					field.Set(ptr)
				}
			}
		}

		shipments = append(shipments, shipment)
	}

	return shipments, nil
}

func ReadAssetsSheet(file string, knownHeaders map[string]interface{}) ([]EmissionAsset, error) {

	if knownHeaders == nil {
		return nil, errors.New("no known headers")
	}

	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, errors.New("file not found")
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Emission assets")
	if err != nil {
		return nil, errors.New("no rows found")
	}

	// Map headers to their indices
	headerIndex := make(map[string]int)
	for i, header := range rows[0] {
		if _, ok := knownHeaders[header]; ok {
			headerIndex[header] = i
		}
	}

	var assets []EmissionAsset

	// Iterate over rows, skipping the header
	for _, row := range rows[1:] {
		asset := EmissionAsset{}
		v := reflect.ValueOf(&asset).Elem()

		for json_field, col_index := range headerIndex {
			if col_index < len(row) {
				value := row[col_index]
				field := v.FieldByNameFunc(func(field_name string) bool {
					field, _ := v.Type().FieldByName(field_name)
					return field.Tag.Get("json") == json_field
				})
				if field.IsValid() && field.CanSet() {
					ptr := reflect.New(field.Type().Elem())
					ptr.Elem().SetString(value)
					field.Set(ptr)
				}
			}
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func FilterStructFields(obj interface{}, keys []string) (map[string]interface{}, error) {
	value := reflect.ValueOf(obj)
	typeValue := value.Type()

	result := make(map[string]interface{})

	for i := 0; i < value.NumField(); i++ {
		field := typeValue.Field(i)
		fieldValue := value.Field(i)

		jsonTag := field.Tag.Get("json")

		for _, key := range keys {
			if jsonTag == key {
				result[jsonTag] = fieldValue.Interface()
				break
			}
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no fields found")

	}

	return result, nil
}

func MapSheetHeadersToStructure(sheetName string, file string) (map[string]interface{}, error) {

	headers := GetSheetHeaders(sheetName, file)

	if len(headers) == 0 {
		return nil, errors.New("no headers found for sheet " + sheetName)
	}

	if strings.ToLower(sheetName) == "shipments" {
		result, err := FilterStructFields(Shipment{}, headers)
		if err != nil {
			return nil, err
		}

		return result, nil

	}

	if strings.ToLower(sheetName) == "emission assets" {
		result, err := FilterStructFields(EmissionAsset{}, headers)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("no known sheet can be processed. processing: " + sheetName)
}
