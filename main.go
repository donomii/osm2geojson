package main

//import "github.com/pkg/profile"
import "compress/bzip2"
import "compress/gzip"
import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//{ "type": "Feature", "id": "2292591621", "geometry": { "type": "Point", "coordinates": [ 121.534116, 25.0146649 ] }, "properties": {  "name": "NET", "shop": "clothes", "wheelchair": "limited", "addr:street": "羅斯福路四段", "addr:housenumber": "64", "toilets:wheelchair": "no", "wheelchair:description": "門口有一個階梯,每層樓上下只有樓梯,沒有電梯

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type GeoJSON struct {
	Type       string            `json:"type"`
	Id         string            `json:"id"`
	Geometry   Geometry          `json:"geometry"`
	Properties map[string]string `json:"properties"`
}

func checkErr(err error) {
	log.Println(err)
}

func main() {
	//defer profile.Start(profile.TraceProfile).Stop()
	var xmlFile *bufio.Reader
	var inFile string
	var outFile string
	//log.Println(os.Args)
	if len(os.Args) > 1 {
		inFile = os.Args[1]
		log.Println(inFile)
		if inFile == "-" {
			xmlFile = bufio.NewReader(os.Stdin)
		} else {
			f, err := os.Open(inFile)
			checkErr(err)
			xmlFile = bufio.NewReader(f)
		}
		//defer xmlFile.Close()
	} else {
		//log.Println("Reading from stdin")
		xmlFile = bufio.NewReader(os.Stdin)
	}

	if strings.HasSuffix(inFile, "bz2") {
		xmlFile = bufio.NewReader(bzip2.NewReader(xmlFile))
	}
	if strings.HasSuffix(inFile, "gz") {
		//Hooray, more golang bullshit
		g, _ := gzip.NewReader(xmlFile)
		xmlFile = bufio.NewReader(g)
	}

	outBuff := bufio.NewWriter(os.Stdout)

	if len(os.Args) > 2 {
		outFile = os.Args[2]
		log.Println(outFile)
		if outFile == "-" {
		} else {
			f, err := os.Create(outFile)
			checkErr(err)
			outBuff = bufio.NewWriter(f)
		}

		//defer xmlFile.Close()
	}

	out := json.NewEncoder(outBuff)

	decoder := xml.NewDecoder(xmlFile)
	var current_element *xml.StartElement
	var tags map[string]string
	tags = map[string]string{}
	a := map[string]string{}
	for {
		gen_token, err := decoder.Token()
		if gen_token == nil {
			log.Println("Error while decoding: ", err)
			break
		}
		//fmt.Printf("%V\n", gen_token)
		switch token := gen_token.(type) {
		case xml.StartElement:
			if token.Name.Local == "node" {
				a = map[string]string{}
				tags = map[string]string{}
				if se, ok := gen_token.(xml.StartElement); ok {
					for _, v := range se.Attr {
						a[fmt.Sprintf("%v", v.Name)] = v.Value
					}
				}

				current_element = &token
				//fmt.Println(gen_token.Attr.());
			} else if token.Name.Local == "tag" {
				if current_element != nil {
					var key string
					var value string
					for _, attr := range token.Attr {
						switch attr.Name.Local {
						case "k":
							key = attr.Value
						case "v":
							value = attr.Value
						}
					}
					// fmt.Printf("%s = %s\n", key, value)
					if key != "" && value != "" {
						tags[key] = value
					}
				}
			}
		case xml.EndElement:
			if token.Name.Local == "node" {
				lat, ok := a["lat"]
				if !ok {
					lat = a["{ lat}"]
				}
				lon, ok := a["lon"]
				if !ok {
					lon = a["{ lon}"]
				}
				id, _ := a["{ id}"]
				flat, _ := strconv.ParseFloat(lon, 64)
				flon, _ := strconv.ParseFloat(lat, 64)
				geom := Geometry{
					Type:        "Point",
					Coordinates: []float64{flat, flon},
				}

				g := GeoJSON{
					Type:       "Feature",
					Id:         id,
					Geometry:   geom,
					Properties: tags,
				}
				out.Encode(g)
				current_element = nil
				tags = nil
			}
		}
	}
	outBuff.Flush()
	log.Println("Job's a good'un, boss!")
}
