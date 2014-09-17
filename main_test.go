package main

import (
	"bytes"
	"math"
	"testing"

	pb "github.com/dgryski/carbonzipper/carbonzipperpb"
)

func TestInterval(t *testing.T) {

	var tests = []struct {
		t       string
		seconds int32
		sign    int
	}{
		{"1s", 1, 1},
		{"2d", 2 * 60 * 60 * 24, 1},
		{"10hours", 60 * 60 * 10, 1},

		{"1s", -1, -1},
		{"+2d", 2 * 60 * 60 * 24, -1},
		{"-10hours", -60 * 60 * 10, -1},
	}

	for _, tt := range tests {
		if secs, _ := intervalString(tt.t, tt.sign); secs != tt.seconds {
			t.Errorf("intervalString(%q)=%d, want %d\n", tt.t, secs, tt.seconds)
		}
	}
}

func TestJSONResponse(t *testing.T) {

	tests := []struct {
		results []*pb.FetchResponse
		out     []byte
	}{
		{
			[]*pb.FetchResponse{
				makeResponse("metric1", []float64{1, 1.5, 2.25, math.NaN()}, 100, 100),
				makeResponse("metric2", []float64{2, 2.5, 3.25, 4, 5}, 100, 100),
			},
			[]byte(`[{"target":"metric1","datapoints":[[1,100],[1.5,200],[2.25,300],[null,400]]},{"target":"metric2","datapoints":[[2,100],[2.5,200],[3.25,300],[4,400],[5,500]]}]`),
		},
	}

	for _, tt := range tests {
		b := marshalJSON(tt.results)
		if !bytes.Equal(b, tt.out) {
			t.Errorf("marshalJSON(%+v)=%+v, want %+v", tt.results, string(b), string(tt.out))
		}
	}
}

func TestRawResponse(t *testing.T) {

	tests := []struct {
		results []*pb.FetchResponse
		out     []byte
	}{
		{
			[]*pb.FetchResponse{
				makeResponse("metric1", []float64{1, 1.5, 2.25, math.NaN()}, 100, 100),
				makeResponse("metric2", []float64{2, 2.5, 3.25, 4, 5}, 100, 100),
			},
			[]byte(`metric1,100,500,100|1,1.5,2.25,None` + "\n" + `metric2,100,600,100|2,2.5,3.25,4,5` + "\n"),
		},
	}

	for _, tt := range tests {
		b := marshalRaw(tt.results)
		if !bytes.Equal(b, tt.out) {
			t.Errorf("marshalRaw(%+v)=%+v, want %+v", tt.results, string(b), string(tt.out))
		}
	}
}
