package marstek

import "fmt"

type Request struct {
	ID     int64  `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params"`
}

type Response struct {
	ID     int64          `json:"id"`
	Src    string         `json:"src,omitempty"`
	Result any            `json:"result,omitempty"`
	Error  *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RPCError struct {
	Code    int
	Message string
	Data    any
}

func (e *RPCError) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("RPC error %d: %s (data: %v)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

const (
	ErrParseError     = -32700
	ErrInvalidRequest = -32600
	ErrMethodNotFound = -32601
	ErrInvalidParams  = -32602
	ErrInternalError  = -32603
)

type DeviceInfo struct {
	Device   string `json:"device"`
	Version  int    `json:"ver"`
	BLEMAC   string `json:"ble_mac"`
	WifiMAC  string `json:"wifi_mac"`
	WifiName string `json:"wifi_name"`
	IP       string `json:"ip"`
}

type WifiStatus struct {
	ID      int     `json:"id"`
	WifiMAC string  `json:"wifi_mac"`
	SSID    *string `json:"ssid"`
	RSSI    int     `json:"rssi"`
	StaIP   *string `json:"sta_ip"`
	StaGate *string `json:"sta_gate"`
	StaMask *string `json:"sta_mask"`
	StaDNS  *string `json:"sta_dns"`
}

type BLEStatus struct {
	ID     int    `json:"id"`
	State  string `json:"state"`
	BLEMAC string `json:"ble_mac"`
}

type BatteryStatus struct {
	ID            int      `json:"id"`
	SOC           int      `json:"soc"`
	ChargFlag     bool     `json:"charg_flag"`
	DischrgFlag   bool     `json:"dischrg_flag"`
	BatTemp       *float64 `json:"bat_temp"`
	BatCapacity   *float64 `json:"bat_capacity"`
	RatedCapacity *float64 `json:"rated_capacity"`
}

type PVStatus struct {
	ID        int     `json:"id"`
	PVPower   float64 `json:"pv_power"`
	PVVoltage float64 `json:"pv_voltage"`
	PVCurrent float64 `json:"pv_current"`
	PVState   int     `json:"PV_state"`

	PV1Power   *float64 `json:"pv1_power,omitempty"`
	PV1Voltage *float64 `json:"pv1_voltage,omitempty"`
	PV1Current *float64 `json:"pv1_current,omitempty"`
	PV1State   *int     `json:"pv1_state,omitempty"`

	PV2Power   *float64 `json:"pv2_power,omitempty"`
	PV2Voltage *float64 `json:"pv2_voltage,omitempty"`
	PV2Current *float64 `json:"pv2_current,omitempty"`
	PV2State   *int     `json:"pv2_state,omitempty"`

	PV3Power   *float64 `json:"pv3_power,omitempty"`
	PV3Voltage *float64 `json:"pv3_voltage,omitempty"`
	PV3Current *float64 `json:"pv3_current,omitempty"`
	PV3State   *int     `json:"pv3_state,omitempty"`

	PV4Power   *float64 `json:"pv4_power,omitempty"`
	PV4Voltage *float64 `json:"pv4_voltage,omitempty"`
	PV4Current *float64 `json:"pv4_current,omitempty"`
	PV4State   *int     `json:"pv4_state,omitempty"`
}

type ESStatus struct {
	ID                   *int     `json:"id,omitempty"`
	BatSOC               *int     `json:"bat_soc,omitempty"`
	BatCap               *float64 `json:"bat_cap,omitempty"`
	PVPower              *float64 `json:"pv_power,omitempty"`
	OngridPower          *float64 `json:"ongrid_power,omitempty"`
	OffgridPower         *float64 `json:"offgrid_power,omitempty"`
	BatPower             *float64 `json:"bat_power,omitempty"`
	TotalPVEnergy        *float64 `json:"total_pv_energy,omitempty"`
	TotalGridOutputEnergy *float64 `json:"total_grid_output_energy,omitempty"`
	TotalGridInputEnergy  *float64 `json:"total_grid_input_energy,omitempty"`
	TotalLoadEnergy      *float64 `json:"total_load_energy,omitempty"`
}

type ESModeStatus struct {
	ID           *int     `json:"id,omitempty"`
	Mode         string   `json:"mode"`
	OngridPower  *float64 `json:"ongrid_power,omitempty"`
	OffgridPower *float64 `json:"offgrid_power,omitempty"`
	BatSOC       *int     `json:"bat_soc,omitempty"`
	CTState      *int     `json:"ct_state,omitempty"`
	APower       *float64 `json:"a_power,omitempty"`
	BPower       *float64 `json:"b_power,omitempty"`
	CPower       *float64 `json:"c_power,omitempty"`
	TotalPower   *float64 `json:"total_power,omitempty"`
	InputEnergy  *float64 `json:"input_energy,omitempty"`
	OutputEnergy *float64 `json:"output_energy,omitempty"`
}

type EMStatus struct {
	ID           int      `json:"id"`
	CTState      *int     `json:"ct_state,omitempty"`
	APower       *float64 `json:"a_power,omitempty"`
	BPower       *float64 `json:"b_power,omitempty"`
	CPower       *float64 `json:"c_power,omitempty"`
	TotalPower   *float64 `json:"total_power,omitempty"`
	InputEnergy  *float64 `json:"input_energy,omitempty"`
	OutputEnergy *float64 `json:"output_energy,omitempty"`
}

const (
	ModeAuto    = "Auto"
	ModeAI      = "AI"
	ModeManual  = "Manual"
	ModePassive = "Passive"
	ModeUPS     = "Ups"
)

type SetModeParams struct {
	ID     int        `json:"id"`
	Config ModeConfig `json:"config"`
}

type ModeConfig struct {
	Mode       string         `json:"mode"`
	AutoCfg    *AutoConfig    `json:"auto_cfg,omitempty"`
	AICfg      *AIConfig      `json:"ai_cfg,omitempty"`
	ManualCfg  *ManualConfig  `json:"manual_cfg,omitempty"`
	PassiveCfg *PassiveConfig `json:"passive_cfg,omitempty"`
	UPSCfg     *UPSConfig     `json:"ups_cfg,omitempty"`
}

type AutoConfig struct {
	Enable int `json:"enable"`
}

type AIConfig struct {
	Enable int `json:"enable"`
}

type ManualConfig struct {
	TimeNum   int    `json:"time_num"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	WeekSet   int    `json:"week_set"`
	Power     int    `json:"power"`
	Enable    int    `json:"enable"`
}

type PassiveConfig struct {
	Power  int `json:"power"`
	CDTime int `json:"cd_time"`
}

type UPSConfig struct {
	Enable int `json:"enable"`
}

type SetResult struct {
	ID        int  `json:"id"`
	SetResult bool `json:"set_result"`
}

const (
	WeekMonday    = 1 << 0
	WeekTuesday   = 1 << 1
	WeekWednesday = 1 << 2
	WeekThursday  = 1 << 3
	WeekFriday    = 1 << 4
	WeekSaturday  = 1 << 5
	WeekSunday    = 1 << 6
	WeekAll       = WeekMonday | WeekTuesday | WeekWednesday | WeekThursday | WeekFriday | WeekSaturday | WeekSunday
)

func WeekDays(days ...int) int {
	var result int
	for _, d := range days {
		result |= d
	}
	return result
}
