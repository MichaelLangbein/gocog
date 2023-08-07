package selfmade

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
	"encoding/binary"
	"fmt"
)

func ReadByteOrder(word []byte) (binary.ByteOrder, error) {
	if word[0] == 0x49 && word[1] == 0x49 {
		return binary.LittleEndian, nil
	}
	if word[0] == 0x4d && word[1] == 0x4d {
		return binary.BigEndian, nil
	}
	return nil, fmt.Errorf("cannot interpret as byte-order: %x ", word[0:2])
}

func ReadVersion(word []byte, byteOrder binary.ByteOrder) (uint16, error) {
	var version = byteOrder.Uint16(word)
	if version != 42 {
		return version, fmt.Errorf("unexpected version: %d", word[0:2])
	}
	return version, nil
}

type TagID uint16

const (
	NewSubfileType               TagID = 254   //	A general indication of the kind of data contained in this subfile.	Baseline	Usage rule in JHOVE TIFF module.
	SubfileType                  TagID = 255   //	A general indication of the kind of data contained in this subfile.	Baseline
	ImageWidth                   TagID = 256   //	The number of columns in the image, i.e., the number of pixels per row.	Baseline 	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1
	ImageLength                  TagID = 257   //	The number of rows of pixels in the image.	Baseline	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1
	BitsPerSample                TagID = 258   //	Number of bits per component.	Baseline
	Compression                  TagID = 259   //	Compression scheme used on the image data.	Baseline	Sample values: 1=uncompressed and 4=CCITT Group 4.
	PhotometricInterpretation    TagID = 262   //	The color space of the image data.	Baseline 	Sample values: 1=black is zero and 2=RGB. Document also states "RGB is assumed to be sRGB; if RGB, an ICC profile should be present in the 34675 tag."
	Thresholding                 TagID = 263   //	For black and white TIFF files that represent shades of gray, the technique used to convert from gray to black and white pixels.	Baseline	Usage rule in JHOVE TIFF module.
	CellWidth                    TagID = 264   //	The width of the dithering or halftoning matrix used to create a dithered or halftoned bilevel file.	Baseline
	CellLength                   TagID = 265   //	The length of the dithering or halftoning matrix used to create a dithered or halftoned bilevel file.	Baseline	Usage rule in JHOVE TIFF module.
	FillOrder                    TagID = 266   //	The logical order of bits within a byte.	Baseline
	DocumentName                 TagID = 269   //	The name of the document from which this image was scanned.	Extended 	Also used by HD Photo.
	ImageDescription             TagID = 270   //	A string that describes the subject of the image.	Baseline 	Also used by HD Photo
	Make                         TagID = 271   //	The scanner manufacturer.	Baseline 	Also used by HD Photo.
	Model                        TagID = 272   //	The scanner model name or number.	Baseline 	Also used by HD Photo.
	StripOffsets                 TagID = 273   //	For each strip, the byte offset of that strip.	Baseline	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1 (Files outside of these classes may use tiles and tags 322, 323, 324, and 325; Comments welcome.)
	Orientation                  TagID = 274   //	The orientation of the image with respect to the rows and columns.	Baseline	Mandatory for TIFF/EP.
	SamplesPerPixel              TagID = 277   //	The number of components per pixel.	Baseline	Mandatory for TIFF 6.0 classes R and Y.1
	RowsPerStrip                 TagID = 278   //	The number of rows per strip.	Baseline	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1 (Files outside of these classes may use tiles and tags 322, 323, 324, and 325; Comments welcome.)
	StripByteCounts              TagID = 279   //	For each strip, the number of bytes in the strip after compression.	Baseline	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1 (Files outside of these classes may use tiles and tags 322, 323, 324, and 325; Comments welcome.)
	MinSampleValue               TagID = 280   //	The minimum component value used.	Baseline
	MaxSampleValue               TagID = 281   //	The maximum component value used.	Baseline
	XResolution                  TagID = 282   //	The number of pixels per ResolutionUnit in the ImageWidth direction.	Baseline	Xresolution is a Rational; ImageWidth (Tag 256) is the numerator and the length of the source (measured in the units specified in ResolutionUnit (Tag 296)) is the denominator.
	YResolution                  TagID = 283   //	The number of pixels per ResolutionUnit in the ImageLength direction.	Baseline	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1
	PlanarConfiguration          TagID = 284   //	How the components of each pixel are stored.	Baseline	Mandatory for TIFF/EP.
	PageName                     TagID = 285   //	The name of the page from which this image was scanned.	Extended 	Also used by HD Photo
	XPosition                    TagID = 286   //	X position of the image.	Extended
	YPosition                    TagID = 287   //	Y position of the image.	Extended
	FreeOffsets                  TagID = 288   //	For each string of contiguous unused bytes in a TIFF file, the byte offset of the string.	Baseline
	FreeByteCounts               TagID = 289   //	For each string of contiguous unused bytes in a TIFF file, the number of bytes in the string.	Baseline
	GrayResponseUnit             TagID = 290   //	The precision of the information contained in the GrayResponseCurve.	Baseline
	GrayResponseCurve            TagID = 291   //	For grayscale data, the optical density of each possible pixel value.	Baseline
	T4Options                    TagID = 292   //	Options for Group 3 Fax compression	Extended
	T6Options                    TagID = 293   //	Options for Group 4 Fax compression	Extended
	ResolutionUnit               TagID = 296   //	The unit of measurement for XResolution and YResolution.	Baseline	Mandatory for TIFF 6.0 classes B, G, P, R, and Y.1
	PageNumber                   TagID = 297   //	The page number of the page from which this image was scanned.	Extended 	Also used by HD Photo
	TransferFunction             TagID = 301   //	Describes a transfer function for the image in tabular style.	Extended
	Software                     TagID = 305   //	Name and version number of the software package(s) used to create the image.	Baseline 	Also used by HD Photo.
	DateTime                     TagID = 306   //	Date and time of image creation.	Baseline 	Also used by HD Photo.
	Artist                       TagID = 315   //	Person who created the image.	Baseline 	Also used by HD Photo
	HostComputer                 TagID = 316   //	The computer and/or operating system in use at the time of image creation.	Baseline 	Also used by HD Photo
	Predictor                    TagID = 317   //	A mathematical operator that is applied to the image data before an encoding scheme is applied.	Extended
	WhitePoint                   TagID = 318   //	The chromaticity of the white point of the image.	Extended
	PrimaryChromaticities        TagID = 319   //	The chromaticities of the primaries of the image.	Extended
	ColorMap                     TagID = 320   //	A color map for palette color images.	Baseline	Mandatory for TIFF 6.0 class P.1
	HalftoneHints                TagID = 321   //	Conveys to the halftone function the range of gray levels within a colorimetrically-specified image that should retain tonal detail.	Extended
	TileWidth                    TagID = 322   //	The tile width in pixels. This is the number of columns in each tile.	Extended	Mandatory for TIFF 6.0 files that use tiles. (Files that use strips employ tags 273, 278, and 279.)
	TileLength                   TagID = 323   //	The tile length (height) in pixels. This is the number of rows in each tile.	Extended	Referenced in JHOVE TIFF module for files that use tiles. (Files that use strips employ tags 273, 278, and 279.)
	TileOffsets                  TagID = 324   //	For each tile, the byte offset of that tile, as compressed and stored on disk.	Extended	Mandatory for TIFF 6.0 files that use tiles. (Files that use strips employ tags 273, 278, and 279.)
	TileByteCounts               TagID = 325   //	For each tile, the number of (compressed) bytes in that tile.	Extended	Mandatory for TIFF 6.0 files that use tiles. (Files that use strips employ tags 273, 278, and 279.)
	BadFaxLines                  TagID = 326   //	Used in the TIFF-F standard, denotes the number of 'bad' scan lines encountered by the facsimile device.	Extended
	CleanFaxData                 TagID = 327   //	Used in the TIFF-F standard, indicates if 'bad' lines encountered during reception are stored in the data, or if 'bad' lines have been replaced by the receiver.	Extended
	ConsecutiveBadFaxLines       TagID = 328   //	Used in the TIFF-F standard, denotes the maximum number of consecutive 'bad' scanlines received.	Extended
	SubIFDs                      TagID = 330   //	Offset to child IFDs.	Extended
	InkSet                       TagID = 332   //	The set of inks used in a separated (PhotometricInterpretation=5) image.	Extended
	InkNames                     TagID = 333   //	The name of each ink used in a separated image.	Extended
	NumberOfInks                 TagID = 334   //	The number of inks.	Extended
	DotRange                     TagID = 336   //	The component values that correspond to a 0% dot and 100% dot.	Extended	Usage rule in JHOVE TIFF module.
	TargetPrinter                TagID = 337   //	A description of the printing environment for which this separation is intended.	Extended
	ExtraSamples                 TagID = 338   //	Description of extra components.	Baseline
	SampleFormat                 TagID = 339   //	Specifies how to interpret each data sample in a pixel.	Extended
	SMinSampleValue              TagID = 340   //	Specifies the minimum sample value.	Extended
	SMaxSampleValue              TagID = 341   //	Specifies the maximum sample value.	Extended
	TransferRange                TagID = 342   //	Expands the range of the TransferFunction.	Extended
	ClipPath                     TagID = 343   //	Mirrors the essentials of PostScript's path creation functionality.	Extended	Usage rule in JHOVE TIFF module.
	XClipPathUnits               TagID = 344   //	The number of units that span the width of the image, in terms of integer ClipPath coordinates.	Extended	Usage rule in JHOVE TIFF module.
	YClipPathUnits               TagID = 345   //	The number of units that span the height of the image, in terms of integer ClipPath coordinates.	Extended
	Indexed                      TagID = 346   //	Aims to broaden the support for indexed images to include support for any color space.	Extended
	JPEGTables                   TagID = 347   //	JPEG quantization and/or Huffman tables.	Extended
	OPIProxy                     TagID = 351   //	OPI-related.	Extended
	GlobalParametersIFD          TagID = 400   //	Used in the TIFF-FX standard to point to an IFD containing tags that are globally applicable to the complete TIFF file.	Extended
	ProfileType                  TagID = 401   //	Used in the TIFF-FX standard, denotes the type of data stored in this file or IFD.	Extended
	FaxProfile                   TagID = 402   //	Used in the TIFF-FX standard, denotes the 'profile' that applies to this file.	Extended
	CodingMethods                TagID = 403   //	Used in the TIFF-FX standard, indicates which coding methods are used in the file.	Extended
	VersionYear                  TagID = 404   //	Used in the TIFF-FX standard, denotes the year of the standard specified by the FaxProfile field.	Extended
	ModeNumber                   TagID = 405   //	Used in the TIFF-FX standard, denotes the mode of the standard specified by the FaxProfile field.	Extended
	Decode                       TagID = 433   //	Used in the TIFF-F and TIFF-FX standards, holds information about the ITULAB (PhotometricInterpretation = 10) encoding.	Extended
	DefaultImageColor            TagID = 434   //	Defined in the Mixed Raster Content part of RFC 2301, is the default color needed in areas where no image is available.	Extended
	JPEGProc                     TagID = 512   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGInterchangeFormat        TagID = 513   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGInterchangeFormatLength  TagID = 514   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGRestartInterval          TagID = 515   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGLosslessPredictors       TagID = 517   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGPointTransforms          TagID = 518   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGQTables                  TagID = 519   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGDCTables                 TagID = 520   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	JPEGACTables                 TagID = 521   //	Old-style JPEG compression field. TechNote2 invalidates this part of the specification.	Extended
	YCbCrCoefficients            TagID = 529   //	The transformation from RGB to YCbCr image data.	Extended	Mandatory for TIFF/EP YCbCr images.
	YCbCrSubSampling             TagID = 530   //	Specifies the subsampling factors used for the chrominance components of a YCbCr image.	Extended	Mandatory for TIFF/EP YCbCr images.
	YCbCrPositioning             TagID = 531   //	Specifies the positioning of subsampled chrominance components relative to luminance samples.	Extended	Mandatory for TIFF/EP YCbCr images.
	ReferenceBlackWhite          TagID = 532   //	Specifies a pair of headroom and footroom image data values (codes) for each pixel component.	Extended	Mandatory for TIFF 6.0 class Y.1
	StripRowCounts               TagID = 559   //	Defined in the Mixed Raster Content part of RFC 2301, used to replace RowsPerStrip for IFDs with variable-sized strips.	Extended
	XMP                          TagID = 700   //	XML packet containing XMP metadata	Extended 	Also used by HD Photo
	ImageRating                  TagID = 18246 //	Ratings tag used by Windows	Exif private IFD
	ImageRatingPercent           TagID = 18249 //		Ratings tag used by Windows, value as percent	Exif private IFD
	ImageID                      TagID = 32781 //	OPI-related.	Extended
	Wang                         TagID = 32932 // Annotation	Annotation data, as used in 'Imaging for Windows'.	Private
	CFARepeatPatternDim          TagID = 33421 //	For camera raw files from sensors with CFA overlay. 	Mandatory in TIFF/EP for CFA files.
	CFAPattern                   TagID = 33422 //	For camera raw files from sensors with CFA overlay. 	Mandatory in TIFF/EP for CFA files.
	BatteryLevel                 TagID = 33423 //	Encodes camera battery level at time of image capture.
	Copyright                    TagID = 33432 //	Copyright notice.	Baseline 	Also used by HD Photo.
	ExposureTime                 TagID = 33434 //	Exposure time, given in seconds.	Exif Private IFD
	FNumber                      TagID = 33437 //	The F number.	Exif Private IFD
	MDFileTag                    TagID = 33445 // 	Specifies the pixel data format encoding in the Molecular Dynamics GEL file format.	Private
	MDScalePixel                 TagID = 33446 // 	Specifies a scale factor in the Molecular Dynamics GEL file format.	Private
	MDColorTable                 TagID = 33447 // 	Used to specify the conversion from 16bit to 8bit in the Molecular Dynamics GEL file format.	Private
	MDLabName                    TagID = 33448 // 	Name of the lab that scanned this file, as used in the Molecular Dynamics GEL file format.	Private
	MDSampleInfo                 TagID = 33449 // 	Information about the sample, as used in the Molecular Dynamics GEL file format.	Private
	MDPrepDate                   TagID = 33450 // 	Date the sample was prepared, as used in the Molecular Dynamics GEL file format.	Private
	MDPrepTime                   TagID = 33451 // 	Time the sample was prepared, as used in the Molecular Dynamics GEL file format.	Private
	MDFileUnits                  TagID = 33452 // 	Units for data in this file, as used in the Molecular Dynamics GEL file format.	Private
	ModelPixelScaleTag           TagID = 33550 //	Used in interchangeable GeoTIFF_1_0 files.	Private
	IPTC                         TagID = 33723 ///NAA	IPTC-NAA (International Press Telecommunications Council-Newspaper Association of America) metadata. 	TIFF/EP spec, p. 33 	Tag name and values defined by IPTC-NAA Info Interchange Model & Digital Newsphoto Parameter Record.
	INGRPacketDataTag            TagID = 33918 // 	Intergraph Application specific storage.	Private
	INGRFlagRegisters            TagID = 33919 // 	Intergraph Application specific flags.	Private
	IrasB                        TagID = 33920 // Transformation Matrix	Originally part of Intergraph's GeoTIFF tags, but likely understood by IrasB only.	Private
	ModelTiepointTag             TagID = 33922 //	Originally part of Intergraph's GeoTIFF tags, but now used in interchangeable GeoTIFF_1_0 files.	Private 	In GeoTIFF_1_0, either this tag or 34264 must be defined, but not both
	Site                         TagID = 34016 //		Site where image created. 	TIFF/IT spec, 7.2.3
	ColorSequence                TagID = 34017 //		Sequence of colors if other than CMYK. 	TIFF/IT spec, 7.2.8.3.2 	For BP and BL only2
	IT8Header                    TagID = 34018 //		Certain inherited headers. 	TIFF/IT spec, 7.2.3 	Obsolete
	RasterPadding                TagID = 34019 //		Type of raster padding used, if any. 	TIFF/IT spec, 7.2.6
	BitsPerRunLength             TagID = 34020 //		Number of bits for short run length encoding. 	TIFF/IT spec, 7.2.6 	For LW only2
	BitsPerExtendedRunLength     TagID = 34021 //		Number of bits for long run length encoding. 	TIFF/IT spec, 7.2.6 	For LW only2
	ColorTable                   TagID = 34022 //		Color value in a color pallette. 	TIFF/IT spec, 7.2.8.4 	For BP and BL only2
	ImageColorIndicator          TagID = 34023 //		Indicates if image (foreground) color or transparency is specified. 	TIFF/IT spec, 7.2.9 	For MP, BP, and BL only2
	BackgroundColorIndicator     TagID = 34024 //		Background color specification. 	TIFF/IT spec, 7.2.9 	For BP and BL only2
	ImageColorValue              TagID = 34025 //		Specifies image (foreground) color. 	TIFF/IT spec, 7.2.8.4 	For MP, BP, and BL only2
	BackgroundColorValue         TagID = 34026 //		Specifies background color. 	TIFF/IT spec, 7.2.8.4 	For BP and BL only2
	PixelIntensityRange          TagID = 34027 //		Specifies data values for 0 percent and 100 percent pixel intensity. 	TIFF/IT spec, 7.2.8.4 	For MP only2
	TransparencyIndicator        TagID = 34028 //		Specifies if transparency is used in HC file. 	TIFF/IT spec, 7.2.8.4 	For HC only2
	ColorCharacterization        TagID = 34029 //		Specifies ASCII table or other reference per ISO 12641 and ISO 12642. 	TIFF/IT spec, 7.2.8.4
	HCUsage                      TagID = 34030 //		Indicates the type of information in an HC file. 	TIFF/IT spec, 7.2.6 	For HC only2
	TrapIndicator                TagID = 34031 //		Indicates whether or not trapping has been applied to the file. 	TIFF/IT spec, 7.2.6
	CMYKEquivalent               TagID = 34032 //		Specifies CMYK equivalent for specific separations. 	TIFF/IT spec, 7.2.8.3.4
	Reserved1                    TagID = 34033 //		For future TIFF/IT use	TIFF/IT spec
	Reserved2                    TagID = 34034 //		For future TIFF/IT use	TIFF/IT spec
	Reserved3                    TagID = 34035 //		For future TIFF/IT use	TIFF/IT spec
	ModelTransformationTag       TagID = 34264 //	Used in interchangeable GeoTIFF_1_0 files.	Private 	In GeoTIFF_1_0, either this tag or 33922 must be defined, but not both
	Photoshop                    TagID = 34377 //	Collection of Photoshop 'Image Resource Blocks'.	Private
	Exif                         TagID = 34665 // IFD	A pointer to the Exif IFD.	Private 	Also used by HD Photo.
	InterColorProfile            TagID = 34675 //	ICC profile data. 	Also called ICC Profile.
	ImageLayer                   TagID = 34732 //	Defined in the Mixed Raster Content part of RFC 2301, used to denote the particular function of this Image in the mixed raster scheme.	Extended
	GeoKeyDirectoryTag           TagID = 34735 //	Used in interchangeable GeoTIFF_1_0 files.	Private 	Mandatory in GeoTIFF_1_0
	GeoDoubleParamsTag           TagID = 34736 //	Used in interchangeable GeoTIFF_1_0 files.	Private
	GeoAsciiParamsTag            TagID = 34737 //	Used in interchangeable GeoTIFF_1_0 files.	Private
	ExposureProgram              TagID = 34850 //	The class of the program used by the camera to set exposure when the picture is taken.	Exif Private IFD
	SpectralSensitivity          TagID = 34852 //	Indicates the spectral sensitivity of each channel of the camera used.	Exif Private IFD
	GPSInfo                      TagID = 34853 //	A pointer to the Exif-related GPS Info IFD. 	Also called GPS IFD.
	ISOSpeedRatings              TagID = 34855 //	Indicates the ISO Speed and ISO Latitude of the camera or input device as specified in ISO 12232.	Exif Private IFD
	OECF                         TagID = 34856 //	Indicates the Opto-Electric Conversion Function (OECF) specified in ISO 14524.	Exif Private IFD
	Interlace                    TagID = 34857 //	Indicates the field number of multifield images.
	TimeZoneOffset               TagID = 34858 //	Encodes time zone of camera clock relative to GMT.
	SelfTimeMode                 TagID = 34859 //	Number of seconds image capture was delayed from button press.
	SensitivityType              TagID = 34864 //	The SensitivityType tag indicates PhotographicSensitivity tag, which one of the parameters of ISO 12232. Although it is an optional tag, it should be recorded when a PhotographicSensitivity tag is recorded. Value = 4, 5, 6, or 7 may be used in case that the values of plural parameters are the same.	Exif private IFD
	StandardOutputSensitivity    TagID = 34865 //	This tag indicates the standard output sensitivity value of a camera or input device defined in ISO 12232. When recording this tag, the PhotographicSensitivity and SensitivityType tags shall also be recorded.	Exif private IFD
	RecommendedExposureIndex     TagID = 34866 //	This tag indicates the recommended exposure index value of a camera or input device defined in ISO 12232. When recording this tag, the PhotographicSensitivity and SensitivityType tags shall also be recorded.	Exif private IFD
	ISOSpeed                     TagID = 34867 //	This tag indicates the ISO speed value of a camera or input device that is defined in ISO 12232. When recording this tag, the PhotographicSensitivity and SensitivityType tags shall also be recorded.	Exif private IFD
	ISOSpeedLatitudeyyy          TagID = 34868 //	This tag indicates the ISO speed latitude yyy value of a camera or input device that is defined in ISO 12232. However, this tag shall not be recorded without ISOSpeed and ISOSpeedLatitudezzz.	Exif private IFD
	ISOSpeedLatitudezzz          TagID = 34869 //	This tag indicates the ISO speed latitude zzz value of a camera or input device that is defined in ISO 12232. However, this tag shall not be recorded without ISOSpeed and ISOSpeedLatitudeyyy.	Exif private IFD
	HylaFAXFaxRecvParams         TagID = 34908 // 	Used by HylaFAX.	Private
	HylaFAXFaxSubAddress         TagID = 34909 // 	Used by HylaFAX.	Private
	HylaFAXFaxRecvTime           TagID = 34910 // 	Used by HylaFAX.	Private
	ExifVersion                  TagID = 36864 //	The version of the supported Exif standard.	Exif Private IFD 	Mandatory in the Exif IFD.
	DateTimeOriginal             TagID = 36867 //	The date and time when the original image data was generated.	Exif Private IFD
	DateTimeDigitized            TagID = 36868 //	The date and time when the image was stored as digital data.	Exif Private IFD
	ComponentsConfiguration      TagID = 37121 //	Specific to compressed data; specifies the channels and complements PhotometricInterpretation	Exif Private IFD
	CompressedBitsPerPixel       TagID = 37122 //	Specific to compressed data; states the compressed bits per pixel.	Exif Private IFD
	ShutterSpeedValue            TagID = 37377 //	Shutter speed.	Exif Private IFD
	ApertureValue                TagID = 37378 //	The lens aperture.	Exif Private IFD
	BrightnessValue              TagID = 37379 //	The value of brightness.	Exif Private IFD
	ExposureBiasValue            TagID = 37380 //	The exposure bias.	Exif Private IFD
	MaxApertureValue             TagID = 37381 //	The smallest F number of the lens.	Exif Private IFD
	SubjectDistance              TagID = 37382 //	The distance to the subject, given in meters.	Exif Private IFD
	MeteringMode                 TagID = 37383 //	The metering mode.	Exif Private IFD
	LightSource                  TagID = 37384 //	The kind of light source.	Exif Private IFD
	Flash                        TagID = 37385 //	Indicates the status of flash when the image was shot.	Exif Private IFD
	FocalLength                  TagID = 37386 //	The actual focal length of the lens, in mm.	Exif Private IFD
	FlashEnergy                  TagID = 37387 //	Amount of flash energy (BCPS).
	SpatialFrequencyResponse     TagID = 37388 //	SFR of the camera.
	Noise                        TagID = 37389 //	Noise measurement values.
	FocalPlaneXResolution        TagID = 37390 //	Number of pixels per FocalPlaneResolutionUnit (37392) in ImageWidth direction for main image.
	FocalPlaneYResolution        TagID = 37391 //	Number of pixels per FocalPlaneResolutionUnit (37392) in ImageLength direction for main image.
	FocalPlaneResolutionUnit     TagID = 37392 //	Unit of measurement for FocalPlaneXResolution(37390) and FocalPlaneYResolution(37391).
	ImageNumber                  TagID = 37393 //	Number assigned to an image, e.g., in a chained image burst.
	SecurityClassification       TagID = 37394 //	Security classification assigned to the image.
	ImageHistory                 TagID = 37395 //	Record of what has been done to the image.
	SubjectLocation              TagID = 37396 //	Indicates the location and area of the main subject in the overall scene.	Exif Private IFD
	ExposureIndex                TagID = 37397 //	Encodes the camera exposure index setting when image was captured.
	TIFF                         TagID = 37398 ///EPStandardID	For current spec, tag value equals 1 0 0 0. 	Mandatory in TIFF/EP.
	SensingMethod                TagID = 37399 //	Type of image sensor. 	Mandatory in TIFF/EP.
	MakerNote                    TagID = 37500 //	Manufacturer specific information.	Exif Private IFD
	UserComment                  TagID = 37510 //	Keywords or comments on the image; complements ImageDescription.	Exif Private IFD
	SubsecTime                   TagID = 37520 //	A tag used to record fractions of seconds for the DateTime tag.	Exif Private IFD
	SubsecTimeOriginal           TagID = 37521 //	A tag used to record fractions of seconds for the DateTimeOriginal tag.	Exif Private IFD
	SubsecTimeDigitized          TagID = 37522 //	A tag used to record fractions of seconds for the DateTimeDigitized tag.	Exif Private IFD
	ImageSourceData              TagID = 37724 //	Used by Adobe Photoshop.	Private
	XPTitle                      TagID = 40091 //	Title tag used by Windows, encoded in UCS2	Exif Private IFD
	XPComment                    TagID = 40092 //	Comment tag used by Windows, encoded in UCS2	Exif Private IFD
	XPAuthor                     TagID = 40093 //	Author tag used by Windows, encoded in UCS2	Exif Private IFD
	XPKeywords                   TagID = 40094 //	Keywords tag used by Windows, encoded in UCS2	Exif Private IFD
	XPSubject                    TagID = 40095 //	Subject tag used by Windows, encoded in UCS2	Exif Private IFD
	FlashpixVersion              TagID = 40960 //	The Flashpix format version supported by a FPXR file.	Exif Private IFD 	Mandatory in the Exif IFD
	ColorSpace                   TagID = 40961 //	The color space information tag is always recorded as the color space specifier.	Exif Private IFD 	Mandatory in the Exif IFD
	PixelXDimension              TagID = 40962 //	Specific to compressed data; the valid width of the meaningful image.	Exif Private IFD
	PixelYDimension              TagID = 40963 //	Specific to compressed data; the valid height of the meaningful image.	Exif Private IFD
	RelatedSoundFile             TagID = 40964 //	Used to record the name of an audio file related to the image data.	Exif Private IFD
	Interoperability             TagID = 40965 // IFD	A pointer to the Exif-related Interoperability IFD.	Private
	FlashEnergy1                 TagID = 41483 //	Indicates the strobe energy at the time the image is captured, as measured in Beam Candle Power Seconds	Exif Private IFD
	SpatialFrequencyResponse1    TagID = 41484 //	Records the camera or input device spatial frequency table and SFR values in the direction of image width, image height, and diagonal direction, as specified in ISO 12233.	Exif Private IFD
	FocalPlaneXResolution1       TagID = 41486 //	Indicates the number of pixels in the image width (X) direction per FocalPlaneResolutionUnit on the camera focal plane.	Exif Private IFD
	FocalPlaneYResolution1       TagID = 41487 //	Indicates the number of pixels in the image height (Y) direction per FocalPlaneResolutionUnit on the camera focal plane.	Exif Private IFD
	FocalPlaneResolutionUnit1    TagID = 41488 //	Indicates the unit for measuring FocalPlaneXResolution and FocalPlaneYResolution.	Exif Private IFD
	SubjectLocation1             TagID = 41492 //	Indicates the location of the main subject in the scene.	Exif Private IFD
	ExposureIndex1               TagID = 41493 //	Indicates the exposure index selected on the camera or input device at the time the image is captured.	Exif Private IFD
	SensingMethod1               TagID = 41495 //	Indicates the image sensor type on the camera or input device.	Exif Private IFD
	FileSource                   TagID = 41728 //	Indicates the image source.	Exif Private IFD
	SceneType                    TagID = 41729 //	Indicates the type of scene.	Exif Private IFD
	CFAPattern1                  TagID = 41730 //	Indicates the color filter array (CFA) geometric pattern of the image sensor when a one-chip color area sensor is used.	Exif Private IFD
	CustomRendered               TagID = 41985 //	Indicates the use of special processing on image data, such as rendering geared to output.	Exif Private IFD
	ExposureMode                 TagID = 41986 //	Indicates the exposure mode set when the image was shot.	Exif Private IFD
	WhiteBalance                 TagID = 41987 //	Indicates the white balance mode set when the image was shot.	Exif Private IFD
	DigitalZoomRatio             TagID = 41988 //	Indicates the digital zoom ratio when the image was shot.	Exif Private IFD
	FocalLengthIn35mmFilm        TagID = 41989 //	Indicates the equivalent focal length assuming a 35mm film camera, in mm.	Exif Private IFD
	SceneCaptureType             TagID = 41990 //	Indicates the type of scene that was shot.	Exif Private IFD
	GainControl                  TagID = 41991 //	Indicates the degree of overall image gain adjustment.	Exif Private IFD
	Contrast                     TagID = 41992 //	Indicates the direction of contrast processing applied by the camera when the image was shot.	Exif Private IFD
	Saturation                   TagID = 41993 //	Indicates the direction of saturation processing applied by the camera when the image was shot.	Exif Private IFD
	Sharpness                    TagID = 41994 //	Indicates the direction of sharpness processing applied by the camera when the image was shot.	Exif Private IFD
	DeviceSettingDescription     TagID = 41995 //	This tag indicates information on the picture-taking conditions of a particular camera model.	Exif Private IFD
	SubjectDistanceRange         TagID = 41996 //	Indicates the distance to the subject.	Exif Private IFD
	ImageUniqueID                TagID = 42016 //	Indicates an identifier assigned uniquely to each image.	Exif Private IFD
	CameraOwnerName              TagID = 42032 //	Camera owner name as ASCII string.	Exif Private IFD
	BodySerialNumber             TagID = 42033 //	Camera body serial number as ASCII string.	Exif Private IFD
	LensSpecification            TagID = 42034 //	This tag notes minimum focal length, maximum focal length, minimum F number in the minimum focal length, and minimum F number in the maximum focal length, which are specification information for the lens that was used in photography. When the minimum F number is unknown, the notation is 0/0.	Exif Private IFD
	LensMake                     TagID = 42035 //	Lens manufacturer name as ASCII string.	Exif Private IFD
	LensModel                    TagID = 42036 //	Lens model name and number as ASCII string.	Exif Private IFD
	LensSerialNumber             TagID = 42037 //	Lens serial number as ASCII string.	Exif Private IFD
	GDAL_METADATA                TagID = 42112 //	Used by the GDAL library, holds an XML list of name=value 'metadata' values about the image as a whole, and about specific samples.	Private
	GDAL_NODATA                  TagID = 42113 //	Used by the GDAL library, contains an ASCII encoded nodata or background pixel value.	Private
	PixelFormat                  TagID = 48129 //	A 128-bit Globally Unique Identifier (GUID) that identifies the image pixel format.	HD Photo Feature Spec, p. 17
	Transformation               TagID = 48130 //	Specifies the transformation to be applied when decoding the image to present the desired representation.	HD Photo Feature Spec, p. 23
	Uncompressed                 TagID = 48131 //	Specifies that image data is uncompressed.	HD Photo Feature Spec, p. 23
	ImageTypePhoto               TagID = 48132 //	Specifies the image type of each individual frame in a multi-frame file.	HD Photo Feature Spec, p. 27
	ImageWidthPhoto              TagID = 48256 //	Specifies the number of columns in the transformed photo, or the number of pixels per scan line.	HD Photo Feature Spec, p. 21
	ImageHeight                  TagID = 48257 //	Specifies the number of pixels or scan lines in the transformed photo.	HD Photo Feature Spec, p. 21
	WidthResolution              TagID = 48258 //	Specifies the horizontal resolution of a transformed image expressed in pixels per inch.	HD Photo Feature Spec, p. 21
	HeightResolution             TagID = 48259 //	Specifies the vertical resolution of a transformed image expressed in pixels per inch.	HD Photo Feature Spec, p. 21
	ImageOffset                  TagID = 48320 //	Specifies the byte offset pointer to the beginning of the photo data, relative to the beginning of the file.	HD Photo Feature Spec, p. 22
	ImageByteCount               TagID = 48321 //	Specifies the size of the photo in bytes.	HD Photo Feature Spec, p. 22
	AlphaOffset                  TagID = 48322 //	Specifies the byte offset pointer the beginning of the planar alpha channel data, relative to the beginning of the file.	HD Photo Feature Spec, p. 22
	AlphaByteCount               TagID = 48323 //	Specifies the size of the alpha channel data in bytes.	HD Photo Feature Spec, p. 23
	ImageDataDiscard             TagID = 48324 //	Signifies the level of data that has been discarded from the image as a result of a compressed domain transcode to reduce the file size.	HD Photo Feature Spec, p. 25
	AlphaDataDiscard             TagID = 48325 //	Signifies the level of data that has been discarded from the planar alpha channel as a result of a compressed domain transcode to reduce the file size.	HD Photo Feature Spec, p. 26
	ImageType                    TagID = 48132 //	Specifies the image type of each individual frame in a multi-frame file.	HD Photo Feature Spec, p. 27
	OceScanjobDescription        TagID = 50215 // 	Used in the Oce scanning process.	Private
	OceApplicationSelector       TagID = 50216 // 	Used in the Oce scanning process.	Private
	OceIdentificationNumber      TagID = 50217 // 	Used in the Oce scanning process.	Private
	OceImageLogicCharacteristics TagID = 50218 // 	Used in the Oce scanning process.	Private
	PrintImageMatching           TagID = 50341 //	Description needed.	Exif Private IFD
	DNGVersion                   TagID = 50706 // 	Encodes DNG four-tier version number; for version 1.1.0.0, the tag contains the bytes 1, 1, 0, 0. Used in IFD 0 of DNG files.
	DNGBackwardVersion           TagID = 50707 // 	Defines oldest version of spec with which file is compatible. Used in IFD 0 of DNG files.
	UniqueCameraModel            TagID = 50708 // 	Unique, non-localized nbame for camera model. Used in IFD 0 of DNG files.
	LocalizedCameraModel         TagID = 50709 // 	Similar to 50708, with localized camera name. Used in IFD 0 of DNG files.
	CFAPlaneColor                TagID = 50710 // 	Mapping between values in the CFAPattern tag and the plane numbers in LinearRaw space. Used in Raw IFD of DNG files. 	Required for non-RGB CFA images.
	CFALayout                    TagID = 50711 // 	Spatial layout of the CFA. Used in Raw IFD of DNG files.
	LinearizationTable           TagID = 50712 // 	Lookup table that maps stored values to linear values. Used in Raw IFD of DNG files.
	BlackLevelRepeatDim          TagID = 50713 // 	Repeat pattern size for BlackLevel tag. Used in Raw IFD of DNG files.
	BlackLevel                   TagID = 50714 // 	Specifies the zero light encoding level.Used in Raw IFD of DNG files.
	BlackLevelDeltaH             TagID = 50715 // 	Specifies the difference between zero light encoding level for each column and the baseline zero light encoding level. Used in Raw IFD of DNG files.
	BlackLevelDeltaV             TagID = 50716 // 	Specifies the difference between zero light encoding level for each row and the baseline zero light encoding level. Used in Raw IFD of DNG files.
	WhiteLevel                   TagID = 50717 // 	Specifies the fully saturated encoding level for the raw sample values. Used in Raw IFD of DNG files.
	DefaultScale                 TagID = 50718 // 	For cameras with non-square pixels, specifies the default scale factors for each direction to convert the image to square pixels. Used in Raw IFD of DNG files.
	DefaultCropOrigin            TagID = 50719 // 	Specifies the origin of the final image area, ignoring the extra pixels at edges used to prevent interpolation artifacts. Used in Raw IFD of DNG files.
	DefaultCropSize              TagID = 50720 // 	Specifies size of final image area in raw image coordinates. Used in Raw IFD of DNG files.
	ColorMatrix1                 TagID = 50721 // 	Defines a transformation matrix that converts XYZ values to reference camera native color space values, under the first calibration illuminant. Used in IFD 0 of DNG files.
	ColorMatrix2                 TagID = 50722 // 	Defines a transformation matrix that converts XYZ values to reference camera native color space values, under the second calibration illuminant. Used in IFD 0 of DNG files.
	CameraCalibration1           TagID = 50723 // 	Defines a calibration matrix that transforms reference camera native space values to individual camera native space values under the first calibration illuminant. Used in IFD 0 of DNG files.
	CameraCalibration2           TagID = 50724 // 	Defines a calibration matrix that transforms reference camera native space values to individual camera native space values under the second calibration illuminant. Used in IFD 0 of DNG files.
	ReductionMatrix1             TagID = 50725 // 	Defines a dimensionality reduction matrix for use as the first stage in converting color camera native space values to XYZ values, under the first calibration illuminant. Used in IFD 0 of DNG files.
	ReductionMatrix2             TagID = 50726 // 	Defines a dimensionality reduction matrix for use as the first stage in converting color camera native space values to XYZ values, under the second calibration illuminant. Used in IFD 0 of DNG files.
	AnalogBalance                TagID = 50727 // 	Pertaining to white balance, defines the gain, either analog or digital, that has been applied to the stored raw values. Used in IFD 0 of DNG files.
	AsShotNeutral                TagID = 50728 // 	Specifies the selected white balance at the time of capture, encoded as the coordinates of a perfectly neutral color in linear reference space values. Used in IFD 0 of DNG files.
	AsShotWhiteXY                TagID = 50729 // 	Specifies the selected white balance at the time of capture, encoded as x-y chromaticity coordinates. Used in IFD 0 of DNG files.
	BaselineExposure             TagID = 50730 // 	Specifies in EV units how much to move the zero point for exposure compensation. Used in IFD 0 of DNG files.
	BaselineNoise                TagID = 50731 // 	Specifies the relative noise of the camera model at a baseline ISO value of 100, compared to reference camera model. Used in IFD 0 of DNG files.
	BaselineSharpness            TagID = 50732 // 	Specifies the relative amount of sharpening required for this camera model, compared to reference camera model. Used in IFD 0 of DNG files.
	BayerGreenSplit              TagID = 50733 // 	For CFA images, specifies, in arbitrary units, how closely the values of the green pixels in the blue/green rows track the values of the green pixels in the red/green rows. Used in Raw IFD of DNG files.
	LinearResponseLimit          TagID = 50734 // 	Specifies the fraction of the encoding range above which the response may become significantly non-linear. Used in IFD 0 of DNG files.
	CameraSerialNumber           TagID = 50735 // 	Serial number of camera. Used in IFD 0 of DNG files.
	LensInfo                     TagID = 50736 // 	Information about the lens. Used in IFD 0 of DNG files.
	ChromaBlurRadius             TagID = 50737 // 	Normally for non-CFA images, provides a hint about how much chroma blur ought to be applied. Used in Raw IFD of DNG files.
	AntiAliasStrength            TagID = 50738 // 	Provides a hint about the strength of the camera's anti-aliasing filter. Used in Raw IFD of DNG files.
	ShadowScale                  TagID = 50739 //	 	Used by Adobe Camera Raw to control sensitivity of its shadows slider. Used in IFD 0 of DNG files.
	DNGPrivateData               TagID = 50740 // 	Provides a way for camera manufacturers to store private data in DNG files for use by their own raw convertors. Used in IFD 0 of DNG files.
	MakerNoteSafety              TagID = 50741 // 	Lets the DNG reader know whether the Exif MakerNote tag is safe to preserve. Used in IFD 0 of DNG files.
	CalibrationIlluminant1       TagID = 50778 // 	Illuminant used for first set of calibration tags. Used in IFD 0 of DNG files.
	CalibrationIlluminant2       TagID = 50779 // 	Illuminant used for second set of calibration tags. Used in IFD 0 of DNG files.
	BestQualityScale             TagID = 50780 // 	Specifies the amount by which the values of the DefaultScale tag need to be multiplied to achieve best quality image size. Used in Raw IFD of DNG files.
	RawDataUniqueID              TagID = 50781 //	 	Contains a 16-byte unique identifier for the raw image file in the DNG file. Used in IFD 0 of DNG files.
	Alias                        TagID = 50784 // Layer Metadata	Alias Sketchbook Pro layer usage description.	Private
	OriginalRawFileName          TagID = 50827 //	 	Name of original file if the DNG file results from conversion from a non-DNG raw file. Used in IFD 0 of DNG files.
	OriginalRawFileData          TagID = 50828 //	 	If the DNG file was converted from a non-DNG raw file, then this tag contains the original raw data. Used in IFD 0 of DNG files.
	ActiveArea                   TagID = 50829 //	 	Defines the active (non-masked) pixels of the sensor. Used in Raw IFD of DNG files.
	MaskedAreas                  TagID = 50830 //	 	List of non-overlapping rectangle coordinates of fully masked pixels, which can optimally be used by DNG readers to measure the black encoding level. Used in Raw IFD of DNG files.
	AsShotICCProfile             TagID = 50831 //	 	Contains ICC profile that, in conjunction with the AsShotPreProfileMatrix tag, specifies a default color rendering from camera color space coordinates (linear reference values) into the ICC profile connection space. Used in IFD 0 of DNG files.
	AsShotPreProfileMatrix       TagID = 50832 //	 	Specifies a matrix that should be applied to the camera color space coordinates before processing the values through the ICC profile specified in the AsShotICCProfile tag. Used in IFD 0 of DNG files.
	CurrentICCProfile            TagID = 50833 //	 	The CurrentICCProfile and CurrentPreProfileMatrix tags have the same purpose and usage as
	CurrentPreProfileMatrix      TagID = 50834 //	 	The CurrentICCProfile and CurrentPreProfileMatrix tags have the same purpose and usage as
	ColorimetricReference        TagID = 50879 // 	The DNG color model documents a transform between camera colors and CIE XYZ values. This tag describes the colorimetric reference for the CIE XYZ values. 0 = The XYZ values are scene-referred. 1 = The XYZ values are output-referred, using the ICC profile perceptual dynamic range. Used in IFD 0 of DNG files.
	CameraCalibrationSignature   TagID = 50931 // 	A UTF-8 encoded string associated with the CameraCalibration1 and CameraCalibration2 tags. Used in IFD 0 of DNG files.
	ProfileCalibrationSignature  TagID = 50932 // 	A UTF-8 encoded string associated with the camera profile tags. Used in IFD 0 or CameraProfile IFD of DNG files.
	ExtraCameraProfiles          TagID = 50933 // 	A list of file offsets to extra Camera Profile IFDs. The format of a camera profile begins with a 16-bit byte order mark (MM or II) followed by a 16-bit "magic" number equal to 0x4352 (CR), a 32-bit IFD offset, and then a standard TIFF format IFD. Used in IFD 0 of DNG files.
	AsShotProfileName            TagID = 50934 // 	A UTF-8 encoded string containing the name of the "as shot" camera profile, if any. Used in IFD 0 of DNG files.
	NoiseReductionApplied        TagID = 50935 // 	Indicates how much noise reduction has been applied to the raw data on a scale of 0.0 to 1.0. A 0.0 value indicates that no noise reduction has been applied. A 1.0 value indicates that the "ideal" amount of noise reduction has been applied, i.e. that the DNG reader should not apply additional noise reduction by default. A value of 0/0 indicates that this parameter is unknown. Used in Raw IFD of DNG files.
	ProfileName                  TagID = 50936 // 	A UTF-8 encoded string containing the name of the camera profile. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileHueSatMapDims         TagID = 50937 // 	Specifies the number of input samples in each dimension of the hue/saturation/value mapping tables. The data for these tables are stored in ProfileHueSatMapData1 and ProfileHueSatMapData2 tags. Allowed values include the following: HueDivisions >= 1; SaturationDivisions >= 2; ValueDivisions >=1. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileHueSatMapData1        TagID = 50938 // 	Contains the data for the first hue/saturation/value mapping table. Each entry of the table contains three 32-bit IEEE floating-point values. The first entry is hue shift in degrees; the second entry is saturation scale factor; and the third entry is a value scale factor. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileHueSatMapData2        TagID = 50939 // 	Contains the data for the second hue/saturation/value mapping table. Each entry of the table contains three 32-bit IEEE floating-point values. The first entry is hue shift in degrees; the second entry is saturation scale factor; and the third entry is a value scale factor. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileToneCurve             TagID = 50940 // 	Contains a default tone curve that can be applied while processing the image as a starting point for user adjustments. The curve is specified as a list of 32-bit IEEE floating-point value pairs in linear gamma. Each sample has an input value in the range of 0.0 to 1.0, and an output value in the range of 0.0 to 1.0. The first sample is required to be (0.0, 0.0), and the last sample is required to be (1.0, 1.0). Interpolated the curve using a cubic spline. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileEmbedPolicy           TagID = 50941 // 	Contains information about the usage rules for the associated camera profile. The valid values and meanings are: 0 = “allow copying”; 1 = “embed if used”; 2 = “embed never”; and 3 = “no restrictions”. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileCopyright             TagID = 50942 // 	Contains information about the usage rules for the associated camera profile. The valid values and meanings are: 0 = “allow copying”; 1 = “embed if used”; 2 = “embed never”; and 3 = “no restrictions”. Used in IFD 0 or Camera Profile IFD of DNG files.
	ForwardMatrix1               TagID = 50964 // 	Defines a matrix that maps white balanced camera colors to XYZ D50 colors. Used in IFD 0 or Camera Profile IFD of DNG files.
	ForwardMatrix2               TagID = 50965 // 	Defines a matrix that maps white balanced camera colors to XYZ D50 colors. Used in IFD 0 or Camera Profile IFD of DNG files.
	PreviewApplicationName       TagID = 50966 // 	A UTF-8 encoded string containing the name of the application that created the preview stored in the IFD. Used in Preview IFD of DNG files.
	PreviewApplicationVersion    TagID = 50967 // 	A UTF-8 encoded string containing the version number of the application that created the preview stored in the IFD. Used in Preview IFD of DNG files.
	PreviewSettingsName          TagID = 50968 // 	A UTF-8 encoded string containing the name of the conversion settings (for example, snapshot name) used for the preview stored in the IFD. Used in Preview IFD of DNG files.
	PreviewSettingsDigest        TagID = 50969 // 	A unique ID of the conversion settings (for example, MD5 digest) used to render the preview stored in the IFD. Used in Preview IFD of DNG files.
	PreviewColorSpace            TagID = 50970 // 	This tag specifies the color space in which the rendered preview in this IFD is stored. The valid values include: 0 = Unknown; 1 = Gray Gamma 2.2; 2 = sRGB; 3 = Adobe RGB; and 4 = ProPhoto RGB. Used in Preview IFD of DNG files.
	PreviewDateTime              TagID = 50971 // 	This tag is an ASCII string containing the name of the date/time at which the preview stored in the IFD was rendered, encoded using ISO 8601 format. Used in Preview IFD of DNG files.
	RawImageDigest               TagID = 50972 // 	MD5 digest of the raw image data. All pixels in the image are processed in row-scan order. Each pixel is zero padded to 16 or 32 bits deep (16-bit for data less than or equal to 16 bits deep, 32-bit otherwise). The data is processed in little-endian byte order. Used in IFD 0 of DNG files.
	OriginalRawFileDigest        TagID = 50973 // 	MD5 digest of the data stored in the OriginalRawFileData tag. Used in IFD 0 of DNG files.
	SubTileBlockSize             TagID = 50974 // 	Normally, pixels within a tile are stored in simple row-scan order. This tag specifies that the pixels within a tile should be grouped first into rectangular blocks of the specified size. These blocks are stored in row-scan order. Within each block, the pixels are stored in row-scan order. Used in Raw IFD of DNG files.
	RowInterleaveFactor          TagID = 50975 // 	Specifies that rows of the image are stored in interleaved order. The value of the tag specifies the number of interleaved fields. Used in Raw IFD of DNG files.
	ProfileLookTableDims         TagID = 50981 // 	Specifies the number of input samples in each dimension of a default "look" table. The data for this table is stored in the ProfileLookTableData tag. Allowed values include: HueDivisions >= 1; SaturationDivisions >= 2; and ValueDivisions >= 1. Used in IFD 0 or Camera Profile IFD of DNG files.
	ProfileLookTableData         TagID = 50982 // 	Default "look" table that can be applied while processing the image as a starting point for user adjustment. This table uses the same format as the tables stored in the ProfileHueSatMapData1 and ProfileHueSatMapData2 tags, and is applied in the same color space. However, it should be applied later in the processing pipe, after any exposure compensation and/or fill light stages, but before any tone curve stage. Each entry of the table contains three 32-bit IEEE floating-point values. The first entry is hue shift in degrees, the second entry is a saturation scale factor, and the third entry is a value scale factor. Used in IFD 0 or Camera Profile IFD of DNG files.
	OpcodeList1                  TagID = 51008 // 	Specifies the list of opcodes (image processing operation codes) that should be applied to the raw image, as read directly from the file. Used in Raw IFD of DNG files.
	OpcodeList2                  TagID = 51009 // 	Specifies the list of opcodes (image processing operation codes) that should be applied to the raw image, just after it has been mapped to linear reference values. Used in Raw IFD of DNG files.
	OpcodeList3                  TagID = 51022 // 	Specifies the list of opcodes (image processing operation codes) that should be applied to the raw image, just after it has been demosaiced. Used in Raw IFD of DNG files.
	NoiseProfile                 TagID = 51041 // 	Describes the amount of noise in a raw image; models the amount of signal-dependent photon (shot) noise and signal-independent sensor readout noise, two common sources of noise in raw images. Used in Raw IFD of DNG files.
	OriginalDefaultFinalSize     TagID = 51089 // 	If this file is a proxy for a larger original DNG file, this tag specifics the default final size of the larger original file from which this proxy was generated. The default value for this tag is default final size of the current DNG file, which is DefaultCropSize * DefaultScale. 	DNG spec (1.4, 2012), p. 74
	OriginalBestQualityFinalSize TagID = 51090 // 	If this file is a proxy for a larger original DNG file, this tag specifics the best quality final size of the larger original file from which this proxy was generated. The default value for this tag is the OriginalDefaultFinalSize, if specified. Otherwise the default value for this tag is the best quality size of the current DNG file, which is DefaultCropSize * DefaultScale * BestQualityScale. 	DNG spec (1.4, 2012), p. 75
	OriginalDefaultCropSize      TagID = 51091 // 	If this file is a proxy for a larger original DNG file, this tag specifics the DefaultCropSize of the larger original file from which this proxy was generated. The default value for this tag is the OriginalDefaultFinalSize, if specified. Otherwise, the default value for this tag is the DefaultCropSize of the current DNG file. 	DNG spec (1.4, 2012), p. 75
	ProfileHueSatMapEncoding     TagID = 51107 // 	Provides a way for color profiles to specify how indexing into a 3D HueSatMap is performed during raw conversion. This tag is not applicable to 2.5D HueSatMap tables (i.e., where the Value dimension is 1). The currently defined values are: 0 = Linear encoding (method described in DNG spec); 1 = sRGB encoding (method described in DNG spec). 	DNG spec (1.4, 2012), p. 73
	ProfileLookTableEncoding     TagID = 51108 // 	Provides a way for color profiles to specify how indexing into a 3D LookTable is performed during raw conversion. This tag is not applicable to a 2.5D LookTable (i.e., where the Value dimension is 1). The currently defined values are: 0 = Linear encoding (method described in DNG spec); 1 = sRGB encoding (method described in DNG spec). 	DNG spec (1.4, 2012), p. 72-3
	BaselineExposureOffset       TagID = 51109 // 	Provides a way for color profiles to increase or decrease exposure during raw conversion. BaselineExposureOffset specifies the amount (in EV units) to add to th e BaselineExposure tag during image rendering. For example, if the BaselineExposure value fo r a given camera model is +0.3, and the BaselineExposureOffset value for a given camera profile used to render an image for that camera model is -0.7, then th e actual default exposure value used during rendering will be +0.3 - 0.7 = -0.4. 	DNG spec (1.4, 2012), p. 71
	DefaultBlackRender           TagID = 51110 // 	This optional tag in a color profile provides a hint to the raw converter regarding how to handle the black point (e.g., flare subtraction) during rendering. The currently defined values are: 0 = Auto; 1 = None. If set to Auto, the raw converter should perform black subtraction during rendering. The amount and method of black subtraction may be automatically determined and may be image-dependent. If set to None, the raw converter should not perform any black subtraction during rendering. This may be desirable when using color lookup tables (e.g., LookTable) or tone curves in camera profiles to perform a fixed, consistent level of black subtraction. 	DNG spec (1.4, 2012), p. 71
	NewRawImageDigest            TagID = 51111 // 	This tag is a modified MD5 digest of the raw image data. It has been updated from the algorithm used to compute the RawImageDigest tag be more multi-processor friendly, and to support lossy compression algorithms. The details of the algorithm used to compute this tag are documented in the Adobe DNG SDK source code. 	DNG spec (1.4, 2012), p. 76
	RawToPreviewGain             TagID = 51112 // 	The gain (what number the sample values are multiplied by) between the main raw IFD and the preview IFD containing this tag. 	DNG spec (1.4, 2012), p. 76
	DefaultUserCrop              TagID = 51125 // 	Specifies a default user crop rectangle in relative coordinates. The values must satisfy: 0.0 <= top < bottom <= 1.0; 0.0 <= left < right <= 1.0. The default values of (top = 0, left = 0, bottom = 1, right = 1) correspond exactly to the default crop rectangle (as specified by the DefaultCropOrigin and DefaultCropSize tags). 	DNG spec (1.4, 2012), p. 70
)

