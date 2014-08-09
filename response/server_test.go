package response

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"crypto/md5"
	"encoding/base64"
	"time"
	//"strings"
  . "github.com/cgentry/gvac"
)


func setChecksum( w http.ResponseWriter , testData string ) string {
	var sum string

	d := md5.New()
	d.Write( []byte(testData))
	sum = base64.StdEncoding.EncodeToString( d.Sum(nil) )
	w.Header().Set( "Content-MD5" , sum )

	return sum
}

func setDate( w http.ResponseWriter, minutes int ) string {
	offset := time.Duration(minutes) * time.Minute
	stamp := time.Now().UTC().Add( offset ).Format( http.TimeFormat)
	w.Header().Set( GAV_HEADER_TIMESTAMP , stamp )
	return stamp
	
}

/**
 ** Test the date setting functions
 */
func TestDate_Now( t * testing.T ){
	s := NewServer()
	w := httptest.NewRecorder()

	Convey( "Date should be good" , t , func(){
		s.SetDate( w )

		tm  := w.Header().Get( `Timestamp`)
		So( tm , ShouldNotEqual, `` )

		ds, err := s.GetDateString( w )
		So( err, ShouldBeNil )
		So( ds , ShouldEqual, tm )

		tc,err := http.ParseTime( tm )
		So( err , ShouldBeNil )
		So( tc , ShouldNotEqual, 0 )
		dv, ev := s.GetDate( w )
		So( ev , ShouldBeNil )
		So( dv.String(), ShouldEqual, tc.String() )
	})
}

// Check for content signing...
func TestContent_MD5( t * testing.T ){
	s := NewServer()
	w := httptest.NewRecorder()
	body := "This is a test...1234567890"

	m5 := CalculateContentMD5( []byte(body) )

	Convey( "Sign content and verify" , t , func(){
		s.SetContentMD5( w , []byte( body ))
		So( m5, ShouldEqual , w.Header().Get( "Content-MD5"))
	})
}
