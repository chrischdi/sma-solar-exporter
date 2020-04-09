package inverter

import (
	"fmt"

	"k8s.io/klog"

	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
)

type Channel struct {
	Handle       int
	DeviceHandle int
	Name         string
	Unit         string
	Minimum      float64
	Maximum      float64
	Description  string
}

func NewChannel(deviceHandle, channelHandle int) (*Channel, error) {
	name, err := yasdi.GetChannelName(channelHandle)
	if err != nil {
		return nil, err
	}

	unit, err := yasdi.GetChannelUnit(channelHandle)
	if err != nil {
		return nil, err
	}

	min, max, err := yasdi.GetChannelValRange(channelHandle)
	if err != nil {
		if err.Error() != "there is no value range" {
			klog.Warningf("channel %s has no value range", name)
		}
		min = -1
		max = -1
	}

	channel := Channel{
		Handle:       channelHandle,
		DeviceHandle: deviceHandle,
		Name:         name,
		Unit:         unit,
		Minimum:      min,
		Maximum:      max,
		Description:  ChannelDescriptions[name],
	}
	return &channel, nil
}

func (c *Channel) GetValues() (float64, string, error) {
	return yasdi.GetChannelValue(c.Handle, c.DeviceHandle, 5)
}

func (c *Channel) String() string {
	return fmt.Sprintf("%4d | %15s | %5s | %58s", c.Handle, c.Name, c.Unit, c.Description)
}

// ChannelDescriptions is a mapping from channel names to their description.
// [0]: http://files.sma.de/dl/1348/NG_PAR-TB-en-22.pdf
// [1]: https://www.energymatters.com.au/images/SMA/SBNG_PAR-TEN084410.pdf
var ChannelDescriptions = map[string]string{
	"A.Ms.Amp":  "Input A DC current",
	"A.Ms.Vol":  "Input A DC voltage",
	"A.Ms.Watt": "Input A DC power",
	"B.Ms.Amp":  "Input B DC current",
	"B.Ms.Vol":  "Input B DC voltage",
	"B.Ms.Watt": "Input B DC power",

	"E-Total": "Total amount of feeding-in energy",

	"GridMs.PhV.phsA": "Grid voltage on phase L1",
	"GridMs.PhV.phsB": "Grid voltage on phase L2",
	"GridMs.PhV.phsC": "Grid voltage on phase L3",

	"Pac": "Delivered active power",

	"Cntry":           "Country Specific Standard",
	"Error":           "Error message",
	"GridMs.Hz":       "Grid frequency",
	"HP.swRev":        "Firmware version of the central assembly",
	"Inv.TmpLimStt":   "Status display for the derating due to excess temperatures",
	"MainModel":       "Inverter Device class",
	"Mlt.BatCha.Pwr":  "Minimum On power for MFR battery bank",
	"Mlt.BatCha.Tmm":  "Minimum time before reconnection of MFR battery bank",
	"Mlt.ComCtl.Sw":   "Status of MFR with control via communication",
	"Mlt.MinOnPwr":    "Minimum On power for MFR self-consumption",
	"Mlt.MinOnPwrTmm": "Minimum power On time, MFR self-consumption",
	"Mlt.MinOnTmm":    "Minimum On time for MFR self-consumption",
	"Mlt.OpMode":      "Operating mode of multifunction relay",
	"Mode":            "Operating condition",
	"Model":           "Inverter Device type",
	"MPPShdw.IsOn":    "OptiTrac Global Peak switched on",
	"Mt.TotOpTmh":     "Total number of grid-feeding operational hours",
	"Mt.TotTmh":       "Total number of operating hours of inverter",
	"Op.EvtCntUsr":    "Number of events for user",
	"Op.EvtNo":        "Current event number",
	"Op.GriSwStt":     "Grid relay status",
	"Op.OpModSet":     "Operating condition",
	"Op.TmsRmg":       "Waiting time until feed-in",
	"Pkg.swRev":       "Software version of components in inverter",
	"Plimit":          "Maximum active power device",
	"Serial Number":   "Serial number",
	"SerNumSet":       "Serial number",
	"SY_Systemzeit":   "System time as linux timestamp",
	"SY_Zeitzone":     "System timezone",

	"Wind_a0": "Power characteristic curves coefficient for Udc^0",
	"Wind_a1": "Power characteristic curves coefficient for Udc^1",
	"Wind_a2": "Power characteristic curves coefficient for Udc^2",
	"Wind_a3": "Power characteristic curves coefficient for Udc^3",
}
