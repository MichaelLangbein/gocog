package main

// https://github.com/geotiffjs/geotiff.js/blob/master/src/geotiff.js
// http://paulbourke.net/dataformats/tiff/tiff_summary.pdf
// https://www.awaresystems.be/imaging/tiff/tifftags.html
// https://www.loc.gov/preservation/digital/formats/content/tiff_tags.shtml
// https://medium.com/planet-stories/a-handy-introduction-to-cloud-optimized-geotiffs-1f2c9e716ec3

/*
 Plan:
 1. create file .................. done
 2. host file .................... done (making sure http-range requests are supported)
 3. fetch file ................... done
 4. fetch head only .............. done
 5. parse head ................... done
 6. parse with official tiff lib . done
 6. get tile
 ............. official lib doesn't really do this?
 ............. try those:
		- godal
				- Wants to fetch full image immediately, I think
		- https://github.com/gden173/geotiff/blob/main/geotiff/reader.go
				- AtCoord
				- AtPoints ......... both require file to have been read yet.
		- https://github.com/jh-sz/cog
				- Read(48) ... reads only the 49th tile
				....... only supports single band images
				....... only supports uint8, but LS8 has uint16
		- https://github.com/superztc/gocog/tree/master
				- DecodeLevel
				- DecodeLevelSumImage
				.......... more extensive, but depends on terrascope ... which is apparently out of business
				..... or is it? https://viewer.terrascope.be
				......... hmmm ... doesn't compile



 7. decompress tile
 8. get all tiles
*/

import (
	"fmt"
	"gocog/gocog"
	"gocog/selfmade"
	// "github.com/airbusgeo/godal"
)

func main() {

	fileUrl := "http://localhost:8000/LC08_L1TP_193026_20230806_20230806_02_RT_B1.TIF"
	cogReader := selfmade.MakeFetchingReader(fileUrl)

	config, _ := gocog.DecodeConfig(cogReader)
	fmt.Print(config)
	geoConfig, _ := gocog.DecodeGeoInfo(cogReader)
	fmt.Print(geoConfig)
	img, _ := gocog.DecodeLevel(cogReader, 3)
	fmt.Print(img)

	// godal.RegisterVSIHandler("http://", cogReader)
	// file, _ := godal.Open(fileUrl)
	// fmt.Print(file.Description())
}
