package filetype

const TYPE_UNDEFINE = "UNDEFINE"

type Extension = string

const (
	Unkonw Extension = "unknow"
	MSI    Extension = "msi"
	APK    Extension = "apk"

	DOC  Extension = "doc"
	DOCM Extension = "docm"
	DOCX Extension = "docx"
	DOT  Extension = "dot"
	DOTM Extension = "dotm"
	DOTX Extension = "dotx"
	RTF  Extension = "rtf"
	ODT  Extension = "odt"

	PPT  Extension = "ppt"
	PPTM Extension = "pptm"
	PPTX Extension = "pptx"
	POT  Extension = "pot"
	POTM Extension = "potm"
	POTX Extension = "potx"
	PPS  Extension = "pps"
	PPSM Extension = "ppsm"
	PPSX Extension = "ppsx"
	PPA  Extension = "ppa"
	PPAM Extension = "ppam"
	ODP  Extension = "odp"

	XLS  Extension = "xls"
	XLSM Extension = "xlsm"
	XLSX Extension = "xlsx"
	XLSB Extension = "xlsb"
	XLT  Extension = "xlt"
	XLTM Extension = "xltm"
	XLTX Extension = "xltx"
	XLA  Extension = "xla"
	XLAM Extension = "xlam"
	XPS  Extension = "xps"
	SLK  Extension = "slk"
	ODS  Extension = "ods"

	VSD   Extension = "vsd"
	FPX   Extension = "fpx"
	CDFV2 Extension = "cdfv2"

	// wps
	WPS Extension = "wps"
	WPT Extension = "wpt"
	ET  Extension = "et"
	EET Extension = "eet"
	DPS Extension = "dps"
	DPT Extension = "dpt"

	MSG Extension = "msg"
	EML Extension = "eml"

	ZIP Extension = "zip"
	JAR Extension = "jar"
	XML Extension = "xml"
)

var extTypes = map[Extension]string{
	MSI: TYPE_UNDEFINE,
	APK: TYPE_UNDEFINE,

	DOC:  TYPE_UNDEFINE,
	DOCM: TYPE_UNDEFINE,
	DOCX: TYPE_UNDEFINE,
	DOT:  TYPE_UNDEFINE,
	DOTM: TYPE_UNDEFINE,
	DOTX: TYPE_UNDEFINE,
	RTF:  TYPE_UNDEFINE,
	ODT:  TYPE_UNDEFINE,

	PPT:  TYPE_UNDEFINE,
	PPTM: TYPE_UNDEFINE,
	PPTX: TYPE_UNDEFINE,
	POT:  TYPE_UNDEFINE,
	POTM: TYPE_UNDEFINE,
	POTX: TYPE_UNDEFINE,
	PPS:  TYPE_UNDEFINE,
	PPSM: TYPE_UNDEFINE,
	PPSX: TYPE_UNDEFINE,
	PPA:  TYPE_UNDEFINE,
	PPAM: TYPE_UNDEFINE,
	ODP:  TYPE_UNDEFINE,

	XLS:  TYPE_UNDEFINE,
	XLSM: TYPE_UNDEFINE,
	XLSX: TYPE_UNDEFINE,
	XLSB: TYPE_UNDEFINE,
	XLT:  TYPE_UNDEFINE,
	XLTM: TYPE_UNDEFINE,
	XLTX: TYPE_UNDEFINE,
	XLA:  TYPE_UNDEFINE,
	XLAM: TYPE_UNDEFINE,
	XPS:  TYPE_UNDEFINE,
	SLK:  TYPE_UNDEFINE,
	ODS:  TYPE_UNDEFINE,

	VSD:   TYPE_UNDEFINE,
	FPX:   TYPE_UNDEFINE,
	CDFV2: TYPE_UNDEFINE,

	// wps
	WPS: TYPE_UNDEFINE,
	WPT: TYPE_UNDEFINE,
	ET:  TYPE_UNDEFINE,
	EET: TYPE_UNDEFINE,
	DPS: TYPE_UNDEFINE,
	DPT: TYPE_UNDEFINE,

	// mail
	MSG: TYPE_UNDEFINE,
	EML: TYPE_UNDEFINE,

	// archive
	ZIP: TYPE_UNDEFINE,
	JAR: TYPE_UNDEFINE,
	XML: TYPE_UNDEFINE,
}

func ExtClass(ext Extension) string {
	return extTypes[ext]
}
