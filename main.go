package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/donomii/goof"
)

//{ "type": "Feature", "id": "2292591621", "geometry": { "type": "Point", "coordinates": [ 121.534116, 25.0146649 ] }, "properties": {  "name": "NET", "shop": "clothes", "wheelchair": "limited", "addr:street": "羅斯福路四段", "addr:housenumber": "64", "toilets:wheelchair": "no", "wheelchair:description": "門口有一個階梯,每層樓上下只有樓梯,沒有電梯

var strict bool

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
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}
}

func main() {
	//defer profile.Start(profile.TraceProfile).Stop()
	var xmlFile *bufio.Reader
	var inFile string
	var outFile string
	var compression string

	flag.StringVar(&compression, "compression", "", "Input is compressed with bz2 or gz")
	flag.BoolVar(&strict, "strict", false, "Emit correct geojson format.  By default, emit grep-friendly geojson.")
	flag.Parse()

	args := flag.Args()
	///log.Println(args)
	if len(args) > 0 {
		inFile = args[0]
		log.Println("Reading from", inFile)
		if inFile == "-" {
			inFile = ""
		}
		xmlFile = bufio.NewReader(goof.OpenInput(inFile, compression))
		//defer xmlFile.Close()
	} else {
		log.Println("Reading from stdin")
		xmlFile = bufio.NewReader(goof.OpenInput("", compression))
	}

	outBuff := bufio.NewWriter(os.Stdout)

	if len(args) > 1 {
		outFile = args[1]
		log.Println("Writing to", outFile)
		if outFile == "-" {
		} else {
			f, err := os.Create(outFile)
			checkErr(err)
			outBuff = bufio.NewWriter(f)
		}

		//defer xmlFile.Close()
	}

	if strict {
		fmt.Fprintf(outBuff, "[")
	}
	//out := json.NewEncoder(outBuff)
	firstItem := true

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
		//log.Printf("%V\n", gen_token)
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
				//log.Println(gen_token.Attr.());
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
					// log.Printf("%s = %s\n", key, value)
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
					Coordinates: []float64{flon, flat},
				}

				g := GeoJSON{
					Type:       "Feature",
					Id:         id,
					Geometry:   geom,
					Properties: tags,
				}
				byt, err := json.Marshal(g)
				if err != nil {
					panic(fmt.Sprintf("Could not encode output data because: %v", err))
				}
				//out.Encode(g)
				if firstItem {
					firstItem = false
				} else {
					if strict {
						fmt.Fprintf(outBuff, ",")
					}
				}
				fmt.Fprintf(outBuff, "%s", string(byt))
				if !strict {
					fmt.Fprintf(outBuff, "\n")
				}
				outBuff.Flush()
				current_element = nil
				tags = nil
			}
		}
	}
	if strict {
		fmt.Fprintf(outBuff, "]")
	}
	outBuff.Flush()
	log.Println("Job's a good'un, boss!")
}
