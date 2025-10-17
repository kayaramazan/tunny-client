// Package common provides wire protocol helpers for serializing HTTP requests over yamux streams.
package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
)

// WriteHTTPRequestHead writes the request line + headers as bytes with a length prefix.
// It returns (#bytes including CRLF CRLF) written (excluding the 8-byte length prefix).
func WriteHTTPRequestHead(dst io.Writer, r *http.Request) (int64, error) {
	var buf bytes.Buffer
	// Request-Line: METHOD SP REQUEST-URI SP HTTP/VERSION CRLF
	uri := r.URL.RequestURI()
	if uri == "" {
		uri = "/"
	}
	fmt.Fprintf(&buf, "%s %s HTTP/%d.%d\r\n", r.Method, uri, r.ProtoMajor, r.ProtoMinor)
	// Copy headers
	r.Header.Write(&buf)
	// Required CRLF
	buf.WriteString("\r\n")

	// Length prefix
	if err := binary.Write(dst, binary.BigEndian, uint64(buf.Len())); err != nil {
		return 0, err
	}
	n, err := dst.Write(buf.Bytes())
	return int64(n), err
}

// ReadHTTPRequestHead reads a length-prefixed head and returns a reader
// positioned after the head, containing the head bytes for http.ReadRequest/Response.
func ReadHTTPRequestHead(src io.Reader) (*bufio.Reader, error) {
	var n uint64
	if err := binary.Read(src, binary.BigEndian, &n); err != nil {
		return nil, err
	}
	if n > 10<<20 { // 10MB guard
		return nil, fmt.Errorf("head too large")
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(src, b); err != nil {
		return nil, err
	}
	return bufio.NewReader(bytes.NewReader(b)), nil
}
