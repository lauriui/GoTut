package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Users struct which contains
// an array of users

func logError(Error string) {
	log_date := time.Now()
	years, month, day := log_date.Date()
	filename := "logs/" + client + "." + strconv.Itoa(years) + strconv.Itoa(int(month)) + strconv.Itoa(day) + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println(Error)

}

func collectDataSan(array Switch) {
	defer wg.Done()
	var connection *ssh.Client
	var err error
	connection, err = connectToHostPW(username, password, array.Ip)
	if err != nil {
		connection, err = connectToHostKB(username, password, array.Ip)
		if err != nil {
			errorString := "CollectData: ConnectToHostKB: " + array.Name + ": " + err.Error()
			logError(errorString)
			return
		}
	}
	data, err := runCommand(connection, "switchshow")
	if err != nil {
		logError(err.Error())
	}
	nameData, err := runCommand(connection, "switchshow -portname")
	if err != nil {
		logError(err.Error())
	}
	fwData, err := runCommand(connection, "version")
	if err != nil {
		logError(err.Error())
	}
	defer connection.Close()
	parseDataSan(data, nameData, fwData, array)
	// var poolData Pools
	// poolData, err = parseData(data, fw, array.Model, array.Name, array.Site, array.Type, client)
	// if err != nil {
	// 	logError(err.Error())
	// }
	// pools.Pools = append(pools.Pools, poolData.Pools...)
}

func collectDataStorage(array Array) {
	defer wg.Done()
	if array.Client != client {
		return
	}
	var connection *ssh.Client
	var err error
	connection, err = connectToHostPW(username, password, array.Ip)
	if err != nil {
		connection, err = connectToHostKB(username, password, array.Ip)
		if err != nil {
			errorString := "CollectData: ConnectToHostKB: " + array.Name + ": " + err.Error()
			logError(errorString)
			return
		}
	}
	data, err := runCommand(connection, array.Data)
	if err != nil {
		logError(err.Error())
	}
	fw, err := runCommand(connection, array.fw)
	if err != nil {
		logError(err.Error())
	}
	defer connection.Close()
	var poolData Pools
	poolData, err = parseData(data, fw, array.Model, array.Name, array.Site, array.Type, client)
	if err != nil {
		logError(err.Error())
	}
	pools.Pools = append(pools.Pools, poolData.Pools...)
}

