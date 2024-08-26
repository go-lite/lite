package lite

type ContentType string

const (
	ContentTypeHTML        ContentType = "text/html"
	ContentTypeCSS         ContentType = "text/css"
	ContentTypeXML         ContentType = "application/xml"
	ContentTypeGIF         ContentType = "image/gif"
	ContentTypeJPEG        ContentType = "image/jpeg"
	ContentTypeJS          ContentType = "text/javascript"
	ContentTypeATOM        ContentType = "application/atom+xml"
	ContentTypeRSS         ContentType = "application/rss+xml"
	ContentTypeMML         ContentType = "text/mathml"
	ContentTypeTXT         ContentType = "text/plain"
	ContentTypeJAD         ContentType = "text/vnd.sun.j2me.app-descriptor"
	ContentTypeWML         ContentType = "text/vnd.wap.wml"
	ContentTypeHTC         ContentType = "text/x-component"
	ContentTypeAVIF        ContentType = "image/avif"
	ContentTypePNG         ContentType = "image/png"
	ContentTypeSVG         ContentType = "image/svg+xml"
	ContentTypeTIFF        ContentType = "image/tiff"
	ContentTypeWBMP        ContentType = "image/vnd.wap.wbmp"
	ContentTypeWEBP        ContentType = "image/webp"
	ContentTypeICO         ContentType = "image/x-icon"
	ContentTypeJNG         ContentType = "image/x-jng"
	ContentTypeBMP         ContentType = "image/x-ms-bmp"
	ContentTypeWOFF        ContentType = "font/woff"
	ContentTypeWOFF2       ContentType = "font/woff2"
	ContentTypeJAR         ContentType = "application/java-archive"
	ContentTypeJSON        ContentType = "application/json"
	ContentTypeHQX         ContentType = "application/mac-binhex40"
	ContentTypeDOC         ContentType = "application/msword"
	ContentTypePDF         ContentType = "application/pdf"
	ContentTypePS          ContentType = "application/postscript"
	ContentTypeRTF         ContentType = "application/rtf"
	ContentTypeM3U8        ContentType = "application/vnd.apple.mpegurl"
	ContentTypeKML         ContentType = "application/vnd.google-earth.kml+xml"
	ContentTypeKMZ         ContentType = "application/vnd.google-earth.kmz"
	ContentTypeXLS         ContentType = "application/vnd.ms-excel"
	ContentTypeEOT         ContentType = "application/vnd.ms-fontobject"
	ContentTypePPT         ContentType = "application/vnd.ms-powerpoint"
	ContentTypeODG         ContentType = "application/vnd.oasis.opendocument.graphics"
	ContentTypeODP         ContentType = "application/vnd.oasis.opendocument.presentation"
	ContentTypeODS         ContentType = "application/vnd.oasis.opendocument.spreadsheet"
	ContentTypeODT         ContentType = "application/vnd.oasis.opendocument.text"
	ContentTypePPTX        ContentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	ContentTypeXLSX        ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypeDOCX        ContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	ContentTypeWMLC        ContentType = "application/vnd.wap.wmlc"
	ContentTypeWASM        ContentType = "application/wasm"
	ContentType7Z          ContentType = "application/x-7z-compressed"
	ContentTypeCCO         ContentType = "application/x-cocoa"
	ContentTypeJARDIFF     ContentType = "application/x-java-archive-diff"
	ContentTypeJNLP        ContentType = "application/x-java-jnlp-file"
	ContentTypeRUN         ContentType = "application/x-makeself"
	ContentTypePL          ContentType = "application/x-perl"
	ContentTypePRC         ContentType = "application/x-pilot"
	ContentTypeRAR         ContentType = "application/x-rar-compressed"
	ContentTypeRPM         ContentType = "application/x-redhat-package-manager"
	ContentTypeSEA         ContentType = "application/x-sea"
	ContentTypeSWF         ContentType = "application/x-shockwave-flash"
	ContentTypeSIT         ContentType = "application/x-stuffit"
	ContentTypeTCL         ContentType = "application/x-tcl"
	ContentTypeCRT         ContentType = "application/x-x509-ca-cert"
	ContentTypeXPI         ContentType = "application/x-xpinstall"
	ContentTypeXHTML       ContentType = "application/xhtml+xml"
	ContentTypeXSPF        ContentType = "application/xspf+xml"
	ContentTypeZIP         ContentType = "application/zip"
	ContentTypeMIDI        ContentType = "audio/midi"
	ContentTypeMP3         ContentType = "audio/mpeg"
	ContentTypeOGG         ContentType = "audio/ogg"
	ContentTypeM4A         ContentType = "audio/x-m4a"
	ContentTypeRA          ContentType = "audio/x-realaudio"
	ContentType3GP         ContentType = "video/3gpp"
	ContentTypeTS          ContentType = "video/mp2t"
	ContentTypeMP4         ContentType = "video/mp4"
	ContentTypeMPEG        ContentType = "video/mpeg"
	ContentTypeMOV         ContentType = "video/quicktime"
	ContentTypeWEBM        ContentType = "video/webm"
	ContentTypeFLV         ContentType = "video/x-flv"
	ContentTypeM4V         ContentType = "video/x-m4v"
	ContentTypeMNG         ContentType = "video/x-mng"
	ContentTypeASX         ContentType = "video/x-ms-asf"
	ContentTypeWMV         ContentType = "video/x-ms-wmv"
	ContentTypeAVI         ContentType = "video/x-msvideo"
	ContentTypeXFormData   ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData    ContentType = "multipart/form-data"
	ContentTypeOctetStream ContentType = "application/octet-stream"
)