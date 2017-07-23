package main

import (
	"fmt"
	"net"
	"os"
)

/* A Simple function to verify error */
func Fail(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

type Point struct {
	t float64
	rates [3]float64
}

type LinearEstimate struct {
	count float64
	sum_x, sum_y float64
	sum_xx, sum_xy float64
}

func (est *LinearEstimate) Add(x,y float64) {
	est.count++
	est.sum_x += x
    est.sum_y += y
    est.sum_xx += x * x
    est.sum_xy += x * y
}

// MoveX adjusts the state as if all data had been provided at dx greater X values.
func (est *LinearEstimate) MoveX(dx float64) {
	// order matters
	est.sum_xx += est.count * dx * dx + 2 * dx * est.sum_x
	est.sum_x += dx * est.count
    est.sum_xy += dx * est.sum_y
}

func (est *LinearEstimate) Estimate(x float64) float64 {
	n := est.count
	sx := est.sum_x
	sy := est.sum_y
	sxy := est.sum_xy
	sxx := est.sum_xx
	return ((n * x - sx) * sxy + (sxx - x * sx) * sy) / (n * sxx - sx * sx)
}

// Ideally this should be a kalman filter
type State struct {
	points []Point // last N points
	linEst LinearEstimate
}

func main() {
	
	var est2 LinearEstimate
	est2.Add(2, 1)
	est2.Add(3, 3)
	est2.Add(4, 5)
	fmt.Printf("%v\n", est2)
	fmt.Println("est(5) = ", est2.Estimate(5))

	var est LinearEstimate

	est.Add(1, 1)
	est.Add(2, 3)
	est.Add(3, 5)
	fmt.Printf("%v\n", est)
	fmt.Println("est(4) = ", est.Estimate(4))
	est.MoveX(1)
	fmt.Printf("%v\n", est)
	fmt.Println("est(5) = ", est.Estimate(5))

	est.Add(5, 8)
	fmt.Println("est(5) = ", est.Estimate(5))

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				fmt.Println(ip)
			case *net.IPAddr:
				ip = v.IP
				fmt.Println(ip)
			}
			// process IP address
		}
	}
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	Fail(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	Fail(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		fmt.Println(string(buf[0:n]))

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