func updateStoragePools(pool Pool) {
	defer wg.Done()
	// fmt.Println(pool.PoolName, " ", pool.ArrayName)

	telia.Total += pool.PoolCapacity
	telia.TotalFree += pool.PoolCapacityFree
	if pool.Site == "P16" {
		telia.P16Total += pool.PoolCapacity
		telia.P16Free += pool.PoolCapacityFree
		if strings.Contains(pool.Type, "Internal") {
			telia.P16InternalTotal += pool.PoolCapacity
			telia.P16InternalFree += pool.PoolCapacityFree
			if strings.Contains(pool.Type, "SSD") || strings.Contains(pool.Type, "MIX") {
				telia.P16InternalSSDTotal += pool.PoolCapacity
				telia.P16InternalSSDFree += pool.PoolCapacityFree
				telia.P16InternalSSDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			} else if strings.Contains(pool.Type, "SAS") {
				telia.P16InternalHDDTotal += pool.PoolCapacity
				telia.P16InternalHDDFree += pool.PoolCapacityFree
				telia.P16InternalHDDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			}
		} else if strings.Contains(pool.Type, "Shared") {
			telia.P16ExternalTotal += pool.PoolCapacity
			telia.P16ExternalFree += pool.PoolCapacityFree
			if strings.Contains(pool.Type, "SSD") || strings.Contains(pool.Type, "MIX") {
				telia.P16ExternalSSDTotal += pool.PoolCapacity
				telia.P16ExternalSSDFree += pool.PoolCapacityFree
				telia.P16ExternalSSDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			} else if strings.Contains(pool.Type, "SAS") {
				telia.P16ExternalHDDTotal += pool.PoolCapacity
				telia.P16ExternalHDDFree += pool.PoolCapacityFree
				telia.P16ExternalHDDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			}
		}
	} else if pool.Site == "Z141" {
		telia.Z141Total += pool.PoolCapacity
		telia.Z141Free += pool.PoolCapacityFree
		if strings.Contains(pool.Type, "Internal") {
			telia.Z141InternalTotal += pool.PoolCapacity
			telia.Z141InternalFree += pool.PoolCapacityFree
			if strings.Contains(pool.Type, "SSD") || strings.Contains(pool.Type, "MIX") {
				telia.Z141InternalSSDTotal += pool.PoolCapacity
				telia.Z141InternalSSDFree += pool.PoolCapacityFree
				telia.Z141InternalSSDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			} else if strings.Contains(pool.Type, "SAS") {
				telia.Z141InternalHDDTotal += pool.PoolCapacity
				telia.Z141InternalHDDFree += pool.PoolCapacityFree
				telia.Z141InternalHDDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			}
		} else if strings.Contains(pool.Type, "Shared") {
			telia.Z141ExternalTotal += pool.PoolCapacity
			telia.Z141ExternalFree += pool.PoolCapacityFree
			if strings.Contains(pool.Type, "SSD") || strings.Contains(pool.Type, "MIX") {
				telia.Z141ExternalSSDTotal += pool.PoolCapacity
				telia.Z141ExternalSSDFree += pool.PoolCapacityFree
				telia.Z141ExternalSSDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			} else if strings.Contains(pool.Type, "SAS") {
				telia.Z141ExternalHDDTotal += pool.PoolCapacity
				telia.Z141ExternalHDDFree += pool.PoolCapacityFree
				telia.Z141ExternalHDDMinLun += int(pool.PoolCapacityFree / 10000000000000)
			}
		}
	} else if pool.Site == "Stretched" {

		if strings.Contains(pool.PoolName, "P16") {
			telia.P16Total += pool.PoolCapacity
			telia.P16Free += pool.PoolCapacityFree
			telia.StretchedP16Total += pool.PoolCapacity
			telia.StretchedP16Free += pool.PoolCapacityFree
			telia.StretchedP16MinLun += int(pool.PoolCapacityFree / 10000000000000)

		} else if strings.Contains(pool.PoolName, "Z141") {
			telia.Z141Total += pool.PoolCapacity
			telia.Z141Free += pool.PoolCapacityFree
			telia.StretchedZ141Total += pool.PoolCapacity
			telia.StretchedZ141Free += pool.PoolCapacityFree
			telia.StretchedZ141MinLun += int(pool.PoolCapacityFree / 10000000000000)
		}
	}

	// {
	// 	for _, array := range arrayVolume {
	// 		for _, vol := range array.Vols {
	// 			clientString := "BackupVolumes,client=\"" + vol.Client + "\",Array=\"" + vol.ArrayName + "\",ID=\"" + vol.ArrayName + vol.Id + "\",Site=\"" + vol.Site + "\" Total=" + fmt.Sprintf("%f", vol.VolTotalSize) + ",TotalUsed=" + fmt.Sprintf("%f", vol.VolAllocatedSize) + ",Type=\"" + vol.Type + "\",Name=\"" + vol.VolName + "\",Hosts=\"" + strings.Join(vol.Hosts, ",") + "\" " + ts
	// 			resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(clientString)))
	// 			if err != nil {
	// 				log.Fatalln(err)
	// 			}
	// 			defer resp.Body.Close()
	// 			fmt.Println(clientString)
	// 			var res map[string]interface{}
	// 			json.NewDecoder(resp.Body).Decode(&res)
	// 			fmt.Println(fmt.Sprint(resp.StatusCode))
	// 		}
	// 	}

	// 	for _, pool := range pools.Pools {
	var clientData string
	if client == "Telia" {
		clientData = "teliaData"
	}
	{
		totalString := clientData + ",ID=\"" + pool.Id + pool.ArrayName +
			",Site=" + pool.Site +
			",type=" + pool.Type +
			" Array=\"" + pool.ArrayName +
			"\",Firmware=\"" + strings.ReplaceAll(strings.ReplaceAll(pool.Firmware, " ", ""), ",", "") +
			"\",Pool=\"" + pool.PoolName +
			"\",TotalCapacity=" + fmt.Sprintf("%f", pool.PoolCapacity) +
			",FreeCapacity=" + fmt.Sprintf("%f", pool.PoolCapacityFree) +
			",UsedCapacity=" + fmt.Sprintf("%f", pool.PoolCapacityUsed) +
			",AllocationPCT=" + fmt.Sprintf("%f", pool.PoolCapacityPCT) +
			" " + ts
		resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(totalString)))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		fmt.Println(totalString)
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		fmt.Println(fmt.Sprint(resp.StatusCode))
	}

}

