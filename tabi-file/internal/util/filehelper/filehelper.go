package filehelper

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// LastString to get last string
func LastString(typeString string) string {
	splitString := strings.Split(typeString, "/")
	return splitString[len(splitString)-1]
}

// Empty empty()
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// InArray in_array()
// haystack supported types: slice, array or map
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type must be slice, array or map")
	}

	return false
}

// IsOverType check file type
func IsOverType(fileType string) bool {
	// check file type, detectcontenttype only needs the first 512 bytes

	avoidFileExtension := "^(bat|exe|cmd|sh|php([0-9])?|pl|cgi|386|dll|com|torrent|app|jar|pif|vb|vbscript|asp|cer|csr|jsp|drv|sys|ade|adp|bas|chm|cpl|crt|csh|fxp|hlp|hta|inf|ins|isp|jse|htaccess|htpasswd|ksh|lnk|mdb|mde|mdt|mdw|msc|msi|msp|mst|ops|pcd|prg|reg|scr|sct|shb|shs|url|vbe|vbs|wsc|wsf|wsh)$"

	//allowedExts := []string{"jpg", "jpeg", "gif", "png", "mov","wmv","wav","amr","htm", "mp4", "m4a", "3gp", "ogg", "csv", "doc", "docx", "txt", "xls", "xlsx","pdf","rar","zip","7up","eml"}

	contentType := LastString(fileType)

	if Empty(contentType) {
		return true
	}

	match, _ := regexp.MatchString(avoidFileExtension, strings.ToLower(contentType))

	return match

}

// AvoidCharacters avoidCharacters()
func AvoidCharacters(fileName string) bool {
	avoidCharacters := []string{"//", "\\", "{", "}", "^", "%", "`", "[", "]", "<", ">", "~", "#", "#"}
	return InArray(fileName, avoidCharacters)
}

// GetUUID get uuid
func GetUUID() string {
	return uuid.New().String()
}
