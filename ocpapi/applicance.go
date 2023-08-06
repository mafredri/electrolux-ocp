package ocpapi

type ApplianceInfo struct {
	PNC         string `json:"pnc"`
	Brand       string `json:"brand"`
	Market      string `json:"market"`
	ProductArea string `json:"productArea"`
	DeviceType  string `json:"deviceType"`
	Project     string `json:"project"`
	Model       string `json:"model"`
	Variant     string `json:"variant"`
	Colour      string `json:"colour"`
}

type Appliance struct {
	ApplianceID     string        `json:"applianceId"`
	ApplianceData   ApplianceData `json:"applianceData"`
	Properties      Properties    `json:"properties"`
	Status          string        `json:"status"`
	ConnectionState string        `json:"connectionState"`
}

type ApplianceData struct {
	ApplianceName string `json:"applianceName"`
	Created       string `json:"created"`
	ModelName     string `json:"modelName"`
}

type Properties struct {
	Desired  Desired     `json:"desired"`
	Reported Reported    `json:"reported"`
	Metadata interface{} `json:"metadata"`
}

type Desired struct {
	TimeZoneStandardName      string          `json:"TimeZoneStandardName"`
	LocationRequest           bool            `json:"LocationReq"`
	FirmwareVersionNIU        *string         `json:"FrmVer_NIU,omitempty"`
	Tasks                     interface{}     `json:"tasks,omitempty"`
	TimeZoneDaylightRule      string          `json:"TimeZoneDaylightRule"`
	Metadata                  DesiredMetadata `json:"$metadata"`
	Version                   int64           `json:"$version"`
	MinRefreshIntervalSeconds *int64          `json:"MinRefreshInterval_s,omitempty"`
	ReportExtraProperties     *bool           `json:"ReportExtraProperties,omitempty"`
	PM25_Hysteresis           *int64          `json:"PM2_5_Hysteresis,omitempty"`
	Monitoring                *bool           `json:"Monitoring,omitempty"`
	MonitoringStop            *int64          `json:"Monitoring_Stop,omitempty"`
	MonitoringStart           *int64          `json:"Monitoring_Start,omitempty"`
	VmNoNIU                   *string         `json:"VmNo_NIU,omitempty"`
}

type DesiredMetadata struct {
	DesiredMetadataUpdated
	TimeZoneStandardName  DesiredMetadataUpdated  `json:"TimeZoneStandardName"`
	LocationRequest       DesiredMetadataUpdated  `json:"LocationReq"`
	FirmwareVersionNIU    *DesiredMetadataUpdated `json:"FrmVer_NIU,omitempty"`
	Tasks                 *DesiredMetadataUpdated `json:"tasks,omitempty"`
	TimeZoneDaylightRule  DesiredMetadataUpdated  `json:"TimeZoneDaylightRule"`
	MinRefreshIntervalS   *DesiredMetadataUpdated `json:"MinRefreshInterval_s,omitempty"`
	ReportExtraProperties *DesiredMetadataUpdated `json:"ReportExtraProperties,omitempty"`
	PM25_Hysteresis       *DesiredMetadataUpdated `json:"PM2_5_Hysteresis,omitempty"`
	Monitoring            *DesiredMetadataUpdated `json:"Monitoring,omitempty"`
	MonitoringStop        *DesiredMetadataUpdated `json:"Monitoring_Stop,omitempty"`
	MonitoringStart       *DesiredMetadataUpdated `json:"Monitoring_Start,omitempty"`
	VmNoNIU               *DesiredMetadataUpdated `json:"VmNo_NIU,omitempty"`
}

type DesiredMetadataUpdated struct {
	LastUpdated        string `json:"$lastUpdated"`
	LastUpdatedVersion int64  `json:"$lastUpdatedVersion"`
}