func updateStorageClient() {
	{
		clientString := "clientData,client=\"" + telia.Name +
			"\" Total=" + fmt.Sprintf("%f", telia.Total) +
			",TotalFree=" + fmt.Sprintf("%f", telia.TotalFree) +
			",P16Total=" + fmt.Sprintf("%f", telia.P16Total) +
			",P16Free=" + fmt.Sprintf("%f", telia.P16Free) +
			",Z141Total=" + fmt.Sprintf("%f", telia.Z141Total) +
			",Z141Free=" + fmt.Sprintf("%f", telia.Z141Free) +
			",P16InternalTotal=" + fmt.Sprintf("%f", telia.P16InternalTotal) +
			",P16InternalFree=" + fmt.Sprintf("%f", telia.P16InternalFree) +
			",P16InternalSSDTotal=" + fmt.Sprintf("%f", telia.P16InternalSSDTotal) +
			",P16InternalHDDTotal=" + fmt.Sprintf("%f", telia.P16InternalHDDTotal) +
			",P16InternalSSDFree=" + fmt.Sprintf("%f", telia.P16InternalSSDFree) +
			",P16InternalHDDFree=" + fmt.Sprintf("%f", telia.P16InternalHDDFree) +
			",P16InternalSSDMinLun=" + fmt.Sprint(telia.P16InternalSSDMinLun) +
			",P16InternalHDDMinLun=" + fmt.Sprint(telia.P16InternalHDDMinLun) +
			",P16ExternalTotal=" + fmt.Sprintf("%f", telia.P16ExternalTotal) +
			",P16ExternalFree=" + fmt.Sprintf("%f", telia.P16ExternalFree) +
			",P16ExternalSSDTotal=" + fmt.Sprintf("%f", telia.P16ExternalSSDTotal) +
			",P16ExternalHDDTotal=" + fmt.Sprintf("%f", telia.P16ExternalHDDTotal) +
			",P16ExternalSSDFree=" + fmt.Sprintf("%f", telia.P16ExternalSSDFree) +
			",P16ExternalHDDFree=" + fmt.Sprintf("%f", telia.P16ExternalHDDFree) +
			",P16ExternalSSDMinLun=" + fmt.Sprint(telia.P16ExternalSSDMinLun) +
			",P16ExternalHDDMinLun=" + fmt.Sprint(telia.P16ExternalHDDMinLun) +
			",Z141InternalTotal=" + fmt.Sprintf("%f", telia.Z141InternalTotal) +
			",Z141InternalFree=" + fmt.Sprintf("%f", telia.Z141InternalFree) +
			",Z141InternalSSDTotal=" + fmt.Sprintf("%f", telia.Z141InternalSSDTotal) +
			",Z141InternalHDDTotal=" + fmt.Sprintf("%f", telia.Z141InternalHDDTotal) +
			",Z141InternalSSDFree=" + fmt.Sprintf("%f", telia.Z141InternalSSDFree) +
			",Z141InternalHDDFree=" + fmt.Sprintf("%f", telia.Z141InternalHDDFree) +
			",Z141InternalSSDMinLun=" + fmt.Sprint(telia.Z141InternalSSDMinLun) +
			",Z141InternalHDDMinLun=" + fmt.Sprint(telia.Z141InternalHDDMinLun) +
			",Z141ExternalTotal=" + fmt.Sprintf("%f", telia.Z141ExternalTotal) +
			",Z141ExternalFree=" + fmt.Sprintf("%f", telia.Z141ExternalFree) +
			",Z141ExternalSSDTotal=" + fmt.Sprintf("%f", telia.Z141ExternalSSDTotal) +
			",Z141ExternalHDDTotal=" + fmt.Sprintf("%f", telia.Z141ExternalHDDTotal) +
			",Z141ExternalSSDFree=" + fmt.Sprintf("%f", telia.Z141ExternalSSDFree) +
			",Z141ExternalHDDFree=" + fmt.Sprintf("%f", telia.Z141ExternalHDDFree) +
			",Z141ExternalSSDMinLun=" + fmt.Sprint(telia.Z141ExternalSSDMinLun) +
			",Z141ExternalHDDMinLun=" + fmt.Sprint(telia.Z141ExternalHDDMinLun) +
			",StretchedP16Total=" + fmt.Sprint(telia.StretchedP16Total) +
			",StretchedP16Free=" + fmt.Sprint(telia.StretchedP16Free) +
			",StretchedP16MinLun=" + fmt.Sprint(telia.StretchedP16MinLun) +
			",StretchedZ141Total=" + fmt.Sprint(telia.StretchedZ141Total) +
			",StretchedZ141Free=" + fmt.Sprint(telia.StretchedZ141Free) +
			",StretchedZ141MinLun=" + fmt.Sprint(telia.StretchedZ141MinLun) +
			" " + ts
		resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(clientString)))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		fmt.Println(clientString)
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		fmt.Println(fmt.Sprint(resp.StatusCode))
	}
}

