package ocpapi

import "time"

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

type ApplianceID string

func (id ApplianceID) String() string {
	return string(id)
}

func (id ApplianceID) Serial() string {
	if len(id) < 17 {
		return ""
	}
	// Example: 950011538111111115087076 -> 11111111
	return string(id)[9:17]
}

func (id ApplianceID) PNC() string {
	if len(id) < 9 {
		return ""
	}
	// Example: 950011538111111115087076 -> 950011538
	return string(id)[0:9]
}

type Appliance struct {
	ApplianceID     ApplianceID   `json:"applianceId"`
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
	Version                   int             `json:"$version"`
	MinRefreshIntervalSeconds *int            `json:"MinRefreshInterval_s,omitempty"`
	ReportExtraProperties     *bool           `json:"ReportExtraProperties,omitempty"`
	PM25Hysteresis            *int            `json:"PM2_5_Hysteresis,omitempty"`
	Monitoring                *bool           `json:"Monitoring,omitempty"`
	MonitoringStop            *int            `json:"Monitoring_Stop,omitempty"`
	MonitoringStart           *int            `json:"Monitoring_Start,omitempty"`
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
	PM25Hysteresis        *DesiredMetadataUpdated `json:"PM2_5_Hysteresis,omitempty"`
	Monitoring            *DesiredMetadataUpdated `json:"Monitoring,omitempty"`
	MonitoringStop        *DesiredMetadataUpdated `json:"Monitoring_Stop,omitempty"`
	MonitoringStart       *DesiredMetadataUpdated `json:"Monitoring_Start,omitempty"`
	VmNoNIU               *DesiredMetadataUpdated `json:"VmNo_NIU,omitempty"`
}

type DesiredMetadataUpdated struct {
	LastUpdated        time.Time `json:"$lastUpdated"`
	LastUpdatedVersion int       `json:"$lastUpdatedVersion"`
}

type Reported struct {
	FrmVerNIU               *string              `json:"FrmVer_NIU,omitempty"`
	Workmode                string               `json:"Workmode"`
	FilterRFID              *string              `json:"FilterRFID,omitempty"`
	FilterLife              *int                 `json:"FilterLife,omitempty"`
	Fanspeed                int                  `json:"Fanspeed"`
	UILight                 bool                 `json:"UILight"`
	SafetyLock              bool                 `json:"SafetyLock"`
	Ionizer                 *bool                `json:"Ionizer,omitempty"`
	FilterType              *int                 `json:"FilterType,omitempty"`
	ErrPM25                 *bool                `json:"ErrPM2_5,omitempty"`
	ErrTVOC                 *bool                `json:"ErrTVOC,omitempty"`
	ErrTempHumidity         *bool                `json:"ErrTempHumidity,omitempty"`
	ErrFanMtr               *bool                `json:"ErrFanMtr,omitempty"`
	ErrCommSensorDisplayBrd *bool                `json:"ErrCommSensorDisplayBrd,omitempty"`
	DoorOpen                *bool                `json:"DoorOpen,omitempty"`
	ErrRFID                 *bool                `json:"ErrRFID,omitempty"`
	SignalStrength          string               `json:"SignalStrength"`
	LogE                    *int                 `json:"logE,omitempty"`
	LogW                    *int                 `json:"logW,omitempty"`
	InterfaceVersion        int                  `json:"InterfaceVer"`
	VmNoNIU                 string               `json:"VmNo_NIU"`
	TVOCBrand               *string              `json:"TVOCBrand,omitempty"`
	Capabilities            ReportedCapabilities `json:"capabilities"`
	Tasks                   interface{}          `json:"tasks"`
	Metadata                ReportedMetadata     `json:"$metadata"`
	Version                 int                  `json:"$version"`
	DeviceID                string               `json:"deviceId"`
	TVOC                    *int                 `json:"TVOC,omitempty"`
	CO2                     *int                 `json:"CO2,omitempty"` // Deprecated: use ECO2 (latest firmware).
	PM1                     *int                 `json:"PM1,omitempty"`
	PM25                    *int                 `json:"PM2_5,omitempty"`
	PM10                    *int                 `json:"PM10,omitempty"`
	Humidity                *int                 `json:"Humidity,omitempty"`
	Temp                    *int                 `json:"Temp,omitempty"`
	RSSI                    *int                 `json:"RSSI,omitempty"`
	ECO2                    *int                 `json:"ECO2,omitempty"`
	FilterLife1             *int                 `json:"FilterLife_1,omitempty"`
	Monitoring              *bool                `json:"Monitoring,omitempty"`
	MonitoringStop          *int                 `json:"Monitoring_Stop,omitempty"`
	MonitoringStart         *int                 `json:"Monitoring_Start,omitempty"`
	UVState                 *string              `json:"UVState,omitempty"`
	UVRuntime               *int                 `json:"UVRuntime,omitempty"`
	ErrCommSensorUIBrd      *string              `json:"ErrCommSensorUIBrd,omitempty"`
	ErrImpellerStuck        *string              `json:"ErrImpellerStuck,omitempty"`
	ErrPmNotResp            *string              `json:"ErrPmNotResp,omitempty"`
	VmNoMCU                 *string              `json:"VmNo_MCU,omitempty"`
	PM25Approximate         *int                 `json:"PM2_5_approximate,omitempty"`
}

