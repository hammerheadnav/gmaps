// Package gmaps provides a utility to decode Google Maps binary-encoded
// polylines.
package gmaps

import (
	"fmt"
)

// Represents a single point stored in a polyline.
type Point struct {
	Lat float64
	Lng float64
}

// An error returned when parsing a polyline that contains an invalid token.
type IncompleteTokenError struct {
	token []byte
}

func (e IncompleteTokenError) Error() string {
	return fmt.Sprintf("Incomplete token (%d bytes): %s (hex 0x%x)", len(e.token), e.token, e.token)
}

// Decode a single token from the given byte string. A token represents a
// float64 and corresponds to a contiguous group of bytes, where all bytes
// except the trailing byte have the 0x20 bit set. Returns pos, the index of the
// character after the final decoded byte; v, the decoded value; and err, either
// nil if decoding succeeded or an IncompleteTokenError if no byte without 0x20
// was found in line.
func decodeOneToken(line []byte, precision float64) (pos int, v float64, err error) {
	var token int64
	var shift uint
	for i, c := range line {
		d := int64(c) - 63
		token |= (d & 0x1f) << uint(shift)
		shift += 5
		if d&0x20 == 0 {
			pos = i + 1
			err = nil
			if token&1 != 0 {
				token = ^(token >> 1)
			} else {
				token = token >> 1
			}
			v = float64(token) / precision
			return
		}
	}

	err = IncompleteTokenError{token: line}
	pos = 0
	return
}

// Parse the given polyline, relative to a starting point (lat, lng),
// and append the starting point and all of the points in the polyline
// to p. The polyline string stores the difference between each
// point's coordinates, but p is filled with the actual points that make
// up the polyline. If p is nil, a new array is instantiated. If the polyline
// string is truncated, err is an IncompleteTokenError; otherwise it is nil.
func DecodePolyline(start Point, line []byte, p *[]Point, precision float64) (err error) {
	if *p == nil {
		*p = make([]Point, 0)
	}

	var latDelta, lngDelta float64
	var pos, consumed int
	for pos < len(line) {
		consumed, latDelta, err = decodeOneToken(line[pos:len(line)], precision)
		pos += consumed
		start.Lat += latDelta
		if err != nil {
			return
		}

		consumed, lngDelta, err = decodeOneToken(line[pos:len(line)], precision)
		pos += consumed
		start.Lng += lngDelta
		if err != nil {
			return
		}

		*p = append(*p, start)
	}
	return
}

func DecodePolylineWithoutStartingPoint(line []byte, p *[]Point, precision float64) (err error) {
	return DecodePolyline(Point{}, line, p, precision)
}
