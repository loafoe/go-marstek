package marstek

import (
	"encoding/json"
	"testing"
)

func TestRequestMarshal(t *testing.T) {
	req := Request{
		ID:     1,
		Method: "ES.GetStatus",
		Params: map[string]int{"id": 0},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := `{"id":1,"method":"ES.GetStatus","params":{"id":0}}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestESStatusUnmarshal(t *testing.T) {
	data := `{
		"id": 0,
		"bat_soc": 46,
		"bat_cap": 5120,
		"pv_power": 0,
		"ongrid_power": -1487,
		"offgrid_power": 0,
		"total_pv_energy": 0,
		"total_grid_output_energy": 184456,
		"total_grid_input_energy": 209268,
		"total_load_energy": 0
	}`

	var status ESStatus
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if status.BatSOC == nil || *status.BatSOC != 46 {
		t.Errorf("expected BatSOC 46, got %v", status.BatSOC)
	}
	if status.OngridPower == nil || *status.OngridPower != -1487 {
		t.Errorf("expected OngridPower -1487, got %v", status.OngridPower)
	}
}

func TestManualConfigMarshal(t *testing.T) {
	cfg := ManualConfig{
		TimeNum:   1,
		StartTime: "08:30",
		EndTime:   "20:30",
		WeekSet:   WeekAll,
		Power:     100,
		Enable:    1,
	}

	params := SetModeParams{
		ID: 0,
		Config: ModeConfig{
			Mode:      ModeManual,
			ManualCfg: &cfg,
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal check failed: %v", err)
	}

	config := result["config"].(map[string]any)
	if config["mode"] != "Manual" {
		t.Errorf("expected mode Manual, got %v", config["mode"])
	}
}

func TestWeekDays(t *testing.T) {
	tests := []struct {
		days     []int
		expected int
	}{
		{[]int{WeekMonday}, 1},
		{[]int{WeekMonday, WeekTuesday}, 3},
		{[]int{WeekMonday, WeekWednesday, WeekFriday}, 21},
		{[]int{WeekAll}, 127},
	}

	for _, tt := range tests {
		result := WeekDays(tt.days...)
		if result != tt.expected {
			t.Errorf("WeekDays(%v) = %d, want %d", tt.days, result, tt.expected)
		}
	}
}

func TestRPCError(t *testing.T) {
	err := &RPCError{
		Code:    -32700,
		Message: "Parse error",
		Data:    403,
	}

	expected := "RPC error -32700: Parse error (data: 403)"
	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}

	err2 := &RPCError{
		Code:    -32601,
		Message: "Method not found",
	}

	expected2 := "RPC error -32601: Method not found"
	if err2.Error() != expected2 {
		t.Errorf("got %q, want %q", err2.Error(), expected2)
	}
}