func (tid TagID) String() string {
	switch tid {
	case NewSubfileType:
		return "NewSubfileType"
	case SubfileType:
		return "SubfileType"
	case ImageWidth:
		return "ImageWidth"
	case ImageLength:
		return "ImageLength"
	case BitsPerSample:
		return "BitsPerSample"
	case Compression:
		return "Compression"
	case PhotometricInterpretation:
		return "PhotometricInterpretation"
	case Thresholding:
		return "Thresholding"
	case CellWidth:
		return "CellWidth"
	case CellLength:
		return "CellLength"
	case FillOrder:
		return "FillOrder"
	case DocumentName:
		return "DocumentName"
	case ImageDescription:
		return "ImageDescription"
	case Make:
		return "Make"
	case Model:
		return "Model"
	case StripOffsets:
		return "StripOffsets"
	case Orientation:
		return "Orientation"
	case SamplesPerPixel:
		return "SamplesPerPixel"
	case RowsPerStrip:
		return "RowsPerStrip"
	case StripByteCounts:
		return "StripByteCounts"
	case MinSampleValue:
		return "MinSampleValue"
	case MaxSampleValue:
		return "MaxSampleValue"
	case XResolution:
		return "XResolution"
	case YResolution:
		return "YResolution"
	case PlanarConfiguration:
		return "PlanarConfiguration"
	case PageName:
		return "PageName"
	case XPosition:
		return "XPosition"
	case YPosition:
		return "YPosition"
	case FreeOffsets:
		return "FreeOffsets"
	case FreeByteCounts:
		return "FreeByteCounts"
	case GrayResponseUnit:
		return "GrayResponseUnit"
	case GrayResponseCurve:
		return "GrayResponseCurve"
	case T4Options:
		return "T4Options"
	case T6Options:
		return "T6Options"
	case ResolutionUnit:
		return "ResolutionUnit"
	case PageNumber:
		return "PageNumber"
	case TransferFunction:
		return "TransferFunction"
	case Software:
		return "Software"
	case DateTime:
		return "DateTime"
	case Artist:
		return "Artist"
	case HostComputer:
		return "HostComputer"
	case Predictor:
		return "Predictor"
	case WhitePoint:
		return "WhitePoint"
	case PrimaryChromaticities:
		return "PrimaryChromaticities"
	case ColorMap:
		return "ColorMap"
	case HalftoneHints:
		return "HalftoneHints"
	case TileWidth:
		return "TileWidth"
	case TileLength:
		return "TileLength"
	case TileOffsets:
		return "TileOffsets"
	case TileByteCounts:
		return "TileByteCounts"
	case BadFaxLines:
		return "BadFaxLines"
	case CleanFaxData:
		return "CleanFaxData"
	case ConsecutiveBadFaxLines:
		return "ConsecutiveBadFaxLines"
	case SubIFDs:
		return "SubIFDs"
	case InkSet:
		return "InkSet"
	case InkNames:
		return "InkNames"
	case NumberOfInks:
		return "NumberOfInks"
	case DotRange:
		return "DotRange"
	case TargetPrinter:
		return "TargetPrinter"
	case ExtraSamples:
		return "ExtraSamples"
	case SampleFormat:
		return "SampleFormat"
	case SMinSampleValue:
		return "SMinSampleValue"
	case SMaxSampleValue:
		return "SMaxSampleValue"
	case TransferRange:
		return "TransferRange"
	case ClipPath:
		return "ClipPath"
	case XClipPathUnits:
		return "XClipPathUnits"
	case YClipPathUnits:
		return "YClipPathUnits"
	case Indexed:
		return "Indexed"
	case JPEGTables:
		return "JPEGTables"
	case OPIProxy:
		return "OPIProxy"
	case GlobalParametersIFD:
		return "GlobalParametersIFD"
	case ProfileType:
		return "ProfileType"
	case FaxProfile:
		return "FaxProfile"
	case CodingMethods:
		return "CodingMethods"
	case VersionYear:
		return "VersionYear"
	case ModeNumber:
		return "ModeNumber"
	case Decode:
		return "Decode"
	case DefaultImageColor:
		return "DefaultImageColor"
	case JPEGProc:
		return "JPEGProc"
	case JPEGInterchangeFormat:
		return "JPEGInterchangeFormat"
	case JPEGInterchangeFormatLength:
		return "JPEGInterchangeFormatLength"
	case JPEGRestartInterval:
		return "JPEGRestartInterval"
	case JPEGLosslessPredictors:
		return "JPEGLosslessPredictors"
	case JPEGPointTransforms:
		return "JPEGPointTransforms"
	case JPEGQTables:
		return "JPEGQTables"
	case JPEGDCTables:
		return "JPEGDCTables"
	case JPEGACTables:
		return "JPEGACTables"
	case YCbCrCoefficients:
		return "YCbCrCoefficients"
	case YCbCrSubSampling:
		return "YCbCrSubSampling"
	case YCbCrPositioning:
		return "YCbCrPositioning"
	case ReferenceBlackWhite:
		return "ReferenceBlackWhite"
	case StripRowCounts:
		return "StripRowCounts"
	case XMP:
		return "XMP"
	case ImageRating:
		return "ImageRating"
	case ImageRatingPercent:
		return "ImageRatingPercent"
	case ImageID:
		return "ImageID"
	case Wang:
		return "Wang"
	case CFARepeatPatternDim:
		return "CFARepeatPatternDim"
	case CFAPattern:
		return "CFAPattern"
	case BatteryLevel:
		return "BatteryLevel"
	case Copyright:
		return "Copyright"
	case ExposureTime:
		return "ExposureTime"
	case FNumber:
		return "FNumber"
	case MDFileTag:
		return "MDFileTag"
	case MDScalePixel:
		return "MDScalePixel"
	case MDColorTable:
		return "MDColorTable"
	case MDLabName:
		return "MDLabName"
	case MDSampleInfo:
		return "MDSampleInfo"
	case MDPrepDate:
		return "MDPrepDate"
	case MDPrepTime:
		return "MDPrepTime"
	case MDFileUnits:
		return "MDFileUnits"
	case ModelPixelScaleTag:
		return "ModelPixelScaleTag"
	case IPTC:
		return "IPTC"
	case INGRPacketDataTag:
		return "INGRPacketDataTag"
	case INGRFlagRegisters:
		return "INGRFlagRegisters"
	case IrasB:
		return "IrasB"
	case ModelTiepointTag:
		return "ModelTiepointTag"
	case Site:
		return "Site"
	case ColorSequence:
		return "ColorSequence"
	case IT8Header:
		return "IT8Header"
	case RasterPadding:
		return "RasterPadding"
	case BitsPerRunLength:
		return "BitsPerRunLength"
	case BitsPerExtendedRunLength:
		return "BitsPerExtendedRunLength"
	case ColorTable:
		return "ColorTable"
	case ImageColorIndicator:
		return "ImageColorIndicator"
	case BackgroundColorIndicator:
		return "BackgroundColorIndicator"
	case ImageColorValue:
		return "ImageColorValue"
	case BackgroundColorValue:
		return "BackgroundColorValue"
	case PixelIntensityRange:
		return "PixelIntensityRange"
	case TransparencyIndicator:
		return "TransparencyIndicator"
	case ColorCharacterization:
		return "ColorCharacterization"
	case HCUsage:
		return "HCUsage"
	case TrapIndicator:
		return "TrapIndicator"
	case CMYKEquivalent:
		return "CMYKEquivalent"
	case Reserved1:
		return "Reserved1"
	case Reserved2:
		return "Reserved2"
	case Reserved3:
		return "Reserved3"
	case ModelTransformationTag:
		return "ModelTransformationTag"
	case Photoshop:
		return "Photoshop"
	case Exif:
		return "Exif"
	case InterColorProfile:
		return "InterColorProfile"
	case ImageLayer:
		return "ImageLayer"
	case GeoKeyDirectoryTag:
		return "GeoKeyDirectoryTag"
	case GeoDoubleParamsTag:
		return "GeoDoubleParamsTag"
	case GeoAsciiParamsTag:
		return "GeoAsciiParamsTag"
	case ExposureProgram:
		return "ExposureProgram"
	case SpectralSensitivity:
		return "SpectralSensitivity"
	case GPSInfo:
		return "GPSInfo"
	case ISOSpeedRatings:
		return "ISOSpeedRatings"
	case OECF:
		return "OECF"
	case Interlace:
		return "Interlace"
	case TimeZoneOffset:
		return "TimeZoneOffset"
	case SelfTimeMode:
		return "SelfTimeMode"
	case SensitivityType:
		return "SensitivityType"
	case StandardOutputSensitivity:
		return "StandardOutputSensitivity"
	case RecommendedExposureIndex:
		return "RecommendedExposureIndex"
	case ISOSpeed:
		return "ISOSpeed"
	case ISOSpeedLatitudeyyy:
		return "ISOSpeedLatitudeyyy"
	case ISOSpeedLatitudezzz:
		return "ISOSpeedLatitudezzz"
	case HylaFAXFaxRecvParams:
		return "HylaFAXFaxRecvParams"
	case HylaFAXFaxSubAddress:
		return "HylaFAXFaxSubAddress"
	case HylaFAXFaxRecvTime:
		return "HylaFAXFaxRecvTime"
	case ExifVersion:
		return "ExifVersion"
	case DateTimeOriginal:
		return "DateTimeOriginal"
	case DateTimeDigitized:
		return "DateTimeDigitized"
	case ComponentsConfiguration:
		return "ComponentsConfiguration"
	case CompressedBitsPerPixel:
		return "CompressedBitsPerPixel"
	case ShutterSpeedValue:
		return "ShutterSpeedValue"
	case ApertureValue:
		return "ApertureValue"
	case BrightnessValue:
		return "BrightnessValue"
	case ExposureBiasValue:
		return "ExposureBiasValue"
	case MaxApertureValue:
		return "MaxApertureValue"
	case SubjectDistance:
		return "SubjectDistance"
	case MeteringMode:
		return "MeteringMode"
	case LightSource:
		return "LightSource"
	case Flash:
		return "Flash"
	case FocalLength:
		return "FocalLength"
	case FlashEnergy:
		return "FlashEnergy"
	case SpatialFrequencyResponse:
		return "SpatialFrequencyResponse"
	case Noise:
		return "Noise"
	case FocalPlaneXResolution:
		return "FocalPlaneXResolution"
	case FocalPlaneYResolution:
		return "FocalPlaneYResolution"
	case FocalPlaneResolutionUnit:
		return "FocalPlaneResolutionUnit"
	case ImageNumber:
		return "ImageNumber"
	case SecurityClassification:
		return "SecurityClassification"
	case ImageHistory:
		return "ImageHistory"
	case SubjectLocation:
		return "SubjectLocation"
	case ExposureIndex:
		return "ExposureIndex"
	case TIFF:
		return "TIFF"
	case SensingMethod:
		return "SensingMethod"
	case MakerNote:
		return "MakerNote"
	case UserComment:
		return "UserComment"
	case SubsecTime:
		return "SubsecTime"
	case SubsecTimeOriginal:
		return "SubsecTimeOriginal"
	case SubsecTimeDigitized:
		return "SubsecTimeDigitized"
	case ImageSourceData:
		return "ImageSourceData"
	case XPTitle:
		return "XPTitle"
	case XPComment:
		return "XPComment"
	case XPAuthor:
		return "XPAuthor"
	case XPKeywords:
		return "XPKeywords"
	case XPSubject:
		return "XPSubject"
	case FlashpixVersion:
		return "FlashpixVersion"
	case ColorSpace:
		return "ColorSpace"
	case PixelXDimension:
		return "PixelXDimension"
	case PixelYDimension:
		return "PixelYDimension"
	case RelatedSoundFile:
		return "RelatedSoundFile"
	case Interoperability:
		return "Interoperability"
	case FlashEnergy1:
		return "FlashEnergy1"
	case SpatialFrequencyResponse1:
		return "SpatialFrequencyResponse1"
	case FocalPlaneXResolution1:
		return "FocalPlaneXResolution1"
	case FocalPlaneYResolution1:
		return "FocalPlaneYResolution1"
	case FocalPlaneResolutionUnit1:
		return "FocalPlaneResolutionUnit1"
	case SubjectLocation1:
		return "SubjectLocation1"
	case ExposureIndex1:
		return "ExposureIndex1"
	case SensingMethod1:
		return "SensingMethod1"
	case FileSource:
		return "FileSource"
	case SceneType:
		return "SceneType"
	case CFAPattern1:
		return "CFAPattern1"
	case CustomRendered:
		return "CustomRendered"
	case ExposureMode:
		return "ExposureMode"
	case WhiteBalance:
		return "WhiteBalance"
	case DigitalZoomRatio:
		return "DigitalZoomRatio"
	case FocalLengthIn35mmFilm:
		return "FocalLengthIn35mmFilm"
	case SceneCaptureType:
		return "SceneCaptureType"
	case GainControl:
		return "GainControl"
	case Contrast:
		return "Contrast"
	case Saturation:
		return "Saturation"
	case Sharpness:
		return "Sharpness"
	case DeviceSettingDescription:
		return "DeviceSettingDescription"
	case SubjectDistanceRange:
		return "SubjectDistanceRange"
	case ImageUniqueID:
		return "ImageUniqueID"
	case CameraOwnerName:
		return "CameraOwnerName"
	case BodySerialNumber:
		return "BodySerialNumber"
	case LensSpecification:
		return "LensSpecification"
	case LensMake:
		return "LensMake"
	case LensModel:
		return "LensModel"
	case LensSerialNumber:
		return "LensSerialNumber"
	case GDAL_METADATA:
		return "GDAL_METADATA"
	case GDAL_NODATA:
		return "GDAL_NODATA"
	case PixelFormat:
		return "PixelFormat"
	case Transformation:
		return "Transformation"
	case Uncompressed:
		return "Uncompressed"
	case ImageWidthPhoto:
		return "ImageWidthPhoto"
	case ImageHeight:
		return "ImageHeight"
	case WidthResolution:
		return "WidthResolution"
	case HeightResolution:
		return "HeightResolution"
	case ImageOffset:
		return "ImageOffset"
	case ImageByteCount:
		return "ImageByteCount"
	case AlphaOffset:
		return "AlphaOffset"
	case AlphaByteCount:
		return "AlphaByteCount"
	case ImageDataDiscard:
		return "ImageDataDiscard"
	case AlphaDataDiscard:
		return "AlphaDataDiscard"
	case ImageType:
		return "ImageType"
	case OceScanjobDescription:
		return "OceScanjobDescription"
	case OceApplicationSelector:
		return "OceApplicationSelector"
	case OceIdentificationNumber:
		return "OceIdentificationNumber"
	case OceImageLogicCharacteristics:
		return "OceImageLogicCharacteristics"
	case PrintImageMatching:
		return "PrintImageMatching"
	case DNGVersion:
		return "DNGVersion"
	case DNGBackwardVersion:
		return "DNGBackwardVersion"
	case UniqueCameraModel:
		return "UniqueCameraModel"
	case LocalizedCameraModel:
		return "LocalizedCameraModel"
	case CFAPlaneColor:
		return "CFAPlaneColor"
	case CFALayout:
		return "CFALayout"
	case LinearizationTable:
		return "LinearizationTable"
	case BlackLevelRepeatDim:
		return "BlackLevelRepeatDim"
	case BlackLevel:
		return "BlackLevel"
	case BlackLevelDeltaH:
		return "BlackLevelDeltaH"
	case BlackLevelDeltaV:
		return "BlackLevelDeltaV"
	case WhiteLevel:
		return "WhiteLevel"
	case DefaultScale:
		return "DefaultScale"
	case DefaultCropOrigin:
		return "DefaultCropOrigin"
	case DefaultCropSize:
		return "DefaultCropSize"
	case ColorMatrix1:
		return "ColorMatrix1"
	case ColorMatrix2:
		return "ColorMatrix2"
	case CameraCalibration1:
		return "CameraCalibration1"
	case CameraCalibration2:
		return "CameraCalibration2"
	case ReductionMatrix1:
		return "ReductionMatrix1"
	case ReductionMatrix2:
		return "ReductionMatrix2"
	case AnalogBalance:
		return "AnalogBalance"
	case AsShotNeutral:
		return "AsShotNeutral"
	case AsShotWhiteXY:
		return "AsShotWhiteXY"
	case BaselineExposure:
		return "BaselineExposure"
	case BaselineNoise:
		return "BaselineNoise"
	case BaselineSharpness:
		return "BaselineSharpness"
	case BayerGreenSplit:
		return "BayerGreenSplit"
	case LinearResponseLimit:
		return "LinearResponseLimit"
	case CameraSerialNumber:
		return "CameraSerialNumber"
	case LensInfo:
		return "LensInfo"
	case ChromaBlurRadius:
		return "ChromaBlurRadius"
	case AntiAliasStrength:
		return "AntiAliasStrength"
	case ShadowScale:
		return "ShadowScale"
	case DNGPrivateData:
		return "DNGPrivateData"
	case MakerNoteSafety:
		return "MakerNoteSafety"
	case CalibrationIlluminant1:
		return "CalibrationIlluminant1"
	case CalibrationIlluminant2:
		return "CalibrationIlluminant2"
	case BestQualityScale:
		return "BestQualityScale"
	case RawDataUniqueID:
		return "RawDataUniqueID"
	case Alias:
		return "Alias"
	case OriginalRawFileName:
		return "OriginalRawFileName"
	case OriginalRawFileData:
		return "OriginalRawFileData"
	case ActiveArea:
		return "ActiveArea"
	case MaskedAreas:
		return "MaskedAreas"
	case AsShotICCProfile:
		return "AsShotICCProfile"
	case AsShotPreProfileMatrix:
		return "AsShotPreProfileMatrix"
	case CurrentICCProfile:
		return "CurrentICCProfile"
	case CurrentPreProfileMatrix:
		return "CurrentPreProfileMatrix"
	case ColorimetricReference:
		return "ColorimetricReference"
	case CameraCalibrationSignature:
		return "CameraCalibrationSignature"
	case ProfileCalibrationSignature:
		return "ProfileCalibrationSignature"
	case ExtraCameraProfiles:
		return "ExtraCameraProfiles"
	case AsShotProfileName:
		return "AsShotProfileName"
	case NoiseReductionApplied:
		return "NoiseReductionApplied"
	case ProfileName:
		return "ProfileName"
	case ProfileHueSatMapDims:
		return "ProfileHueSatMapDims"
	case ProfileHueSatMapData1:
		return "ProfileHueSatMapData1"
	case ProfileHueSatMapData2:
		return "ProfileHueSatMapData2"
	case ProfileToneCurve:
		return "ProfileToneCurve"
	case ProfileEmbedPolicy:
		return "ProfileEmbedPolicy"
	case ProfileCopyright:
		return "ProfileCopyright"
	case ForwardMatrix1:
		return "ForwardMatrix1"
	case ForwardMatrix2:
		return "ForwardMatrix2"
	case PreviewApplicationName:
		return "PreviewApplicationName"
	case PreviewApplicationVersion:
		return "PreviewApplicationVersion"
	case PreviewSettingsName:
		return "PreviewSettingsName"
	case PreviewSettingsDigest:
		return "PreviewSettingsDigest"
	case PreviewColorSpace:
		return "PreviewColorSpace"
	case PreviewDateTime:
		return "PreviewDateTime"
	case RawImageDigest:
		return "RawImageDigest"
	case OriginalRawFileDigest:
		return "OriginalRawFileDigest"
	case SubTileBlockSize:
		return "SubTileBlockSize"
	case RowInterleaveFactor:
		return "RowInterleaveFactor"
	case ProfileLookTableDims:
		return "ProfileLookTableDims"
	case ProfileLookTableData:
		return "ProfileLookTableData"
	case OpcodeList1:
		return "OpcodeList1"
	case OpcodeList2:
		return "OpcodeList2"
	case OpcodeList3:
		return "OpcodeList3"
	case NoiseProfile:
		return "NoiseProfile"
	case OriginalDefaultFinalSize:
		return "OriginalDefaultFinalSize"
	case OriginalBestQualityFinalSize:
		return "OriginalBestQualityFinalSize"
	case OriginalDefaultCropSize:
		return "OriginalDefaultCropSize"
	case ProfileHueSatMapEncoding:
		return "ProfileHueSatMapEncoding"
	case ProfileLookTableEncoding:
		return "ProfileLookTableEncoding"
	case BaselineExposureOffset:
		return "BaselineExposureOffset"
	case DefaultBlackRender:
		return "DefaultBlackRender"
	case NewRawImageDigest:
		return "NewRawImageDigest"
	case RawToPreviewGain:
		return "RawToPreviewGain"
	case DefaultUserCrop:
		return "DefaultUserCrop"
	}
	return fmt.Sprintf("unknown(%d)", tid)
}

