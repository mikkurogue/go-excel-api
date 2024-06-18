package core

import (
	"fmt"
	"go-backend/util"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/xuri/excelize/v2"
)

func ChannelShipmentSheet(file string, result_channel chan<- []Shipment) {
	time.Sleep(1 * time.Second)
	result := ReadShipmentsSheet(file, MapSheetHeadersToStructure("shipments", file))
	result_channel <- result
}

func ChannelAssetSheet(file string, result_channel chan<- []EmissionAsset) {
	time.Sleep(1 * time.Second)
	result := ReadAssetsSheet(file, MapSheetHeadersToStructure("emission assets", file))
	result_channel <- result
}

func ProcessExcel(file string) ExcelToJson {
	color.Magenta("Start processing file")
	defer util.TimeTrack(time.Now(), "ProcessExcel")

	shipment_channel := make(chan []Shipment)
	asset_channel := make(chan []EmissionAsset)

	go ChannelShipmentSheet(file, shipment_channel)
	go ChannelAssetSheet(file, asset_channel)

	encocded_shipments := <-shipment_channel
	encocded_assets := <-asset_channel

	return ExcelToJson{
		Shipments:             encocded_shipments,
		Assets:                encocded_assets,
		MappedShipmentHeaders: MapSheetHeadersToStructure("shipments", file),
		MappedAssetHeaders:    MapSheetHeadersToStructure("emission assets", file),
	}
}

func GetSheetHeaders(sheet_name string, file string) []string {
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

	rows, err := f.GetRows(sheet_name)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return rows[0]
}

func ReadShipmentsSheet(file string, known_headers map[string]interface{}) []Shipment {

	if known_headers == nil {
		return nil
	}

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

	rows, err := f.GetRows("Shipments")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Map headers to their indices
	header_index := make(map[string]int)
	for i, header := range rows[0] {
		if _, ok := known_headers[header]; ok {
			header_index[header] = i
		}
	}

	var shipments []Shipment

	// Iterate over rows, skipping the header
	for _, row := range rows[1:] {
		shipment := Shipment{}
		v := reflect.ValueOf(&shipment).Elem()

		for json_field, col_index := range header_index {
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

	return shipments
}

func ReadAssetsSheet(file string, known_headers map[string]interface{}) []EmissionAsset {

	if known_headers == nil {
		return nil
	}

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

	rows, err := f.GetRows("Emission assets")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Map headers to their indices
	header_index := make(map[string]int)
	for i, header := range rows[0] {
		if _, ok := known_headers[header]; ok {
			header_index[header] = i
		}
	}

	var assets []EmissionAsset

	// Iterate over rows, skipping the header
	for _, row := range rows[1:] {
		asset := EmissionAsset{}
		v := reflect.ValueOf(&asset).Elem()

		for json_field, col_index := range header_index {
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

	return assets
}

func FilterStructFields(obj interface{}, keys []string) map[string]interface{} {
	value := reflect.ValueOf(obj)
	type_value := value.Type()

	result := make(map[string]interface{})

	for i := 0; i < value.NumField(); i++ {
		field := type_value.Field(i)
		fieldValue := value.Field(i)

		jsonTag := field.Tag.Get("json")

		for _, key := range keys {
			if jsonTag == key {
				result[jsonTag] = fieldValue.Interface()
				break
			}
		}
	}
	return result
}

func MapSheetHeadersToStructure(sheet_name string, file string) map[string]interface{} {

	headers := GetSheetHeaders(sheet_name, file)

	if len(headers) == 0 {
		fmt.Println("No headers found for sheet", sheet_name)
		return nil
	}

	var result map[string]interface{}

	if strings.ToLower(sheet_name) == "shipments" {
		result = FilterStructFields(Shipment{}, headers)
	}

	if strings.ToLower(sheet_name) == "emission assets" {
		result = FilterStructFields(EmissionAsset{}, headers)
	}

	return result

	// fmt.Println("Headers for sheet", sheetName, "are:", headers)

}
