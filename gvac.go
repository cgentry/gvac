package gvac

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"strings"
	"sort"
)

const GAV_HEADER_TIMESTAMP	= "Timestamp"
const GAV_HEADER_DATE		= "Date"
const GAV_HEADER_TOKEN		= "Authorization"
const GAV_HEADER_MD5		= "Content-MD5"
const GAV_HEADER_TYPE		= "Content-Type"


/**
 * 	Return the base64 of the MD5 of the body.
 *  if the body is empty, you will receive an empty string
 */
func CalculateContentMD5(body []byte) string {
	var sum string = ""
	if len(body) > 0 {
		d := md5.New()
		d.Write(body)
		m5 := d.Sum(nil)
		sum = base64.StdEncoding.EncodeToString(m5)
	}
	return sum
}


/**
 * Passed a Header Map, this will return a byte array of all matching
 * header keys in sorted order. So, if you pass the prefix 'X-GAV-'
 * and the header contains:
 *		X-GAV-BETA:		Value
 *		X-GAV-ALPHA:	Another
 *
 * You get back:
 *		X-GAV-ALPHAAnotherX-GAV-BETAValue
 */
func GetAppHeaderValues( h http.Header , prefix string  ) []byte {
	rtn := ""
	if prefix == "" {
		var headerList []string
		for key,_ := range h {
			if strings.HasPrefix( key , prefix ){
				headerList = append( headerList , key )
			}
		}
		sort.Strings( headerList )
		for _,k := range headerList {
			rtn = rtn + k + h.Get( k )
		}
	}
	return []byte( rtn )
}
