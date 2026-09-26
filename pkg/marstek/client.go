package marstek

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultPort         = 30000
	DefaultTimeout      = 10 * time.Second
	DefaultRequestDelay = 2 * time.Second
	BroadcastAddr       = "255.255.255.255"
)

type Client struct {
	addr         string
	port         int
	timeout      time.Duration
	retries      int
	requestDelay time.Duration
	lastRequest  time.Time
	conn         *net.UDPConn
	mu           sync.Mutex
	msgID        atomic.Int64
}

type Option func(*Client)

func WithPort(port int) Option {
	return func(c *Client) {
		c.port = port
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithRetries(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}

func WithRequestDelay(delay time.Duration) Option {
	return func(c *Client) {
		c.requestDelay = delay
	}
}

func NewClient(addr string, opts ...Option) *Client {
	c := &Client{
		addr:         addr,
		port:         DefaultPort,
		timeout:      DefaultTimeout,
		retries:      3,
		requestDelay: DefaultRequestDelay,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) nextID() int64 {
	return c.msgID.Add(1)
}

func (c *Client) send(req *Request) (*Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.lastRequest.IsZero() {
		elapsed := time.Since(c.lastRequest)
		if elapsed < c.requestDelay {
			time.Sleep(c.requestDelay - elapsed)
		}
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.addr, c.port))
	if err != nil {
		return nil, fmt.Errorf("resolve address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("dial udp: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(buf[:n], &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	c.lastRequest = time.Now()

	if resp.Error != nil {
		return &resp, &RPCError{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
			Data:    resp.Error.Data,
		}
	}

	return &resp, nil
}

func (c *Client) call(method string, params any, result any) error {
	req := &Request{
		ID:     c.nextID(),
		Method: method,
		Params: params,
	}

	var resp *Response
	var err error

	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*200) * time.Millisecond)
		}
		resp, err = c.send(req)
		if err == nil {
			break
		}
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			continue
		}
		return err
	}

	if err != nil {
		return err
	}

	if result != nil && resp.Result != nil {
		data, err := json.Marshal(resp.Result)
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		if err := json.Unmarshal(data, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}

func DiscoverDevices(port int, timeout time.Duration) ([]DeviceInfo, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", BroadcastAddr, port))
	if err != nil {
		return nil, fmt.Errorf("resolve broadcast address: %w", err)
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		return nil, fmt.Errorf("listen udp: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}

	req := &Request{
		ID:     0,
		Method: "Marstek.GetDevice",
		Params: map[string]string{"ble_mac": "0"},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	if _, err := conn.WriteToUDP(data, addr); err != nil {
		return nil, fmt.Errorf("write broadcast: %w", err)
	}

	var devices []DeviceInfo
	buf := make([]byte, 4096)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			return nil, fmt.Errorf("read: %w", err)
		}

		var resp Response
		if err := json.Unmarshal(buf[:n], &resp); err != nil {
			continue
		}

		if resp.Result != nil {
			data, _ := json.Marshal(resp.Result)
			var info DeviceInfo
			if err := json.Unmarshal(data, &info); err == nil {
				devices = append(devices, info)
			}
		}
	}

	return devices, nil
}

func (c *Client) GetDevice(bleMac string) (*DeviceInfo, error) {
	var result DeviceInfo
	params := map[string]string{"ble_mac": bleMac}
	if err := c.call("Marstek.GetDevice", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetWifiStatus(instanceID int) (*WifiStatus, error) {
	var result WifiStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("Wifi.GetStatus", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetBLEStatus(instanceID int) (*BLEStatus, error) {
	var result BLEStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("BLE.GetStatus", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetBatteryStatus(instanceID int) (*BatteryStatus, error) {
	var result BatteryStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("Bat.GetStatus", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPVStatus(instanceID int) (*PVStatus, error) {
	var result PVStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("PV.GetStatus", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetESStatus(instanceID int) (*ESStatus, error) {
	var result ESStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("ES.GetStatus", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetESMode(instanceID int) (*ESModeStatus, error) {
	var result ESModeStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("ES.GetMode", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetAutoMode(instanceID int) (*SetResult, error) {
	var result SetResult
	params := SetModeParams{
		ID: instanceID,
		Config: ModeConfig{
			Mode:    ModeAuto,
			AutoCfg: &AutoConfig{Enable: 1},
		},
	}
	if err := c.call("ES.SetMode", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetAIMode(instanceID int) (*SetResult, error) {
	var result SetResult
	params := SetModeParams{
		ID: instanceID,
		Config: ModeConfig{
			Mode:  ModeAI,
			AICfg: &AIConfig{Enable: 1},
		},
	}
	if err := c.call("ES.SetMode", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetManualMode(instanceID int, cfg ManualConfig) (*SetResult, error) {
	var result SetResult
	params := SetModeParams{
		ID: instanceID,
		Config: ModeConfig{
			Mode:      ModeManual,
			ManualCfg: &cfg,
		},
	}
	if err := c.call("ES.SetMode", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetPassiveMode(instanceID int, power int, cdTime int) (*SetResult, error) {
	var result SetResult
	params := SetModeParams{
		ID: instanceID,
		Config: ModeConfig{
			Mode: ModePassive,
			PassiveCfg: &PassiveConfig{
				Power:  power,
				CDTime: cdTime,
			},
		},
	}
	if err := c.call("ES.SetMode", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetUPSMode(instanceID int) (*SetResult, error) {
	var result SetResult
	params := SetModeParams{
		ID: instanceID,
		Config: ModeConfig{
			Mode:   ModeUPS,
			UPSCfg: &UPSConfig{Enable: 1},
		},
	}
	if err := c.call("ES.SetMode", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetEMStatus(instanceID int) (*EMStatus, error) {
	var result EMStatus
	params := map[string]int{"id": instanceID}
	if err := c.call("EM.GetStatus", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetDOD(value int) (*SetResult, error) {
	if value < 30 || value > 88 {
		return nil, fmt.Errorf("DOD value must be between 30 and 88, got %d", value)
	}
	var result SetResult
	params := map[string]int{"value": value}
	if err := c.call("DOD.SET", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetGridExportLimit(limit int) (*SetResult, error) {
	validLimits := []int{800, 1200, 1500, 2200, 2500}
	valid := false
	for _, v := range validLimits {
		if limit == v {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("grid export limit must be one of %v, got %d", validLimits, limit)
	}
	var result SetResult
	params := map[string]int{"version": limit}
	if err := c.call("Set.Ver", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetBLEAdvertising(enable bool) (*SetResult, error) {
	var result SetResult
	enableVal := 0
	if enable {
		enableVal = 1
	}
	params := map[string]int{"enable": enableVal}
	if err := c.call("Ble.Adv", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetLED(on bool) (*SetResult, error) {
	var result SetResult
	state := 0
	if on {
		state = 1
	}
	params := map[string]int{"state": state}
	if err := c.call("Led.Ctrl", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