func runCommand(connection *ssh.Client, command string) ([]byte, error) {
	var output []byte
	var err error
	session, err := connection.NewSession()
	if err != nil {
		logError(err.Error())
	}
	output, err = session.CombinedOutput(command)

	return output, err

}
func parseDataSan(switchData, nameData, fwData []byte, device Switch) (output Port, err error) {
	for _, line := range strings.Split(string(switchData), "\n") {
		fmt.Println(string(line))
	}

	return output, err
}

func parseData(inputData []byte, inputFw []byte, model, array, site, type_s, client_s string) (output Pools, err error) {
	splitInputData := strings.Split(string(inputData), "\n")
	switch model {
	case "ibm":
		firmware := strings.Split(strings.Split(string(inputFw), ",")[1], " ")[0]
		for i := 1; i < len(splitInputData); i++ {
			lineSplit := strings.Split(splitInputData[i], ",")
			if len(lineSplit) > 1 && lineSplit[0] != "--" {
				var pool Pool
				pool.ArrayName = array
				pool.Id = strings.ReplaceAll(lineSplit[0], "	", "")
				pool.PoolName = lineSplit[1]
				pool.PoolCapacity, err = strconv.ParseFloat(lineSplit[5], 64)
				pool.PoolCapacityUsed, err = strconv.ParseFloat(lineSplit[8], 64)
				pool.PoolCapacityFree, err = strconv.ParseFloat(lineSplit[7], 64)
				pool.PoolCapacityPCT = pool.PoolCapacityUsed / pool.PoolCapacity
				pool.Firmware = firmware
				pool.Site = site
				pool.Type = type_s
				pool.Client = client_s
				output.Pools = append(output.Pools, pool)
			}
		}

	case "huawei":

		splitFW := strings.Split(string(inputFw), "\n")
		var pversion, patch, firmware string
		for index, line := range splitFW {
			if strings.Contains(line, "Product Version") {
				pversion = strings.ReplaceAll(strings.Split(line, ":")[1], " ", "")
			}
			if strings.Contains(line, "Patch Version") {
				patch = strings.ReplaceAll(strings.Split(line, ":")[1], " ", "")
			}
			index = index
		}
		firmware = pversion + ", " + patch
		// fmt.Println(firmware)
		for i := 3; i < len(splitInputData); i++ {

			line := strings.ReplaceAll(splitInputData[i], "	", "")
			re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
			re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
			line = re_leadclose_whtsp.ReplaceAllString(line, "")
			line = re_inside_whtsp.ReplaceAllString(line, " ")
			splitLine := strings.Split(line, " ")
			if len(splitLine) > 1 && splitLine[0] != "--" {
				var pool Pool
				pool.ArrayName = array
				pool.Id = splitLine[0]
				pool.PoolName = splitLine[1]
				var tcap float64
				if strings.Contains(splitLine[5], "PB") {
					tcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[5], "PB", ""), 64)
					tcap = tcap * 1024 * 1024 * 1024 * 1024 * 1024
				} else if strings.Contains(splitLine[5], "TB") {
					tcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[5], "TB", ""), 64)
					tcap = tcap * 1024 * 1024 * 1024 * 1024
				} else if strings.Contains(splitLine[5], "GB") {
					tcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[5], "GB", ""), 64)
					tcap = tcap * 1024 * 1024 * 1024
				} else if strings.Contains(splitLine[5], "MB") {
					tcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[5], "MB", ""), 64)
					tcap = tcap * 1024 * 1024
				} else if strings.Contains(splitLine[5], "KB") {
					tcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[5], "KB", ""), 64)
					tcap = tcap * 1024
				} else if strings.Contains(splitLine[5], "B") {
					tcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[5], "B", ""), 64)
				}
				pool.PoolCapacity = tcap
				var fcap float64
				if strings.Contains(splitLine[6], "PB") {
					fcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[6], "PB", ""), 64)
					fcap = fcap * 1024 * 1024 * 1024 * 1024 * 1024
				} else if strings.Contains(splitLine[6], "TB") {
					fcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[6], "TB", ""), 64)
					fcap = fcap * 1024 * 1024 * 1024 * 1024
				} else if strings.Contains(splitLine[6], "GB") {
					fcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[6], "GB", ""), 64)
					fcap = fcap * 1024 * 1024 * 1024
				} else if strings.Contains(splitLine[6], "MB") {
					fcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[6], "MB", ""), 64)
					fcap = fcap * 1024 * 1024
				} else if strings.Contains(splitLine[6], "KB") {
					fcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[6], "KB", ""), 64)
					fcap = fcap * 1024
				} else if strings.Contains(splitLine[6], "B") {
					fcap, err = strconv.ParseFloat(strings.ReplaceAll(splitLine[6], "B", ""), 64)
				}
				pool.PoolCapacityFree = fcap
				pool.PoolCapacityUsed = pool.PoolCapacity - pool.PoolCapacityFree
				pool.PoolCapacityPCT = pool.PoolCapacityUsed / pool.PoolCapacity
				pool.Firmware = firmware
				pool.Site = site
				pool.Type = type_s
				pool.Client = client_s
				output.Pools = append(output.Pools, pool)
				// fmt.Println(pool)
			}
		}

	case "3par":
		fmt.Println("3par")

	case "dell":
		response := Response{}
		xml.Unmarshal(inputData, &response)
		for _, pools := range response.Objects {
			if pools.Name == "pools" {
				var pool Pool
				pool.ArrayName = array
				pool.Id = pools.Oid
				var tcap float64
				var fcap float64
				for _, property := range pools.Properties {
					if property.Name == "name" {
						pool.PoolName = property.Value
					}

					if property.Name == "total-size" {
						if strings.Contains(property.Value, "PiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "PiB", ""), 64)
							tcap = tcap * 1024 * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "TiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "TiB", ""), 64)
							tcap = tcap * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "GiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "GiB", ""), 64)
							tcap = tcap * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "MiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "MiB", ""), 64)
							tcap = tcap * 1024 * 1024
						} else if strings.Contains(property.Value, "KiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "KiB", ""), 64)
							tcap = tcap * 1024
						} else if strings.Contains(property.Value, "iB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "iB", ""), 64)
						}
					}
					if property.Name == "total-avail" {
						if strings.Contains(property.Value, "PiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "PiB", ""), 64)
							fcap = fcap * 1024 * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "TiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "TiB", ""), 64)
							fcap = fcap * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "GiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "GiB", ""), 64)
							fcap = fcap * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "MiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "MiB", ""), 64)
							fcap = fcap * 1024 * 1024
						} else if strings.Contains(property.Value, "KiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "KiB", ""), 64)
							fcap = fcap * 1024
						} else if strings.Contains(property.Value, "iB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "iB", ""), 64)
						}
					}
				}

				pool.PoolCapacity = tcap
				pool.PoolCapacityFree = fcap
				pool.PoolCapacityUsed = pool.PoolCapacity - pool.PoolCapacityFree
				pool.PoolCapacityPCT = pool.PoolCapacityUsed / pool.PoolCapacity
				pool.Firmware = "firmware"
				pool.Site = site
				pool.Type = type_s
				pool.Client = client_s
				output.Pools = append(output.Pools, pool)
			}

		}
	}

	return output, err
}