type ReportedCapabilities struct {
	Tasks interface{} `json:"tasks,omitempty"`
}

type ReportedMetadata struct {
	ReportedMetadataUpdated
	FrmVerNIU               *ReportedMetadataUpdated     `json:"FrmVer_NIU,omitempty"`
	Workmode                ReportedMetadataUpdated      `json:"Workmode"`
	FilterRFID              *ReportedMetadataUpdated     `json:"FilterRFID,omitempty"`
	FilterLife              *ReportedMetadataUpdated     `json:"FilterLife,omitempty"`
	Fanspeed                ReportedMetadataUpdated      `json:"Fanspeed"`
	UILight                 ReportedMetadataUpdated      `json:"UILight"`
	SafetyLock              ReportedMetadataUpdated      `json:"SafetyLock"`
	Ionizer                 *ReportedMetadataUpdated     `json:"Ionizer,omitempty"`
	FilterType              *ReportedMetadataUpdated     `json:"FilterType,omitempty"`
	ErrPM25                 *ReportedMetadataUpdated     `json:"ErrPM2_5,omitempty"`
	ErrTVOC                 *ReportedMetadataUpdated     `json:"ErrTVOC,omitempty"`
	ErrTempHumidity         *ReportedMetadataUpdated     `json:"ErrTempHumidity,omitempty"`
	ErrFanMtr               *ReportedMetadataUpdated     `json:"ErrFanMtr,omitempty"`
	ErrCommSensorDisplayBrd *ReportedMetadataUpdated     `json:"ErrCommSensorDisplayBrd,omitempty"`
	DoorOpen                *ReportedMetadataUpdated     `json:"DoorOpen,omitempty"`
	ErrRFID                 *ReportedMetadataUpdated     `json:"ErrRFID,omitempty"`
	SignalStrength          ReportedMetadataUpdated      `json:"SignalStrength"`
	LogE                    *ReportedMetadataUpdated     `json:"logE,omitempty"`
	LogW                    *ReportedMetadataUpdated     `json:"logW,omitempty"`
	InterfaceVersion        ReportedMetadataUpdated      `json:"InterfaceVer"`
	VmNoNIU                 ReportedMetadataUpdated      `json:"VmNo_NIU"`
	TVOCBrand               *ReportedMetadataUpdated     `json:"TVOCBrand,omitempty"`
	Capabilities            ReportedMetadataCapabilities `json:"capabilities"`
	Tasks                   ReportedMetadataUpdated      `json:"tasks"`
	TVOC                    *ReportedMetadataUpdated     `json:"TVOC,omitempty"`
	CO2                     *ReportedMetadataUpdated     `json:"CO2,omitempty"`
	PM1                     *ReportedMetadataUpdated     `json:"PM1,omitempty"`
	PM25                    *ReportedMetadataUpdated     `json:"PM2_5,omitempty"`
	PM10                    *ReportedMetadataUpdated     `json:"PM10,omitempty"`
	Humidity                *ReportedMetadataUpdated     `json:"Humidity,omitempty"`
	Temp                    *ReportedMetadataUpdated     `json:"Temp,omitempty"`
	RSSI                    *ReportedMetadataUpdated     `json:"RSSI,omitempty"`
	ECO2                    *ReportedMetadataUpdated     `json:"ECO2,omitempty"`
	FilterLife1             *ReportedMetadataUpdated     `json:"FilterLife_1,omitempty"`
	Monitoring              *ReportedMetadataUpdated     `json:"Monitoring,omitempty"`
	MonitoringStop          *ReportedMetadataUpdated     `json:"Monitoring_Stop,omitempty"`
	MonitoringStart         *ReportedMetadataUpdated     `json:"Monitoring_Start,omitempty"`
	UVState                 *ReportedMetadataUpdated     `json:"UVState,omitempty"`
	UVRuntime               *ReportedMetadataUpdated     `json:"UVRuntime,omitempty"`
	ErrCommSensorUIBrd      *ReportedMetadataUpdated     `json:"ErrCommSensorUIBrd,omitempty"`
	ErrImpellerStuck        *ReportedMetadataUpdated     `json:"ErrImpellerStuck,omitempty"`
	ErrPmNotResp            *ReportedMetadataUpdated     `json:"ErrPmNotResp,omitempty"`
	VmNoMCU                 *ReportedMetadataUpdated     `json:"VmNo_MCU,omitempty"`
	PM25Approximate         *ReportedMetadataUpdated     `json:"PM2_5_approximate,omitempty"`
}

type ReportedMetadataCapabilities struct {
	ReportedMetadataUpdated
	Tasks *ReportedMetadataUpdated `json:"tasks,omitempty"`
}

type ReportedMetadataUpdated struct {
	LastUpdated time.Time `json:"$lastUpdated"`
}
