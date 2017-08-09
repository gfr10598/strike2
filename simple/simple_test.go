package simple_test

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gfr10598/strike2/simple"
)

func TestLinEst(t *testing.T) {
	lin := simple.NewRollingLinEst(20)
	lin.Add(1.0, 1.0)
	lin.Add(2.0, 2.0)
	lin.Add(3.0, 3.0)
	lin.Add(4.0, 4.0)
	if lin.Slope() != 1.0 {
		t.Error("Bad slope: ", lin.Slope())
	}
	if lin.Estimate(5.0) != 5.0 {
		t.Error("Bad estimate: ", lin.Estimate(5.0))
	}

	for i := 0; i < 50000; i++ {
		lin.Add(float64(i)+rand.Float64(), float64(i)+rand.Float64())
	}
	if !lin.Check(50000.0) {
		t.Fatal()
	}
}

func TestBellUp(t *testing.T) {
	f, err := os.Open("testdata/sensorLog.txt")
	if err != nil {
		t.Fatal()
	}
	rdr := bufio.NewReader(f)

	state := simple.NewState(0.0)
	last_time := 0
	for line, _, err := rdr.ReadLine(); err != io.EOF; line, _, err = rdr.ReadLine() {
		words := strings.Fields(string(line))
		if len(words) < 2 || words[1] != "GYR" {
			continue
		}
		tt, _ := strconv.Atoi(words[0])
		time := .001 * float64(tt)
		w, _ := strconv.ParseFloat(words[4], 64)
		state.Advance(time, w)
		if int(10*time) > last_time {
			fmt.Printf("%10.3f %5.1f, %6.2f, %6.2f, %6.2f\n", time, state.Angle(), state.Rate(), w, state.Mg_I)
		}
		last_time = int(10 * time)
	}
	fmt.Printf("%v\n", state)
}
