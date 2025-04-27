package linkinfo

import "encoding/xml"

type File struct {
	XMLName          xml.Name `xml:"link_info"`
	Banner           string   `xml:"banner"`
	Copyright        string   `xml:"copyright"`
	LinkTime         string   `xml:"link_time"`
	LinkErrors       string   `xml:"link_errors"`
	OutputFile       string   `xml:"output_file"`
	EntryPoint       EntryPoint
	InputFiles       []InputFile       `xml:"input_file_list>input_file"`
	ObjectComponents []ObjectComponent `xml:"object_component_list>object_component"`
	LogicalGroups    []LogicalGroup    `xml:"logical_group_list>logical_group"`
	PlacementMap     PlacementMap
	CRCTables        []CRCTable `xml:"crc_table_list>crc_table"`
	Symbols          []Symbol   `xml:"symbol_table>symbol"`
	Title            string     `xml:"title"`
}

type EntryPoint struct {
	XMLName xml.Name `xml:"entry_point"`
	Name    string   `xml:"name"`
	Address string   `xml:"address"`
}

type InputFile struct {
	XMLName xml.Name `xml:"input_file"`
	ID      string   `xml:"id,attr"`
	Path    string   `xml:"path"`
	Kind    string   `xml:"kind"`
	File    string   `xml:"file"`
	Name    string   `xml:"name"`
}

type InputFileRef struct {
	XMLName xml.Name `xml:"input_file_ref"`
	IDRef   string   `xml:"idref,attr"`
}

type ObjectComponentRef struct {
	XMLName xml.Name `xml:"object_component_ref"`
	IDRef   string   `xml:"idref,attr"`
}

type LogicalGroupRef struct {
	XMLName xml.Name `xml:"logical_group_ref"`
	IDRef   string   `xml:"idref,attr"`
}

type ObjectComponent struct {
	XMLName      xml.Name `xml:"object_component"`
	ID           string   `xml:"id,attr"`
	Name         string   `xml:"name"`
	LoadAddress  string   `xml:"load_address"`
	RunAddress   string   `xml:"run_address"`
	Size         string   `xml:"size"`
	InputFileRef InputFileRef
}

type LogicalGroup struct {
	XMLName     xml.Name `xml:"logical_group"`
	ID          string   `xml:"id,attr"`
	Display     string   `xml:"display,attr"`
	Color       string   `xml:"color,attr"`
	Name        string   `xml:"name"`
	LoadAddress string   `xml:"load_address"`
	RunAddress  string   `xml:"run_address"`
	Size        string   `xml:"size"`
	Contents    LogicalGroupContents
}

type LogicalGroupContents struct {
	XMLName             xml.Name             `xml:"contents"`
	ObjectComponentRefs []ObjectComponentRef `xml:"object_component_ref"`
}

type PlacementMap struct {
	XMLName     xml.Name     `xml:"placement_map"`
	MemoryAreas []MemoryArea `xml:"memory_area"`
}

type MemoryArea struct {
	XMLName      xml.Name `xml:"memory_area"`
	Display      string   `xml:"display,attr"`
	Color        string   `xml:"color,attr"`
	Name         string   `xml:"name"`
	Origin       string   `xml:"origin"`
	Length       string   `xml:"length"`
	UsedSpace    string   `xml:"used_space"`
	UnusedSpace  string   `xml:"unused_space"`
	Attributes   string   `xml:"attributes"`
	UsageDetails UsageDetails
}

type UsageDetails struct {
	XMLName        xml.Name `xml:"usage_details"`
	AllocatedSpace AllocatedSpace
}

type AllocatedSpace struct {
	XMLName          xml.Name          `xml:"allocated_space"`
	StartAddress     string            `xml:"start_address"`
	Size             string            `xml:"size"`
	LogicalGroupRefs []LogicalGroupRef `xml:"logical_group_ref"`
}

type CRCTable struct {
	XMLName xml.Name    `xml:"crc_table"`
	Records []CRCRecord `xml:"crc_rec"`
}

type CRCRecord struct {
	XMLName     xml.Name `xml:"crc_rec"`
	Name        string   `xml:"name"`
	AlgName     string   `xml:"alg_name"`
	AlgID       string   `xml:"alg_id"`
	LoadPageID  string   `xml:"load_page_id"`
	LoadAddress string   `xml:"load_address"`
	LoadSize    string   `xml:"load_size"`
	CRCValue    string   `xml:"crc_value"`
}

type Symbol struct {
	XMLName xml.Name `xml:"symbol"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}
