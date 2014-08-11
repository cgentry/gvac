gvac
===

Go-language Verify Access for web requests and responces.

Similar to AWS hmac signing, this will sign and verify requests using:

1. A shared secret
2. A user id
3. Content-MD5 (if it exists)
4. Content-Type
5. Date/time stamp
6. The body of the message
7. Any specific header strings
8. The "action" code. This is normally the path
9. All the parameters in the request.
10. The fragment (or options)

If there is a body text, there must be a Content-MD5 header value set. This is compared to one created on the locally to ensure validity.

The Date/time stamp can either be a header field 'Timestamp' or the 'Date' header value. I recommend that every program set the 'Timestamp' value
instead of using the 'Date' field. This should be formatted in GMT and in standard HTTP header Date format.


Your code must either use the routines in this library (if using GO) or recreate the exact way the key is generated.

## The action, parameters and fragments:
Given a request string of: /gus/register?user=webapp&domain=local#nopush

* The action is **/gus/register**
* The parameters string would be **user=webapp&domain=local** Notice that the string appears in key-sorted order.
* The fragment is **nopush**

The string is "rebuilt" to look just like the original request.