type Reported struct {
	FirmwareVersionNIU        *string              `json:"FrmVer_NIU,omitempty"`
	Workmode                  string               `json:"Workmode"`
	FilterRFID                *string              `json:"FilterRFID,omitempty"`
	FilterLife                *int64               `json:"FilterLife,omitempty"`
	Fanspeed                  int64                `json:"Fanspeed"`
	UILight                   bool                 `json:"UILight"`
	SafetyLock                bool                 `json:"SafetyLock"`
	Ionizer                   *bool                `json:"Ionizer,omitempty"`
	FilterType                *int64               `json:"FilterType,omitempty"`
	ErrPM25                   *bool                `json:"ErrPM2_5,omitempty"`
	ErrTVOC                   *bool                `json:"ErrTVOC,omitempty"`
	ErrTempHumidity           *bool                `json:"ErrTempHumidity,omitempty"`
	ErrFanMtr                 *bool                `json:"ErrFanMtr,omitempty"`
	ErrCommSensorDisplayBoard *bool                `json:"ErrCommSensorDisplayBrd,omitempty"`
	DoorOpen                  *bool                `json:"DoorOpen,omitempty"`
	ErrRFID                   *bool                `json:"ErrRFID,omitempty"`
	SignalStrength            string               `json:"SignalStrength"`
	LogE                      *int64               `json:"logE,omitempty"`
	LogW                      *int64               `json:"logW,omitempty"`
	InterfaceVersion          int64                `json:"InterfaceVer"`
	VmNoNIU                   string               `json:"VmNo_NIU"`
	TVOCBrand                 *string              `json:"TVOCBrand,omitempty"`
	Capabilities              ReportedCapabilities `json:"capabilities"`
	Tasks                     interface{}          `json:"tasks"`
	Metadata                  ReportedMetadata     `json:"$metadata"`
	Version                   int64                `json:"$version"`
	DeviceID                  string               `json:"deviceId"`
	TVOC                      *int64               `json:"TVOC,omitempty"`
	CO2                       *int64               `json:"CO2,omitempty"` // Deprecated: use ECO2 (latest firmware).
	PM1                       *int64               `json:"PM1,omitempty"`
	PM25                      *int64               `json:"PM2_5,omitempty"`
	PM10                      *int64               `json:"PM10,omitempty"`
	Humidity                  *int64               `json:"Humidity,omitempty"`
	Temp                      *int64               `json:"Temp,omitempty"`
	RSSI                      *int64               `json:"RSSI,omitempty"`
	ECO2                      *int64               `json:"ECO2,omitempty"`
	FilterLife1               *int64               `json:"FilterLife_1,omitempty"`
	Monitoring                *bool                `json:"Monitoring,omitempty"`
	MonitoringStop            *int64               `json:"Monitoring_Stop,omitempty"`
	MonitoringStart           *int64               `json:"Monitoring_Start,omitempty"`
	UVState                   *string              `json:"UVState,omitempty"`
	UVRuntime                 *int64               `json:"UVRuntime,omitempty"`
	ErrCommSensorUIBrd        *string              `json:"ErrCommSensorUIBrd,omitempty"`
	ErrImpellerStuck          *string              `json:"ErrImpellerStuck,omitempty"`
	ErrPmNotResp              *string              `json:"ErrPmNotResp,omitempty"`
	VmNoMCU                   *string              `json:"VmNo_MCU,omitempty"`
	PM25_Approximate          *int64               `json:"PM2_5_approximate,omitempty"`
}

type ReportedCapabilities struct {
	Tasks interface{} `json:"tasks,omitempty"`
}

