[![Build Status](https://travis-ci.org/donomii/osm2geojson.svg?branch=master)](https://travis-ci.org/donomii/osm2geojson)

# osm2geojson
Converts Open Street Map osm files into geojson, or optionally "one feature per line" format

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
	
	or you can use the --compression flag to use the internal decompressor
	
	wget -q -O - http://download.geofabrik.de/antarctica-latest.osm.bz2 |  ./osm2geojson --compression bz2
	

## Read from file, write to stdout

    ./osm2geojson europe.osm.bz2

    osm2geojson europe.osm.bz2 

## Read from file, write to file

    ./osm2geojson europe.osm.bz2 europe.geojson.gz

    osm2geojson europe.osm.bz2 europe.geojson.gz

## Read from stdin, write to file

    ./osm2geojson - europe.geojson.gz

    osm2geojson - europe.geojson.gz

# Correct geojson format

By default, osm2geojson prints each feature on a new line, so you can use command line tools like GREP and AWK to filter the results.  e.g. to get every mountain in Europe:

	./osm2geojson europe.osm.bz2 | grep -i mountain
	
However this is not correct geojson.  To get correct geojson, add ```--strict``` to your command.

	./osm2geojson --strict europe.osm.bz2
	
Note that correct geojson is all on one line, so you can no longer use command line tools.

# Author

Adapted from http://fabsk.eu/misc/osmxml.go

Now maintained by Donomii