type TagDataType uint16

const (
	BYTE      TagDataType = 1  // 8-bit unsigned integer
	ASCII     TagDataType = 2  // 8-bit, NULL-terminated string
	SHORT     TagDataType = 3  // 16-bit unsigned integer
	LONG      TagDataType = 4  // 32-bit unsigned integer
	RATIONAL  TagDataType = 5  // Two 32-bit unsigned integers
	SBYTE     TagDataType = 6  // 8-bit signed integer
	UNDEFINE  TagDataType = 7  // 8-bit byte
	SSHORT    TagDataType = 8  // 16-bit signed integer
	SLONG     TagDataType = 9  // 32-bit signed integer
	SRATIONAL TagDataType = 10 // Two 32-bit signed integers
	FLOAT     TagDataType = 11 // 4-byte single-precision IEEE floating-point value
	DOUBLE    TagDataType = 12 // 8-byte double-precision IEEE floating-point value
)

func (tdt TagDataType) String() string {
	switch tdt {
	case BYTE:
		return "BYTE"
	case ASCII:
		return "ASCII"
	case SHORT:
		return "SHORT"
	case LONG:
		return "LONG"
	case RATIONAL:
		return "RATIONAL"
	case SBYTE:
		return "SBYTE"
	case UNDEFINE:
		return "UNDEFINE"
	case SSHORT:
		return "SSHORT"
	case SLONG:
		return "SLONG"
	case SRATIONAL:
		return "SRATIONAL"
	case FLOAT:
		return "FLOAT"
	case DOUBLE:
		return "DOUBLE"
	}
	return fmt.Sprintf("unknown(%d)", tdt)
}

