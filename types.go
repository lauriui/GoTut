package main

import "encoding/xml"

type Arrays struct {
	Arrays []Array `json:"array"`
}

// User struct which contains a name
// a type and a list of social links
type Array struct {
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	Site       string `json:"site"`
	Type       string `json:"type_arr"`
	Client     string `json:"client"`
	Model      string
	Data       string
	fw         string
	volume     string
	volumeMaps string
}

type Switches struct {
	Switches []Switch `json:"switch"`
}

type Switch struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Site string `json:"site"`
	Room string `json:"room"`
}

type Ports struct {
	Ports []Port
}

type Port struct {
	Switch  string
	Index   string
	Site    string
	Room    string
	Version string
	Port    string
	Address string
	Media   string
	Speed   string
	State   string
	Proto   string
	Name    string
}

type Pools struct {
	Pools []Pool
}

type Vols struct {
	Vols []Vol
}

type Pool struct {
	Id               string
	ArrayName        string
	PoolName         string
	Firmware         string
	Site             string
	Type             string
	Client           string
	PoolCapacity     float64
	PoolCapacityFree float64
	PoolCapacityUsed float64
	PoolCapacityPCT  float64
}

type Vol struct {
	ArrayName        string
	Id               string
	VolName          string
	Site             string
	Type             string
	Client           string
	VolTotalSize     float64
	VolAllocatedSize float64
	Hosts            []string
}

type Client struct {
	Name                  string
	StretchedP16Total     float64
	StretchedP16Free      float64
	StretchedP16MinLun    int
	StretchedZ141Total    float64
	StretchedZ141Free     float64
	StretchedZ141MinLun   int
	P16Total              float64
	P16Free               float64
	P16InternalTotal      float64
	P16InternalFree       float64
	P16InternalSSDTotal   float64
	P16InternalHDDTotal   float64
	P16InternalSSDFree    float64
	P16InternalHDDFree    float64
	P16InternalSSDMinLun  int
	P16InternalHDDMinLun  int
	P16ExternalTotal      float64
	P16ExternalFree       float64
	P16ExternalSSDTotal   float64
	P16ExternalHDDTotal   float64
	P16ExternalSSDFree    float64
	P16ExternalHDDFree    float64
	P16ExternalSSDMinLun  int
	P16ExternalHDDMinLun  int
	Z141Total             float64
	Z141Free              float64
	Z141InternalTotal     float64
	Z141InternalFree      float64
	Z141InternalSSDTotal  float64
	Z141InternalHDDTotal  float64
	Z141InternalSSDFree   float64
	Z141InternalHDDFree   float64
	Z141InternalSSDMinLun int
	Z141InternalHDDMinLun int
	Z141ExternalTotal     float64
	Z141ExternalFree      float64
	Z141ExternalSSDTotal  float64
	Z141ExternalHDDTotal  float64
	Z141ExternalSSDFree   float64
	Z141ExternalHDDFree   float64
	Z141ExternalSSDMinLun int
	Z141ExternalHDDMinLun int
	Total                 float64
	TotalFree             float64
}

type Response struct {
	XMLName xml.Name `xml:"RESPONSE"`
	Version string   `xml:"VERSION,attr"`
	Request string   `xml:"REQUEST,attr"`
	Objects []Object `xml:"OBJECT"`
}

type Object struct {
	XMLName    xml.Name     `xml:"OBJECT"`
	Name       string       `xml:"name,attr"`
	Basetype   string       `xml:"basetype,attr"`
	Oid        string       `xml:"oid,attr"`
	Format     string       `xml:"format,attr"`
	Properties []Properties `xml:"PROPERTY"`
}

type Properties struct {
	XMLName xml.Name `xml:"PROPERTY"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}
