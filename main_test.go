package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"testing"
)

/* Test for method getInterfacesList
 */
func mockNetInterfaces() ([]net.Interface, error) {
	ifaceList := []net.Interface{
		net.Interface{Name: "lo"},
		net.Interface{Name: "eth0"},
		net.Interface{Name: "eth1"},
		net.Interface{Name: "bond0"},
	}
	return ifaceList, nil
}

func TestGetInterfacesList(t *testing.T) {
	// Suppress logs
	out := new(bytes.Buffer)
	writer := io.Writer(out)
	log.SetOutput(writer)
	// Test various regexp possibilities, checking  we have the right number of interfaces in return
	netInterfaces = mockNetInterfaces
	usecases := map[string]int{
		".*":     4,
		"eth.*":  2,
		"[0-9]+": 3,
		"ens.*":  0,
	}
	for regexp, want := range usecases {
		ifaceList, _ := getInterfacesList(regexp)
		if len(ifaceList) != want {
			t.Errorf("Incorrect number of detected interfaces , got: %d, want: %d.", len(ifaceList), want)
		}
	}
}

/* Test for method retrieveStats
 */
type mockEthtool struct{}

func (e *mockEthtool) Stats(ifaceName string) (map[string]uint64, error) {
	var dat = make(map[string]uint64)
	switch ifaceName {
	case "lo":
		return dat, errors.New("operation not supported")
	case "eth0":
		dat["stat1"] = 1
		return dat, nil
	default:
		return dat, errors.New("no such device")
	}
}

func TestRetrieveStats(t *testing.T) {
	ethtool := mockEthtool{}
	ifaceList := map[string]error{
		"lo":    errors.New("operation not supported"), // Interface exists but doesn't support ethtool stats
		"eth0":  nil,                                   // Interface exists and supports ethtool stats
		"bond0": errors.New("no such device"),          // Interface dosn't exists
	}
	for ifaceName, wantedResult := range ifaceList {
		iface := net.Interface{
			Name: ifaceName,
		}
		stats, err := retrieveStats(&ethtool, iface, 1)
		if ifaceName == "lo" && err.Error() != wantedResult.Error() {
			t.Errorf("Should get an 'operation not supported' error for interface %s", ifaceName)
		}
		if ifaceName == "bond0" && err.Error() != wantedResult.Error() {
			t.Errorf("Should get a 'no such device' error for interface %s", ifaceName)
		}
		if ifaceName == "eth0" && err == wantedResult && len(stats) != 46 {
			t.Errorf("Incorrect stats' number for interfaces %s, got: %d, want: %d.", ifaceName, len(stats), 1)
		}
	}

}
