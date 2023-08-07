package main

// https://github.com/geotiffjs/geotiff.js/blob/master/src/geotiff.js
// http://paulbourke.net/dataformats/tiff/tiff_summary.pdf
// https://www.awaresystems.be/imaging/tiff/tifftags.html
// https://www.loc.gov/preservation/digital/formats/content/tiff_tags.shtml

// Plan:
// 1. create file ................. done
// 2. host file ................... done (making sure http-range requests are supported)
// 3. fetch file .................. done
// 4. fetch head only ............. done
// 5. parse head .................. done
// 6. get tile
// 7. decompress tile
// 8. get all tiles

import (
	"fmt"
	"gocog/selfmade"

	"github.com/google/tiff"
)

func main() {

	fileUrl := "http://localhost:8000/testfile.tiff"
	cogReader := selfmade.MakeFetchingReader(fileUrl)

	t, _ := tiff.Parse(cogReader, nil, nil)

	firstIfd := t.IFDs()[0]
	tileOffsetsField := firstIfd.GetField(324)
	// tileOffsetsField.Value().
	fmt.Print(tileOffsetsField)
}