type ReportedMetadata struct {
	ReportedMetadataUpdated
	FirmwareVersionNIU        *ReportedMetadataUpdated     `json:"FrmVer_NIU,omitempty"`
	Workmode                  ReportedMetadataUpdated      `json:"Workmode"`
	FilterRFID                *ReportedMetadataUpdated     `json:"FilterRFID,omitempty"`
	FilterLife                *ReportedMetadataUpdated     `json:"FilterLife,omitempty"`
	Fanspeed                  ReportedMetadataUpdated      `json:"Fanspeed"`
	UILight                   ReportedMetadataUpdated      `json:"UILight"`
	SafetyLock                ReportedMetadataUpdated      `json:"SafetyLock"`
	Ionizer                   *ReportedMetadataUpdated     `json:"Ionizer,omitempty"`
	FilterType                *ReportedMetadataUpdated     `json:"FilterType,omitempty"`
	ErrPM25                   *ReportedMetadataUpdated     `json:"ErrPM2_5,omitempty"`
	ErrTVOC                   *ReportedMetadataUpdated     `json:"ErrTVOC,omitempty"`
	ErrTempHumidity           *ReportedMetadataUpdated     `json:"ErrTempHumidity,omitempty"`
	ErrFanMtr                 *ReportedMetadataUpdated     `json:"ErrFanMtr,omitempty"`
	ErrCommSensorDisplayBoard *ReportedMetadataUpdated     `json:"ErrCommSensorDisplayBrd,omitempty"`
	DoorOpen                  *ReportedMetadataUpdated     `json:"DoorOpen,omitempty"`
	ErrRFID                   *ReportedMetadataUpdated     `json:"ErrRFID,omitempty"`
	SignalStrength            ReportedMetadataUpdated      `json:"SignalStrength"`
	LogE                      *ReportedMetadataUpdated     `json:"logE,omitempty"`
	LogW                      *ReportedMetadataUpdated     `json:"logW,omitempty"`
	InterfaceVersion          ReportedMetadataUpdated      `json:"InterfaceVer"`
	VmNoNIU                   ReportedMetadataUpdated      `json:"VmNo_NIU"`
	TVOCBrand                 *ReportedMetadataUpdated     `json:"TVOCBrand,omitempty"`
	Capabilities              ReportedMetadataCapabilities `json:"capabilities"`
	Tasks                     ReportedMetadataUpdated      `json:"tasks"`
	TVOC                      *ReportedMetadataUpdated     `json:"TVOC,omitempty"`
	CO2                       *ReportedMetadataUpdated     `json:"CO2,omitempty"`
	PM1                       *ReportedMetadataUpdated     `json:"PM1,omitempty"`
	PM25                      *ReportedMetadataUpdated     `json:"PM2_5,omitempty"`
	PM10                      *ReportedMetadataUpdated     `json:"PM10,omitempty"`
	Humidity                  *ReportedMetadataUpdated     `json:"Humidity,omitempty"`
	Temp                      *ReportedMetadataUpdated     `json:"Temp,omitempty"`
	RSSI                      *ReportedMetadataUpdated     `json:"RSSI,omitempty"`
	ECO2                      *ReportedMetadataUpdated     `json:"ECO2,omitempty"`
	FilterLife1               *ReportedMetadataUpdated     `json:"FilterLife_1,omitempty"`
	Monitoring                *ReportedMetadataUpdated     `json:"Monitoring,omitempty"`
	MonitoringStop            *ReportedMetadataUpdated     `json:"Monitoring_Stop,omitempty"`
	MonitoringStart           *ReportedMetadataUpdated     `json:"Monitoring_Start,omitempty"`
	UVState                   *ReportedMetadataUpdated     `json:"UVState,omitempty"`
	UVRuntime                 *ReportedMetadataUpdated     `json:"UVRuntime,omitempty"`
	ErrCommSensorUIBrd        *ReportedMetadataUpdated     `json:"ErrCommSensorUIBrd,omitempty"`
	ErrImpellerStuck          *ReportedMetadataUpdated     `json:"ErrImpellerStuck,omitempty"`
	ErrPmNotResp              *ReportedMetadataUpdated     `json:"ErrPmNotResp,omitempty"`
	VmNoMCU                   *ReportedMetadataUpdated     `json:"VmNo_MCU,omitempty"`
	PM25_Approximate          *ReportedMetadataUpdated     `json:"PM2_5_approximate,omitempty"`
}

type ReportedMetadataCapabilities struct {
	ReportedMetadataUpdated
	Tasks *ReportedMetadataUpdated `json:"tasks,omitempty"`
}

type ReportedMetadataUpdated struct {
	LastUpdated string `json:"$lastUpdated"`
}
