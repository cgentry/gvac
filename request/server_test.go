package request

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	//"github.com/cgentry/gus/record"
	//"bytes"
	"net/http"
	"crypto/md5"
	//"crypto/hmac"
	//"crypto/sha256"
	"encoding/base64"
	//"fmt"
	"time"
	"strings"
  . "github.com/cgentry/gvac"
)

func setChecksum( r * http.Request , testData string ) string {
	var sum string

	d := md5.New()
	d.Write( []byte(testData))
	sum = base64.StdEncoding.EncodeToString( d.Sum(nil) )
	r.Header.Set( "Content-MD5" , sum )

	return sum
}

func setDate( r * http.Request, minutes int ) string {
	offset := time.Duration(minutes) * time.Minute
	stamp := time.Now().UTC().Add( offset ).Format( http.TimeFormat)
	r.Header.Set( GAV_HEADER_TIMESTAMP , stamp )
	return stamp
	
}

func TestDate_Now( t * testing.T ){
	s := NewServer()
	testData := "Test Body should be here"
	r, err := http.NewRequest( "POST" , "http://example.com/test?a=b&c=d#fragment" , strings.NewReader( testData ) )

	Convey( "Date should be good" , t , func(){
		So( err, ShouldBeNil )
		stamp := setDate( r , 0 )
		err:=s.VerifyDate( r )
		So( err, ShouldBeNil )
		val, err := s.GetDateString(r)
		So( val ,ShouldEqual, stamp  )
		So( err , ShouldBeNil )
	})
}

func TestDate_FutureBad( t * testing.T ){
	s := NewServer()
	testData := "Test Body should be here"
	r, err := http.NewRequest( "POST" , "http://example.com/test?a=b&c=d#fragment" , strings.NewReader( testData ) )

	Convey( "Date should be outside of range" , t , func(){
		So( err, ShouldBeNil )
		stamp := setDate( r , 20 )
		err:=s.VerifyDate( r )
		So( err, ShouldNotBeNil )
		val, err := s.GetDateString(r)
		So( val ,ShouldEqual, stamp  )
		So( err , ShouldBeNil )
	})
}
func TestDate_PastBad( t * testing.T ){
	s := NewServer()
	testData := "Test Body should be here"
	r, err := http.NewRequest( "POST" , "http://example.com/test?a=b&c=d#fragment" , strings.NewReader( testData ) )

	Convey( "Date should be outside of range" , t , func(){
		So( err, ShouldBeNil )
		stamp := setDate( r , -20 )
		err:=s.VerifyDate( r )
		So( err, ShouldNotBeNil )
		So( err.Error() , ShouldStartWith , TIMESTAMP_RANGE)

		val, err := s.GetDateString(r)
		So( val ,ShouldEqual, stamp  )
		So( err , ShouldBeNil )

	})
}

func TestDate_FutureOK( t * testing.T ){
	s := NewServer()
	testData := "Test Body should be here"
	r, err := http.NewRequest( "POST" , "http://example.com/test?a=b&c=d#fragment" , strings.NewReader( testData ) )

	Convey( "Date should be outside of range" , t , func(){
		So( err, ShouldBeNil )
		stamp := setDate( r , +14 )
		err:=s.VerifyDate( r )
		So( err, ShouldBeNil )
		val, err := s.GetDateString(r)
		So( val ,ShouldEqual, stamp  )
		So( err , ShouldBeNil )
	})
}

func TestDate_PastOK( t * testing.T ){
	s := NewServer()
	testData := "Test Body should be here"
	r, err := http.NewRequest( "POST" , "http://example.com/test?a=b&c=d#fragment" , strings.NewReader( testData ) )

	Convey( "Date should be outside of range" , t , func(){
		So( err, ShouldBeNil )
		stamp := setDate( r , -14 )
		err:=s.VerifyDate( r )
		So( err, ShouldBeNil )
		val, err := s.GetDateString(r)
		So( val ,ShouldEqual, stamp  )
		So( err , ShouldBeNil )
	})
}

