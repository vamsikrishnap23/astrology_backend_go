package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	content, _ := ioutil.ReadFile("internal/astrology/btr/btr.go")
	str := string(content)

	// Fix the first error: not enough arguments to getNadiPlanet in the scan range (line 336)
	str = strings.Replace(str, `match, checkTatwa, _, _ := evaluateRow(wrappedRow, weekday, ascType, input.Gender, actualStarLord)`,
		`match, checkTatwa, _, _ := evaluateRow(wrappedRow, weekday, ascType, input.Gender, actualStarLord)`, -1) // wait, evaluateRow calls it correctly now. Where is getNadiPlanet called with 2 args?

	// Wait, let's just grep where getNadiPlanet is called.
	fmt.Println("Done")
	ioutil.WriteFile("internal/astrology/btr/btr.go", []byte(str), 0644)
}
