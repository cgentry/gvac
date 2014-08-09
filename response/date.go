package response

import (
	"net/http"
	"time"
	"errors"
	. "github.com/cgentry/gvac"
)

/*
 *  Simple function to put a correct time/date stamp into the header
 *	Parm:	Pointer to http.Request to modify
 *  return: pointer to Secure for chaining
 */
func (s * Secure) SetDate(w http.ResponseWriter) * Secure {
	tm := time.Now().UTC().Format(http.TimeFormat)
	w.Header().Set( GAV_HEADER_TIMESTAMP, tm)
	w.Header().Set( GAV_HEADER_DATE     , tm )
	return s
}

func (s * Secure) GetDateString(w http.ResponseWriter) ( string , error ) {
	requestDate := w.Header().Get(GAV_HEADER_TIMESTAMP)        // Header has "Timestamp:"
	if len(requestDate) == 0 {                    			 	// Umm..NO
		requestDate = w.Header().Get(GAV_HEADER_DATE)        	// Header has "Date:" ?

		if len(requestDate) == 0 {
			return "", errors.New(TIMESTAMP_MISSING)
		}
	}
	return requestDate, nil
}

func (s * Secure) GetDate(w http.ResponseWriter) ( tstamp time.Time , err error) {
	var dt string
	if dt, err = s.GetDateString(w); err == nil {
		tstamp, err = http.ParseTime(dt)    // Always use this parser as it does 3 formats...
	}
	return
}

