package plant

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chrischdi/sma-solar-exporter/pkg/inverter"
)

type DummyInverter struct {
	name         string
	serialNumber uint
	values       map[string]float64
}

func (i *DummyInverter) GetChannelValue(channel string) (float64, error) {
	v, ok := i.values[channel]
	if !ok {
		return 0, fmt.Errorf("error")
	}
	return v, nil
}
func (i *DummyInverter) GetName() string       { return i.name }
func (i *DummyInverter) GetSerialNumber() uint { return i.serialNumber }
func (i *DummyInverter) GetChannelUnit(channelName string) (string, error) {
	return "unknown", nil
}

func init() {
	RegisterMetrics()
}

func TestNewPlant(t *testing.T) {
	tests := []struct {
		name     string
		channels []string
		want     *Plant
		wantErr  bool
	}{
		{
			"unsupported metric",
			[]string{"unsupported"},
			nil,
			true,
		},
		{
			"supported metric",
			[]string{"Pac"},
			&Plant{
				Channels: []string{"Pac"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPlant(nil, nil, tt.channels)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlant_Scrape(t *testing.T) {
	tests := []struct {
		name  string
		plant Plant
	}{
		{
			"happy path",
			Plant{
				Inverters: []inverter.ScrapeableInverter{
					&DummyInverter{
						name:         "foo",
						serialNumber: 1,
						values: map[string]float64{
							"Pac": 100,
						},
					},
				},
				Channels: []string{"Pac"},
			},
		},
		{
			"no metric",
			Plant{
				Inverters: []inverter.ScrapeableInverter{
					&DummyInverter{
						name:         "foo",
						serialNumber: 1,
						values:       map[string]float64{},
					},
				},
				Channels: []string{"doesNotExist"},
			},
		},
		{
			"error on getValue",
			Plant{
				Inverters: []inverter.ScrapeableInverter{
					&DummyInverter{
						name:         "foo",
						serialNumber: 1,
						values:       map[string]float64{},
					},
				},
				Channels: []string{"Pac"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.plant.Scrape()
		})
	}
}
