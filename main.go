package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/safchain/ethtool"
)

var netInterfaces = net.Interfaces

type ethtoolWrapper interface {
	Stats(s string) (map[string]uint64, error)
}

func getInterfacesList(ifaceRegexp string) ([]net.Interface, error) {
	// Retrieve network interface list
	var ifaceList []net.Interface
	dat, err := netInterfaces()
	if err == nil {
		ifaceregexp, _ := regexp.Compile(ifaceRegexp)
		for _, iface := range dat {
			if ifaceregexp.MatchString(iface.Name) {
				ifaceList = append(ifaceList, iface)
			}
		}
		log.Printf("Got %d interfaces\n", len(ifaceList))
	}
	return ifaceList, nil
}

func writeMetrics(output string, metrics []byte) {
	log.Printf("writing metrics in %s", output)
	err := ioutil.WriteFile(output, metrics, 0644)
	if err != nil {
		log.Printf("writeMetrics failed with error %s", err.Error())
	}
}

func retrieveStats(ethandler ethtoolWrapper, iface net.Interface, tick int) ([]byte, error) {
	var metrics = fmt.Sprintf("ethtool_tick %d\n", tick)
	var labels = fmt.Sprintf("{device=%q}", iface.Name)

	stats, err := ethandler.Stats(iface.Name)
	if err != nil {
		return []byte(metrics), err
	}
	for metric, value := range stats {
		metrics += fmt.Sprintf(
			"ethtool_%s%s %d\n",
			strings.Replace(metric, "-", "_", -1),
			strings.Replace(labels, "-", "_", -1),
			value)
	}

	return []byte(metrics), nil
}

func mainLoop(ethandler *ethtool.Ethtool, ifaceList []net.Interface, output *string, sleep *int) {
	var tick = 0
	for true {
		// Increment loop's tick
		tick++

		// Iterate over network interfaces' list to retrieve statistics
		for _, iface := range ifaceList {
			log.Printf("Processing interface %s\n", iface.Name)

			// Retrieve stats from network interface and build metrics list
			metrics, err := retrieveStats(ethandler, iface, tick)
			if err != nil {
				if err.Error() == "operation not supported" {
					log.Printf("You should consider using [-ifaceregexp] option to avoid processing interface [%s]", iface.Name)
				}
				log.Printf("Got error [%s] for interface [%s]", err, iface.Name)
				continue
			}

			// Write metrics
			writeMetrics(*output, metrics)
		}
		log.Printf("End tick %d", tick)

		// Wait `sleep` second before next iteration
		time.Sleep(time.Duration(*sleep) * time.Second)
	}
}

func main() {
	// Parse commandline options
	ifaceRegexp := flag.String("ifaceregexp", ".*", "an interface name or regexp")
	sleep := flag.Int("sleep", 20, "time in second to wait between two statistics gathering")
	output := flag.String("output", "/prom_output/ethtool.prom", "an existing directory to store file containing metrics")
	flag.Parse()

	// Retrieve network interfaces list
	ifaceList, err := getInterfacesList(*ifaceRegexp)
	if err != nil {
		panic("Unable to retrieve network interfaces list with error: " + err.Error())
	}
	if len(ifaceList) == 0 {
		panic("No network interface found. Exiting.")
	}

	// Initializing ethtool
	ethandler, err := ethtool.NewEthtool()
	if err != nil {
		panic("Unable to initialize ethtool with error: " + err.Error())
	}

	// Initializing main loop
	mainLoop(ethandler, ifaceList, output, sleep)

	// Exiting
	log.Println("Exit")
}
