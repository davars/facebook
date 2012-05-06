## Overview

facebook.go facilitates writing Go programs that integrate with Facebook.  Also, this
README is inaccurate and describes features not yet implemented.

## Documentation

* [Authentication](http://developers.facebook.com/docs/authentication/)
* [Graph API](http://developers.facebook.com/docs/reference/api/)

## Installation 

To install facebook.go, simply run `go get github.com/davars/facebook`. 
To use it in a program, use `import "github.com/davars/facebook"`

## Usage

If you haven't done so already, you should create a [Facebook
Application](https://developers.facebook.com/apps).

You'll have to define an Application instance somehow.  You can refer to the
getTestApp function in facebook_test.go for one possible implementation.

The main API mimics that of the net/http package, except that the Request
always goes to the Graph API url and optionally includes an access token.
Also, since Graph API are JSON, all responses are parsed into a
map[string]interface{}. 

The facebook/oauth package provides tools for handling access_tokens and
signed_requests.

The facebook/test_users package provides tools for manipulating test user
accounts.

## License

Copyright (c) 2012 David Jack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