func parseVol(inputData []byte, inputMaps []byte, model, array, site, type_s, client_s string) (output Vols, err error) {

	switch model {
	case "ibm":
		fmt.Println("ibm")
	case "huawei":
		fmt.Println("huawei")
	case "3par":
		fmt.Println("3par")

	case "dell":
		response := Response{}
		responseMaps := Response{}

		xml.Unmarshal(inputData, &response)
		xml.Unmarshal(inputMaps, &responseMaps)
		for _, vols := range response.Objects {
			if vols.Name == "volume" {
				var vol Vol
				vol.ArrayName = array
				vol.Site = site
				vol.Type = type_s
				vol.Client = client_s

				var tcap float64
				var fcap float64
				for _, property := range vols.Properties {
					if property.Name == "volume-name" {
						vol.VolName = property.Value
					}
					if property.Name == "durable-id" {
						vol.Id = property.Value
					}
					if property.Name == "allocated-size" {
						if strings.Contains(property.Value, "PiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "PiB", ""), 64)
							fcap = fcap * 1024 * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "TiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "TiB", ""), 64)
							fcap = fcap * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "GiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "GiB", ""), 64)
							fcap = fcap * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "MiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "MiB", ""), 64)
							fcap = fcap * 1024 * 1024
						} else if strings.Contains(property.Value, "KiB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "KiB", ""), 64)
							fcap = fcap * 1024
						} else if strings.Contains(property.Value, "iB") {
							fcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "iB", ""), 64)
						}
					}
					if property.Name == "total-size" {
						if strings.Contains(property.Value, "PiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "PiB", ""), 64)
							tcap = tcap * 1024 * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "TiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "TiB", ""), 64)
							tcap = tcap * 1024 * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "GiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "GiB", ""), 64)
							tcap = tcap * 1024 * 1024 * 1024
						} else if strings.Contains(property.Value, "MiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "MiB", ""), 64)
							tcap = tcap * 1024 * 1024
						} else if strings.Contains(property.Value, "KiB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "KiB", ""), 64)
							tcap = tcap * 1024
						} else if strings.Contains(property.Value, "iB") {
							tcap, err = strconv.ParseFloat(strings.ReplaceAll(property.Value, "iB", ""), 64)
						}
					}
					vol.VolAllocatedSize = fcap
					vol.VolTotalSize = tcap
					// vol.Hosts
					for _, maps := range responseMaps.Objects {
						if maps.Name == "host-view" {
							volumeCheck := false
							for _, property := range maps.Properties {
								if property.Name == "parent-id" && property.Value == vol.Id {
									volumeCheck = true
								}
							}
							for _, property := range maps.Properties {
								if property.Name == "nickname" && volumeCheck {
									add := true
									for _, member := range vol.Hosts {
										if property.Value == member {
											add = false
										}
									}
									if add {
										vol.Hosts = append(vol.Hosts, property.Value)
									}
								}
							}
						}
					}
				}

				output.Vols = append(output.Vols, vol)
			}

		}
	}

	return output, err
}

