package core

import "github.com/xuri/excelize/v2"

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