type IFD struct {
	NrTags          uint16
	TagData         []Tag
	OffsetToNextIFD uint32
}

type Tag struct {
	TagID              TagID
	TagDataType        TagDataType
	NrValues           uint32
	DataOrOffsetToData uint32
}

func ReadTag(rawTagData []byte, byteReader binary.ByteOrder) Tag {
	tagId := TagID(byteReader.Uint16(rawTagData[:2]))
	tagDataType := TagDataType(byteReader.Uint16(rawTagData[2:4]))
	nrValues := byteReader.Uint32(rawTagData[4:8])
	pointerToTagData := byteReader.Uint32(rawTagData[8:12])
	tag := Tag{tagId, tagDataType, nrValues, pointerToTagData}
	return tag
}

func ReadIFD(rawData []byte, byteReader binary.ByteOrder) IFD {
	nrTags := byteReader.Uint16(rawData[:2])

	var currentPosition = 2
	tags := []Tag{}
	for i := 0; i < int(nrTags); i++ {
		rawTagData := rawData[currentPosition : currentPosition+12]
		tag := ReadTag(rawTagData, byteReader)
		tags = append(tags, tag)
		currentPosition += 12
	}

	offsetToNextIFD := byteReader.Uint32(rawData[currentPosition : currentPosition+4])

	ifd := IFD{nrTags, tags, offsetToNextIFD}
	return ifd
}

func ReadIFDs(rawData []byte, offsetToFirstIFD uint32, byteReader binary.ByteOrder) []IFD {
	ifds := []IFD{}
	var currentPosition = offsetToFirstIFD
	for {
		ifd := ReadIFD(rawData[currentPosition:], byteReader)
		ifds = append(ifds, ifd)
		if ifd.OffsetToNextIFD == 0 {
			break
		}
		currentPosition = ifd.OffsetToNextIFD
	}
	return ifds
}
