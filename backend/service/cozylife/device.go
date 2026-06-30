package cozylife

import "time"

const switchTypeCode = "00"

// ScanResult matches the JSON structure written by the cozylife CLI.
type ScanResult struct {
	IP             string `json:"ip"`
	DID            string `json:"did"`
	PID            string `json:"pid"`
	DMN            string `json:"dmn"`
	DPID           []int  `json:"dpid"`
	DeviceTypeCode string `json:"device_type_code"`
}

type device struct {
	cli *tcpClient
}

func newDeviceFromResult(r *ScanResult, timeout time.Duration) *device {
	return &device{cli: newClient(r.IP, timeout)}
}

func (d *device) turnOn() error  { return d.cli.control(map[string]any{"1": 1}) }
func (d *device) turnOff() error { return d.cli.control(map[string]any{"1": 0}) }

func (d *device) queryOn() (bool, error) {
	state, err := d.cli.query()
	if err != nil {
		return false, err
	}
	v, _ := state["1"].(float64)
	return v > 0, nil
}
