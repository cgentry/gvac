package request

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
	. "github.com/cgentry/gvac"
)

/*
 *  Simple function to put a correct time/date stamp into the header
 *	Parm:	Pointer to http.Request to modify
 *  return: pointer to Secure for chaining
 */
func (s * Secure) SetDate(r * http.Request) * Secure {
	tm := time.Now().UTC().Format(http.TimeFormat)
	r.Header.Set( GAV_HEADER_TIMESTAMP, tm)
	r.Header.Set( GAV_HEADER_DATE     , tm )
	return s
}

func (s * Secure) GetDateString(r * http.Request) ( string , error ) {
	requestDate := r.Header.Get(GAV_HEADER_TIMESTAMP)        // Header has "Timestamp:"
	if len(requestDate) == 0 {                    // Umm..NO
		requestDate = r.Header.Get(GAV_HEADER_DATE)        // Header has "Date:" ?

		if len(requestDate) == 0 {
			return "", errors.New(TIMESTAMP_MISSING)
		}
	}
	return requestDate, nil
}

func (s * Secure) GetDate(r * http.Request) ( tstamp time.Time , err error) {
	var requestDate string
	if requestDate, err = s.GetDateString(r); err == nil {
		// Check to see if timestamp is older than 15min. If so, reject request
		// First, parse this into a time object...
		tstamp, err = http.ParseTime(requestDate)    // Always use this parser as it does 3 formats...
	}
	return
}

/*
 *  Get Signature date from header unless it falls outside of the maximum time slots
 *  If the date is within the range, then you get a string. If it is outside of range
 *  or doesn't exist, you get an error return
 *	Parm:	Pointer to http.Request to modify
 *  return: error message or nil on no error
 */
func (s * Secure) VerifyDate(r *http.Request)  error  {

	tstamp, err := s.GetDate(r)
	if err != nil {
		return err
	}
	// Now...see what the difference is between NOW and the HTTP date
	diff := math.Abs(time.Now().Sub(tstamp).Minutes())        // We want how far in the past it is...

	if diff > s.TimeWindow.Minutes() {
		return fmt.Errorf("%s - %.0f min. max/%.0f in header",
			TIMESTAMP_RANGE, s.TimeWindow.Minutes(), diff)
	}
	return nil                        // Passed all tests...
}