func TestComputeBodyMd5_Good_Sum( t * testing.T ){
	testData := "Test Body should be here"
	r, err := http.NewRequest( "POST" , "http://example.com/test?a=b&c=d#fragment" , strings.NewReader( testData ) )

	Convey( "MD5 values should be the same", t , func(){
		So( err, ShouldBeNil )
		sum := setChecksum( r , testData )
		So( sum , ShouldEqual , CalculateContentMD5( []byte( testData ) ))
	})
}

// When we don't have a body, we should get a blank string back.

func TestComputeBodyMd5_NoBody( t * testing.T ){

	testData := ""

	Convey( "MD5 values should be blank", t , func(){
		So( CalculateContentMD5( []byte( testData) ) , ShouldEqual , "")
	})
}

func TestGetUser_WithValue( t * testing.T ){
	s := NewServer()
	testData := ""
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "No blanks, Should be 12345", t , func(){
		So( err, ShouldBeNil )
		r.Header.Set( "Authorization" , "12345:abcde")
		val , err := s.GetUser( r )
		So( err, ShouldBeNil )
		So( val , ShouldEqual , "12345")
	})
}

func TestGetUser_WithLeadingBlanks( t * testing.T ){
	s := NewServer()
	testData := ""
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "Trim blanks from user", t , func(){
		So( err, ShouldBeNil )
		r.Header.Set( "Authorization" , " 12345 : abcde ")
		val , err := s.GetUser( r )
		So( err, ShouldBeNil )
		So( val , ShouldEqual , "12345")
	})
}

func TestGetUser_NoAuth( t * testing.T ){
	s := NewServer()
	testData := ""
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "Empty string returns error", t , func(){
		So( err, ShouldBeNil )
		val , err := s.GetUser( r )
		So( err, ShouldNotBeNil )
		So( val, ShouldEqual, "" )
	})
}


func TestGetSignature_Simple( t * testing.T ){
	s := NewServer()
	testData := ""
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "No blanks in signature, Should be abcde", t , func(){
		So( err, ShouldBeNil )
		r.Header.Set( "Authorization" , "12345:abcde")
		val , err := s.GetSignature( r )
		So( err, ShouldBeNil )
		So( val , ShouldEqual , "abcde")
	})
}

func TestGetSignature_WithBlanks( t * testing.T ){
	s := NewServer()
	testData := ""
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "Trim blanks from signature", t , func(){
		So( err, ShouldBeNil )
		r.Header.Set( "Authorization" , " 12345 : abcde ")
		val , err := s.GetSignature( r )
		So( err, ShouldBeNil )
		So( val , ShouldEqual , "abcde")
	})
}

func TestGetSignature_NoAuth( t * testing.T ){
	s := NewServer()
	testData := ""
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "Empty string returns error", t , func(){
		So( err, ShouldBeNil )
		val , err := s.GetSignature( r )
		So( err, ShouldNotBeNil )
		So( val, ShouldEqual, "" )
	})
}

func TestSignRequestAndTest( t * testing.T ){
	var secret = []byte( `abcde`)
	var userId = "123"
	s := NewServer()
	testData := "Good morning world!"
	r, err := http.NewRequest( "POST" , "http://example.com/test" , strings.NewReader( testData ) )

	Convey( "Signature should be valid", t , func(){
		So( err, ShouldBeNil )
		err = s.SetSignature( r , userId , secret, []byte( testData) )

		So( err , ShouldBeNil )
		val , err := s.GetSignature( r )
		So( err, ShouldBeNil )
		So( val, ShouldNotEqual, "" )

		So( s.VerifyDate( r ), ShouldBeNil )
		So( s.VerifyContentMD5( r, []byte( testData )), ShouldBeNil )

		So( s.VerifySignature( r , secret , []byte( testData)) , ShouldBeNil)
		val,err = s.GetUser(r)
		So( err, ShouldBeNil )
		So( val , ShouldEqual, userId )

		So( r.Header.Get( GAV_HEADER_MD5) , ShouldEqual, CalculateContentMD5( []byte( testData )))
		So( r.Header.Get( GAV_HEADER_TOKEN ) , ShouldStartWith, userId )
		dt, err := s.GetDateString(r)
		So( r.Header.Get(GAV_HEADER_TIMESTAMP ) , ShouldEqual, dt )
		So( err , ShouldBeNil )
	})
}
