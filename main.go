package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

func san() {
	fmt.Println("Collecting SAN data for:", client)
	getSwitches("bdcsan.json")
	// fmt.Println(switches)

	arraysLength := len(switches)
	wg.Add(arraysLength)
	for _, switche := range switches {
		go collectDataSan(switche)
	}
	wg.Wait()
}

func storage() {
	fmt.Println("Collecting Storage data for:", client)
	getArrays("IBM.json", "ibm", "lsmdiskgrp -bytes -delim ,", "lssystem -delim ,| grep -i code", "", "")
	getArrays("3par.json", "3par", "", "", "", "")
	getArrays("huawei.json", "huawei", "show storage_pool general", "show system general", "", "")
	getArrays("dell.json", "dell", "show volume-maps", "show versions", "show volumes", "show volume-maps")

	arraysLength := len(arrays)
	wg.Add(arraysLength)
	for _, array := range arrays {
		go collectDataStorage(array)
	}
	wg.Wait()

	telia.Name = client
	telia.P16Total = 0
	telia.P16Free = 0
	telia.P16InternalTotal = 0
	telia.P16InternalFree = 0
	telia.P16InternalSSDTotal = 0
	telia.P16InternalHDDTotal = 0
	telia.P16InternalSSDFree = 0
	telia.P16InternalHDDFree = 0
	telia.P16InternalSSDMinLun = 0
	telia.P16InternalHDDMinLun = 0
	telia.P16ExternalTotal = 0
	telia.P16ExternalFree = 0
	telia.P16ExternalSSDTotal = 0
	telia.P16ExternalHDDTotal = 0
	telia.P16ExternalSSDFree = 0
	telia.P16ExternalHDDFree = 0
	telia.P16ExternalSSDMinLun = 0
	telia.P16ExternalHDDMinLun = 0
	telia.Z141Total = 0
	telia.Z141Free = 0
	telia.Z141InternalTotal = 0
	telia.Z141InternalFree = 0
	telia.Z141InternalSSDTotal = 0
	telia.Z141InternalHDDTotal = 0
	telia.Z141InternalSSDFree = 0
	telia.Z141InternalHDDFree = 0
	telia.Z141InternalSSDMinLun = 0
	telia.Z141InternalHDDMinLun = 0
	telia.Z141ExternalTotal = 0
	telia.Z141ExternalFree = 0
	telia.Z141ExternalSSDTotal = 0
	telia.Z141ExternalHDDTotal = 0
	telia.Z141ExternalSSDFree = 0
	telia.Z141ExternalHDDFree = 0
	telia.Z141ExternalSSDMinLun = 0
	telia.Z141ExternalHDDMinLun = 0
	telia.Total = 0
	telia.TotalFree = 0

	poolsLength := len(pools.Pools)
	wg.Add(poolsLength)
	fmt.Println("Finished for loop")
	for _, pool := range pools.Pools {
		go updateStoragePools(pool)
	}
	wg.Wait()
	updateStorageClient()
}

var pools Pools
var arrays []Array
var switches []Switch
var ts = fmt.Sprint(time.Now().UnixNano())
var wg sync.WaitGroup
var client string
var telia Client
var switchPortMetrics = "switch_ports"

func main() {
	var option string

	flag.StringVar(&option, "o", "", "option to collect \"san\" or \"storage\"")
	flag.StringVar(&client, "c", "", "client name to collect data for")
	flag.Parse()

	logError("Start")
	if option == "san" {
		san()
	} else if option == "storage" {
		storage()
	} else {
		fmt.Println("Such option is not available")
	}
	logError("Finish")
}