func connectToHostPW(user, password, host string) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    []ssh.AuthMethod{ssh.Password(password)},
		Timeout: 10 * time.Second,
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group-exchange-sha256")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group16-sha512")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group14-sha256")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group18-sha512")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "curve25519-sha256@libssh.org")

	client, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func connectToHostKB(user, password, host string) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			for i, _ := range answers {
				answers[i] = password
			}
			return answers, nil
		})},
		Timeout: 10 * time.Second,
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group-exchange-sha256")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "curve25519-sha256@libssh.org")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group14-sha256")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "diffie-hellman-group18-sha512")
	sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, "curve25519-sha256@libssh.org")

	client, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
func getArrays(file, model, dataString, fwString, volumeString, volumeMapString string) {
	var jsonFile *os.File
	var err error

	jsonFile, err = os.Open(file)
	if err != nil {
		logError(err.Error())
	}
	fmt.Println("Successfully Opened " + file)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var locArrays Arrays
	json.Unmarshal(byteValue, &locArrays)
	for _, array := range locArrays.Arrays {
		array.Model = model
		array.Data = dataString
		array.fw = fwString
		array.volume = volumeString
		array.volumeMaps = volumeMapString
		arrays = append(arrays, array)
	}

}

func getSwitches(file string) {
	var jsonFile *os.File
	var err error

	jsonFile, err = os.Open(file)
	if err != nil {
		logError(err.Error())
	}
	fmt.Println("Successfully Opened " + file)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var locSwitches Switches
	json.Unmarshal(byteValue, &locSwitches)
	switches = append(switches, locSwitches.Switches...)
}

// 	for i := 0; i < len(arraysDell.Arrays); i++ {

// 		model := "dell"
// 		if arraysDell.Arrays[i].Client == clientArg {
// 			logError("connecting to dell host : " + arraysDell.Arrays[i].Name)
// 			arrayPools := collectData(username, password, arraysDell.Arrays[i].Ip, arraysDell.Arrays[i].Name, arraysDell.Arrays[i].Site, arraysDell.Arrays[i].Type, arraysDell.Arrays[i].Client, model)
// 			pools.Pools = append(pools.Pools, arrayPools.Pools...)
// 			volumes := collectVolumes(username, password, arraysDell.Arrays[i].Ip, arraysDell.Arrays[i].Name, arraysDell.Arrays[i].Site, arraysDell.Arrays[i].Type, arraysDell.Arrays[i].Client, model)
// 			arrayVolume = append(arrayVolume, volumes)
// 		}
