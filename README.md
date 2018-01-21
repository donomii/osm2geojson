# osm2geojson
Converts Open Street Map osm files into geojson

Works on osm files or streams, allowing it to process files directly from the network, without downloading them first.

# To do

Note it currently only extracts points.  More complicated structures are not implemented yet.

# To install

    go get github.com/donomii/osm2geojson
    go build github.com/donomii/osm2geojson

# To use

## As a stream processor

    cat europe.osm | ./osm2geojson

    type europe.osm | osm2geojson

## Read from file, write to stdout

    ./osm2geojson europe.osm.bz2

    osm2geojson europe.osm.bz2 

## Read from file, write to file

    ./osm2geojson europe.osm.bz2 europe.geojson.gz

    osm2geojson europe.osm.bz2 europe.geojson.gz

# Author

Adapted from http://fabsk.eu/misc/osmxml.go

Now maintained by Donomii
