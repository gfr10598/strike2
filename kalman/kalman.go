/*

 */
package kalman

import (
	"fmt"

	"github.com/ChristopherRabotin/gokalman"
	"github.com/gonum/matrix/mat64"
)

func foobar() {
	F := mat64.NewDense(4, 4, []float64{1, 0.01, 0.00005, 0, 0, 1, 0.01, 0, 0, 0, 1, 0, 0, 0, 0, 1.0005125020836})
	G := mat64.NewDense(4, 1, []float64{0.0, 0.0001, 0.01, 0.0})
	// Note that we will be using two difference H matrices, which we'll swap on the fly.
	//	H1 := mat64.NewDense(2, 4, []float64{1, 0, 0, 0, 0, 0, 1, 1})
	H2 := mat64.NewDense(1, 4, []float64{0, 0, 1, 1})
	// Noise
	Q := mat64.NewSymDense(4, []float64{0.0000000000025, 0.000000000625, 0.000000083333333, 0, 0.000000000625, 0.000000166666667, 0.000025, 0, 0.000000083333333, 0.000025, 0.005, 0, 0, 0, 0, 0.530265088355421})
	Q.ScaleSym(1e-3, Q)
	R := mat64.NewSymDense(2, []float64{0.5, 0, 0, 0.05})
	//	noise1 := gokalman.NewNoiseless(Q, R)
	Ra := mat64.NewSymDense(1, []float64{0.05})
	noise2 := gokalman.NewNoiseless(Q, Ra)

	// Vanilla KF
	x0 := mat64.NewVector(4, []float64{0, 0.45, 0, 0.09})
	Covar0 := gokalman.ScaledIdentity(4, 10)
	vanillaKF, vest0, err := gokalman.NewVanilla(x0, Covar0, F, G, H2, noise2)
	fmt.Printf("Vanilla: \n%s", vanillaKF)
	if err != nil {
		panic(err)
	}
	vanillaEstChan := make(chan (gokalman.Estimate), 1)
	vanillaEstChan <- vest0
}
