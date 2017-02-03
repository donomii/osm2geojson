package main

//import "github.com/pkg/profile"
import "regexp"
import "compress/bzip2"
import "compress/gzip"
import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func checkErr(err error) {
	log.Println(err)
}

func main() {
	//defer profile.Start(profile.TraceProfile).Stop()
	var xmlFile *bufio.Reader
	//log.Println(os.Args)
	if len(os.Args) > 1 {
		log.Println(os.Args[1])
		f, err := os.Open(os.Args[1])
		checkErr(err)
		xmlFile = bufio.NewReader(f)
		//defer xmlFile.Close()
	} else {
		//log.Println("Reading from stdin")
		xmlFile = bufio.NewReader(os.Stdin)
	}

	xmlFile = bufio.NewReader(bzip2.NewReader(xmlFile))
	aBuff := bufio.NewWriter(os.Stdout)

	if len(os.Args) > 2 {
		log.Println(os.Args[2])
		f, err := os.Create(os.Args[2])
		checkErr(err)
		aBuff = bufio.NewWriter(f)

		//defer xmlFile.Close()
	}

	bBuff := gzip.NewWriter(aBuff)
	out := bufio.NewWriter(bBuff)

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
				//if pagename, ok := tags["wikipedia"]; ok {
				//fmt.Println(tags)
				//}
				//fmt.Println("Attribs: ", a)
				//for k,v := range a{
				//fmt.Printf(" '%v': '%v' ", k,v)
				//}
				lat, ok := a["lat"]
				if !ok {
					lat = a["{ lat}"]
				}
				lon, ok := a["lon"]
				if !ok {
					lon = a["{ lon}"]
				}
				out.WriteString(fmt.Sprintf("{ \"type\": \"Feature\", "))
				id, ok := a["{ id}"]
				if ok {
					out.WriteString(fmt.Sprintf("\"id\": \"%v\", ", id))
				}
				out.WriteString(fmt.Sprintf("\"geometry\": { \"type\": \"Point\", \"coordinates\": [ %v, %v ] }", lon, lat))
				if len(tags) > 0 {
					out.WriteString(fmt.Sprintf(", \"properties\": { "))
					start := true
					for k, v := range tags {
						if start {
							start = false
						} else {
							out.WriteString(fmt.Sprintf(","))
						}
						re := regexp.MustCompile("\"")
						out.WriteString(fmt.Sprintf(" \"%v\": \"%v\"", re.ReplaceAllLiteralString(k, "\\\""), re.ReplaceAllLiteralString(v, "\\\"")))
					}
					out.WriteString(fmt.Sprintf(" }"))
				}
				out.WriteString(fmt.Sprintf(" }\n"))
				current_element = nil
				tags = nil
			}
		}
	}
	out.Flush()
	bBuff.Flush()
	bBuff.Close()
	aBuff.Flush()
	log.Println("Job's a good'un, boss!")
}
