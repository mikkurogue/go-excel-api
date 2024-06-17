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

type ExcelFile struct {
	*excelize.File
}

type Shipment struct {
	ConsignmentId           *string `json:"consignment_id"`
	SequenceNumber          *string `json:"sequence_number"`
	Date                    *string `json:"date"`
	OriginLatitude          *string `json:"origin_latitude"`
	OriginLongitude         *string `json:"origin_longitude"`
	OriginIATA              *string `json:"origin_iata"`
	OriginUnlocode          *string `json:"origin_unlocode"`
	OriginRailTerminal      *string `json:"origin_rail_terminal"`
	OriginPostalCode        *string `json:"origin_postal_code"`
	OriginCity              *string `json:"origin_city"`
	OriginCountryCode       *string `json:"origin_country_code"`
	DestionationLatitude    *string `json:"destination_latitude"`
	DestinationLongitude    *string `json:"destination_longitude"`
	DestinationIATA         *string `json:"destination_iata"`
	DestinationUnlocode     *string `json:"destination_unlocode"`
	DestinationRailTerminal *string `json:"destination_rail_terminal"`
	DestinationPostalCode   *string `json:"destination_postal_code"`
	DestinationCity         *string `json:"destination_city"`
	DestinationCountryCode  *string `json:"destination_country_code"`
	Quantity                *string `json:"quantity"`
	Unit                    *string `json:"unit"`
	WeightInKg              *string `json:"weight_in_kg"`
	VolumeInM3              *string `json:"volume_in_m3"`
	EmissionInKgCo2e        *string `json:"emission_in_kg_co2e"`
	AssetKey1               *string `json:"asset_key_1"`
	AssetKey2               *string `json:"asset_key_2"`
	AssetKey3               *string `json:"asset_key_3"`
	AssetKey4               *string `json:"asset_key_4"`
	AssetKey5               *string `json:"asset_key_5"`
	AssetKey6               *string `json:"asset_key_6"`
	AssetKey7               *string `json:"asset_key_7"`
	AssetKey8               *string `json:"asset_key_8"`
	DebtorId                *string `json:"debtor_id"`
	CustomerName            *string `json:"customer_name"`
	CustomerGroup           *string `json:"customer_group"`
	TypeOfGoods             *string `json:"type_of_goods"`
	SkuNumber               *string `json:"sku_number"`
	CarrierName             *string `json:"carrier_name"`
	DossierNumber           *string `json:"dossier_number"`
	OriginAddress           *string `json:"origin_address"`
	DestinationAddress      *string `json:"destination_address"`
	NetworkType             *string `json:"network_type"`
	TransportLeg            *string `json:"transport_leg"`
	ServiceType             *string `json:"service_type"`
}

type EmissionAsset struct {
	AssetKey                     *string `json:"asset_key"`
	VehicleId                    *string `json:"vehicle_id"`
	HubName                      *string `json:"hub_name"`
	AssetPeriod                  *string `json:"asset_period"`
	AssetCarrier                 *string `json:"asset_carrier"`
	AssetInfo4                   *string `json:"asset_info_4"`
	AssetInfo5                   *string `json:"asset_info_5"`
	VehicleType                  *string `json:"vehicle_type"`
	VehicleCapacity              *string `json:"vehicle_capacity"`
	VehicleCargoType             *string `json:"vehicle_cargo_type"`
	VehicleEngineType            *string `json:"vehicle_engine_type"`
	IsTemperatureControlled      *string `json:"is_temperature_controlled"`
	IsExternalAsset              *string `json:"is_external_asset"`
	FuelType                     *string `json:"fuel_type"`
	FuelQuantity                 *string `json:"fuel_quantity"`
	FuelConsumptionKmPerLiter    *string `json:"fuel_consumption_km_per_liter"`
	PlannedDistanceInKm          *string `json:"planned_distance_in_km"`
	ActualDistanceInKm           *string `json:"actual_distance_in_km"`
	CO2PerTonKm                  *string `json:"co2_per_ton_km"`
	EmissionLocationM2           *string `json:"emission_location_m2"`
	EmissionLocationGeoReference *string `json:"emission_location_georeference"`
	EmissionAdditionalCategory   *string `json:"emission_additional_category"`
	EmissionOriginCountryCode    *string `json:"emission_origin_country_code"`
	EmissionYear                 *string `json:"emission_year"`
	HubType                      *string `json:"hub_type"`
}

type ExcelToJson struct {
	MappedShipmentHeaders map[string]interface{} `json:"mapped_shipment_headers"`
	MappedAssetHeaders    map[string]interface{} `json:"mapped_asset_headers"`
	Shipments             []Shipment             `json:"shipments"`
	Assets                []EmissionAsset        `json:"assets"`
}

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
