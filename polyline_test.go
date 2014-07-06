package gmaps

import (
	"fmt"
	"testing"
)

func ExampleDecodePolyline() {
	// Start decoding relative to (0, 0); though underspecified in the API docs,
	// that is how their polyline is set up.
	start := Point{Lat: 0, Lng: 0}

	var ps []Point
	e := DecodePolyline(start, []byte("_p~iF~ps|U_ulLnnqC_mqNvxq`@"), &ps)
	fmt.Printf("Error: %v\nDecoded: %v\n", e, ps)
	// Output:
	// Error: <nil>
  // Decoded: [{0 0} {38.5 -120.2} {40.7 -120.95} {43.252 -126.453}]
}

// Test decoding the sample polyline given on
// https://developers.google.com/maps/documentation/utilities/polylinealgorithm.
func TestSamplePolyline(t *testing.T) {
	polyline := []byte("_p~iF~ps|U_ulLnnqC_mqNvxq`@")

	var ps []Point
	e := DecodePolyline(Point{Lat: 0, Lng: 0}, polyline, &ps)
	if e != nil {
		t.Errorf("Want: nil error, got: %s", e.Error())
	}

	if ps == nil {
		t.Error("Want: ps not nil, got: nil")
	}

	expected := []Point{
		Point{Lat: 0, Lng: 0},
		Point{Lat: 38.5, Lng: -120.2},
		Point{Lat: 40.7, Lng: -120.95},
		Point{Lat: 43.252, Lng: -126.453}}

	if len(expected) != len(ps) {
		t.Errorf("Want: %d-length array, got: %v", len(expected), ps)
	}

	for i, _ := range expected {
		if expected[i] != ps[i] {
			t.Errorf("Element %d: want %v, got %v", i, expected[i], ps[i])
		}
	}
}

// Test decoding the example string above when one of the latitudes is
// cut off.
func TestIncompleteLat(t *testing.T) {
	polyline := []byte("_p~iF~ps|U_ul")

	var ps []Point
	e := DecodePolyline(Point{Lat: 0, Lng: 0}, polyline, &ps)
	if e == nil {
		t.Errorf("Want: non-nil error, got: %s", e.Error())
	}

	if e.Error() != "Incomplete token (3 bytes): _ul (hex 0x5f756c)" {
		t.Errorf("Incorrect error message: %s", e.Error())
	}

	if ps == nil {
		t.Error("Want: ps not nil, got: nil")
	}

	expected := []Point{
		Point{Lat: 0, Lng: 0},
		Point{Lat: 38.5, Lng: -120.2}}

	if len(expected) != len(ps) {
		t.Errorf("Want: %d-length array, got: %v", len(expected), ps)
	}

	for i, _ := range expected {
		if expected[i] != ps[i] {
			t.Errorf("Element %d: want %v, got %v", i, expected[i], ps[i])
		}
	}
}

// Test decoding the example string above when one of the longitudes is
// cut off.
func TestIncompleteLng(t *testing.T) {
	polyline := []byte("_p~iF~ps|U_ulLnnq")

	var ps []Point
	e := DecodePolyline(Point{Lat: 0, Lng: 0}, polyline, &ps)
	if e == nil {
		t.Errorf("Want: non-nil error, got: %s", e.Error())
	}

	if e.Error() != "Incomplete token (3 bytes): nnq (hex 0x6e6e71)" {
		t.Errorf("Incorrect error message: %s", e.Error())
	}

	if ps == nil {
		t.Error("Want: ps not nil, got: nil")
	}

	expected := []Point{
		Point{Lat: 0, Lng: 0},
		Point{Lat: 38.5, Lng: -120.2}}

	if len(expected) != len(ps) {
		t.Errorf("Want: %d-length array, got: %v", len(expected), ps)
	}

	for i, _ := range expected {
		if expected[i] != ps[i] {
			t.Errorf("Element %d: want %v, got %v", i, expected[i], ps[i])
		}
	}
}

func TestEmptyString(t *testing.T) {
	var ps []Point
	e := DecodePolyline(Point{Lat: 0, Lng: 0}, []byte(""), &ps)

	if e != nil {
		t.Errorf("Want no error, got: %v", e)
	}

	if len(ps) != 1 {
		t.Errorf("Want len(ps) == 1, got len(ps) == %v", len(ps))
	}

	if ps[0].Lat != 0 {
		t.Errorf("Want ps[0].Lat == 0, got: %v", ps[0].Lat)
	}

	if ps[0].Lng != 0 {
		t.Errorf("Want ps[0].Lng == 0, got: %v", ps[0].Lng)
	}
}
