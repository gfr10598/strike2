// The Simple package uses a simple filter to estimate the bell dynamics.
package simple

import (
	"fmt"
	"math"
	"os"
)

// This will be used to estimate the bottom dead center (BDC),
// based on the point where acceleration is zero.  This will
// be slightly off, because the influence of the rope is assymetric.
//
// We will estimate ω' by simple linear fit to d/dt of the raw
// rate measurement.
type RollingLinEst struct {
	count                    int // count of entries
	index                    int // index of last value
	x                        []float64
	y                        []float64
	xsum, ysum, xysum, xxsum float64
}

func NewRollingLinEst(size int) RollingLinEst {
	result := RollingLinEst{}
	result.x = make([]float64, size)
	result.y = make([]float64, size)
	return result
}

// How much trouble will we have with roundoff error drift?
// Could we address this with a very small decay?
func (est *RollingLinEst) Add(x, y float64) {
	if est.count >= len(est.x) {
		x := est.x[est.index]
		y := est.y[est.index]
		est.xsum -= x
		est.ysum -= y
		est.xxsum -= x * x
		est.xysum -= x * y
	} else {
		est.count++
	}
	est.x[est.index] = x
	est.y[est.index] = y
	est.index = (est.index + 1) % len(est.x)
	est.xsum += x
	est.ysum += y
	est.xxsum += x * x
	est.xysum += x * y
}

func (est *RollingLinEst) Slope() float64 {
	n := float64(est.count)
	if n < 1.0 {
		return 0.0
	}
	sx := est.xsum
	sy := est.ysum
	sxy := est.xysum
	sxx := est.xxsum
	return (n*sxy - sx*sy) / (n*sxx - sx*sx)
}

func (est *RollingLinEst) Estimate(x float64) float64 {
	n := float64(est.count)
	if n < 1.0 {
		return 0.0
	}
	var sx, sy, sxx, sxy float64
	for i := 0; i < est.count; i++ {
		x := est.x[i]
		y := est.y[i]
		sx += x
		sy += y
		sxx += x * x
		sxy += x * y
	}
	return ((n*x-sx)*sxy + (sxx-x*sx)*sy) / (n*sxx - sx*sx)
}

func (est *RollingLinEst) Check(x float64) bool {
	other := NewRollingLinEst(len(est.x))
	for i := 0; i < len(est.x); i++ {
		other.Add(est.x[i], est.y[i])
	}
	result := true
	a := other.Estimate(x)
	b := est.Estimate(x)
	delta := (a - b) / (a + b)
	if math.Abs(delta) > 1e-8 {
		fmt.Println("Estimates don't match ", other.Estimate(x), est.Estimate(x))
		fmt.Println(float64(est.count)*est.xxsum-est.xsum*est.xsum, " ? ",
			float64(other.count)*other.xxsum-other.xsum*other.xsum)
		result = false
	}
	if !result {
		fmt.Printf("%v\n", est)
		fmt.Printf("%v\n", other)
		os.Exit(1)
	}
	return result
}

const (
	α_omega = 0.1  // confidence in estimate vs measured rate.
	α_I     = 0.99 // confidence in estimate of bell torque constant
	α_zero  = 0.9  // attenuation of angle at zero crossings
)

type State struct {
	t      float64       // time in seconds.
	Mg_I   float64       // The estimated acceleration constant, about 12
	estθ   float64       // bell angle
	estω   float64       // bell angular rate
	ω      float64       // last measured angular rate
	linEst RollingLinEst // linear BDC estimator
}

func NewState(t float64) State {
	return State{t, 12, math.Pi, 0, 0, NewRollingLinEst(20)}
}

func (state *State) Angle() float64 {
	return state.estθ * 180 / math.Pi
}

func (state *State) Rate() float64 {
	return state.estω
}

func (state *State) Reset(t float64, ω float64) {
	fmt.Println("Reset", state.estθ, state.estω, t, ω)
	state.estθ = math.Pi
	state.estω = ω
	state.Mg_I = 12
	state.t = t
	state.ω = ω
	state.linEst = NewRollingLinEst(len(state.linEst.x))
}

// TODO - check for bell already up - may need to add or subtract pi.
func (state *State) Advance(t float64, ω float64) {
	dω := ω - state.ω // change in raw rate.
	Δt := t - state.t
	if Δt > 1.0 {
		state.Reset(t, ω)
		return
	}
	if Δt > 0.0 {
		state.linEst.Add(t, dω/Δt)
	}
	// BTW, near BDC, the acceleration curve should be very linear
	// with low noise.  We could use this for conditioning BDC
	// estimation.

	sinθ := math.Sin(state.estθ)
	forward_dω := -state.Mg_I * Δt * sinθ
	forwardω := state.estω + forward_dω

	// For the bell torque constant, we update it only when:
	//   a. |sin(θ)| >> 0
	//   b. |ω| >> 0
	//
	// a. avoids ratios of very small numbers.
	// b. makes it unlikely that it is influenced by rope force.
	//
	// Furthermore, rope force will almost always distort the constant
	// upwards, so we should actually track some low percentile value,
	// perhaps 10th percentile.
	if math.Abs(ω) > 5.0 {
		acc0 := state.linEst.Estimate(state.t)
		acc1 := state.linEst.Estimate(t)
		if (acc0 * acc1) <= 0.0 {
			fmt.Printf("Time: %8.3f Accelerations: %8.3f %8.3f  angle: %6.3f  rate: %6.3f  Slope: %6.3f\n",
				t, acc0, acc1, state.Angle(), state.Rate(), state.linEst.Slope())
			// We have just passed through BDC.
			// If our integrator is far off from BDC, then we should
			// just reset
			if math.Abs(state.estθ) > 1 {
				//	state.Reset(t, ω)
				//	return
			}

			// fmt.Printf("%v\n", state)
			state.estθ *= α_zero
		}
		if math.Abs(sinθ) > .2 {
			apparentI := -dω / (Δt * sinθ)
			if apparentI > state.Mg_I {
				state.Mg_I += .001
			} else {
				state.Mg_I -= .001
			}

			//			state.Mg_I = α_I*state.Mg_I + (1-α_I)*apparentI
		}
	}

	// Finally, update the state variables.
	state.estθ = state.estθ + Δt*(state.estω+forwardω)/2
	state.estω = α_omega*forwardω + (1-α_omega)*ω
	state.ω = ω
	state.t = t

	// Now, if the linest suggests that we are very close to BDC,
	// update the angle estimate accordingly.
}
