package request

import (
	"net/http"
	"errors"
  .  "github.com/cgentry/gvac"
)


func ( s * Secure ) VerifyContentMD5( r * http.Request , body []byte) error {
	calcMD5 := CalculateContentMD5( body )
	if calcMD5 != "" {
		contentMD5 := r.Header.Get( GAV_HEADER_MD5 )
		if contentMD5 == "" {
			return errors.New( MD5_MISSING )
		}
		if calcMD5 != contentMD5 {
			return errors.New( MD5_MISMATCH )
		}
	}
	return nil
}


func ( s * Secure ) SetContentMD5( r * http.Request , body [] byte ) * Secure  {
	r.Header.Set( GAV_HEADER_MD5 , CalculateContentMD5( body ))
	return s
}


