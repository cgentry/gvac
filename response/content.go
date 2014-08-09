package response

import (
	"net/http"
	. "github.com/cgentry/gvac"
)

/*
 * Set the Content-MD5 header
 */
func ( s * Secure ) SetContentMD5( w  http.ResponseWriter , body [] byte ) * Secure  {
	w.Header().Set( GAV_HEADER_MD5 ,CalculateContentMD5(body))
	return s
}

