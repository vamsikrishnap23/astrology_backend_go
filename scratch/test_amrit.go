package main

import "fmt"

func main() {
	varjyamGhatis := []float64{50, 24, 30, 40, 14, 21, 30, 20, 32, 30, 20, 18, 21, 20, 14, 14, 10, 14, 56, 24, 20, 10, 10, 18, 16, 24, 30}
	
	fmt.Print("amruthaGhatis := []float64{")
	for i, v := range varjyamGhatis {
		a := v + 24
		if a >= 60 {
			a -= 60
		}
		fmt.Printf("%.0f", a)
		if i < len(varjyamGhatis)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")
}
