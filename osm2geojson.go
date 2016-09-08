package main

import (
    "encoding/xml"
    "fmt"
    "os"
)
func main() {
    xmlFile := os.Stdin
    defer xmlFile.Close()
    decoder := xml.NewDecoder(xmlFile)
    var current_element *xml.StartElement
    var tags map[string] string
    tags = map[string] string{}
    a:=map[string]string{}
    for {
        gen_token, _ := decoder.Token()
        if gen_token == nil {
            break
        }
        //fmt.Printf("%V\n", gen_token)
        switch token := gen_token.(type) {
        case xml.StartElement:
            if (token.Name.Local=="node") {
                a=map[string]string{}
                tags = map[string]string{}
                if se, ok := gen_token.(xml.StartElement); ok {
                    for _, v := range se.Attr {
                        a[fmt.Sprintf("%v", v.Name)] = v.Value
                    }
                }

                current_element = &token
                //fmt.Println(gen_token.Attr.());
            } else if (token.Name.Local=="tag") {
                if (current_element!=nil) {
                    var key string
                    var value string
                    for _ , attr := range token.Attr {
                        switch attr.Name.Local {
                        case "k":
                            key = attr.Value
                        case "v":
                            value = attr.Value
                        }
                    }
                    // fmt.Printf("%s = %s\n", key, value)
                    if key!="" && value!="" {
                        tags[key] = value
                    }
                }
            }
        case xml.EndElement:
            if (token.Name.Local=="node") {
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
                    lon =  a["{ lon}"]
                }
                fmt.Printf("{ \"type\": \"Feature\", ")
                id, ok := a["{ id}"]
                if ok { fmt.Printf("{ \"id\": \"%v\", ", id) }
                fmt.Printf("\"geometry\": { \"type\": \"Point\", \"coordinates\": [ %v, %v ] }, ", lon, lat)
                if len(tags)>0 {
                    fmt.Printf("\"properties\": { ")
                    start := true
                    for k,v := range tags {
                        if start {
                            start = false
                        } else {
                            fmt.Printf(",")
                        }
                        fmt.Printf(" \"%v\": \"%v\"", k,v)
                    }
                    fmt.Printf(" }")
                }
                fmt.Printf(" }\n")
                current_element = nil
                tags = nil
            }
        }
    }
}
