// The Simple package uses a simple filter to estimate the bell dynamics.
package simple

import "math"

const (
	α_omega = 0.8  // confidence in estimate vs measured rate.
	α_I     = 0.99 // confidence in estimate of bell torque constant
)

type State struct {
	I float64 // The estimated acceleration constant
	θ float64 // bell angle
	ω float64 // bell angular rate
}

// This will be used to estimate the bottom dead center (BDC),
// based on the point where acceleration is zero.  This will
// be slightly off, because the influence of the rope is assymetric.
//
// We will estimate ω' by simple d/dt of the raw rate measurement.

type RollingLinEst struct {
	x                        []float64
	y                        []float64
	xsum, ysum, xysum, xxsum float64
	index                    int32
}

func (est *RollingLinEst) Add(x, y float64) {

}

func (est *RollingLinEst) Slope() {

}

func (state State) Advance(Δt float64, ω float64) {
	estω := state.ω - state.I*Δt*math.Cos(state.θ)
	ω = α_omega*estω + (1-α_omega)*ω
	state.θ = state.θ + Δt*(state.ω+ω)/2

	// For the bell torque constant, we update it only when the bell
	// is moving faster than 3 radians/sec.  This makes it unlikely
	// that it is influenced by rope force.
	// Furthermore, rope force will almost always distort the constant
	// upwards, so we should actually track some low percentile value,
	// perhaps 10th percentile.
	if math.Abs(state.ω) > 3 {

	}
}
