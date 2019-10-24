[![Build Status](https://travis-ci.org/donomii/osm2geojson.svg?branch=master)](https://travis-ci.org/donomii/osm2geojson)
[![GoDoc](https://godoc.org/github.com/donomii/osm2geojson?status.svg)](https://godoc.org/github.com/donomii/osm2geojson)


# osm2geojson
Converts Open Street Map osm files into geojson

Works on osm files or streams, allowing it to process files directly from the network, without downloading them first.

# To do

Note osm2geojson currently extracts only points.  More complicated structures are not implemented yet.

# To install

    go get github.com/donomii/osm2geojson
    go install github.com/donomii/osm2geojson

# To use

## As a stream processor

    cat europe.osm | ./osm2geojson

    type europe.osm | osm2geojson
	
## Converting directly from the internet

	wget -q -O - http://download.geofabrik.de/antarctica-latest.osm.bz2 | bunzip2  | ./osm2geojson

## Read from file, write to stdout

    ./osm2geojson europe.osm.bz2

    osm2geojson europe.osm.bz2 

## Read from file, write to file

    ./osm2geojson europe.osm.bz2 europe.geojson.gz

    osm2geojson europe.osm.bz2 europe.geojson.gz

# Author

Adapted from http://fabsk.eu/misc/osmxml.go

Now maintained by Donomii
