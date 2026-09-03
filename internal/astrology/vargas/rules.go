package vargas

import (
	"math"
)

// VargaPosition holds the calculated divisional sign index (0-11) and the exact longitude inside that divisional sign (0-30).
type VargaPosition struct {
	SignIndex           int
	LongitudeInDivision float64
}

// VargaRule defines the interface for calculating a specific divisional chart.
type VargaRule interface {
	Name() string
	Division() int
	Calculate(longitude float64) VargaPosition
}

var Registry = []VargaRule{
	&RuleD1{},
	&RuleD2{},
	&RuleD3{},
	&RuleD4{},
	&RuleD7{},
	&RuleD9{},
	&RuleD10{},
	&RuleD12{},
	&RuleD16{},
	&RuleD20{},
	&RuleD24{},
	&RuleD27{},
	&RuleD30{},
	&RuleD40{},
	&RuleD45{},
	&RuleD60{},
}

// Common helper to wrap signs between 0 and 11
func wrapSign(s int) int {
	res := s % 12
	if res < 0 {
		res += 12
	}
	return res
}

// --- D1 Rasi ---
type RuleD1 struct{}

func (r *RuleD1) Name() string  { return "D1" }
func (r *RuleD1) Division() int { return 1 }
func (r *RuleD1) Calculate(lon float64) VargaPosition {
	sign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	return VargaPosition{SignIndex: wrapSign(sign), LongitudeInDivision: rem}
}

// --- D2 Hora ---
type RuleD2 struct{}

func (r *RuleD2) Name() string  { return "D2" }
func (r *RuleD2) Division() int { return 2 }
func (r *RuleD2) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	isOddSign := baseSign%2 == 0 // 0=Aries (Odd)
	part := int(math.Floor(rem / 15.0))

	var divSign int
	if isOddSign {
		if part == 0 {
			divSign = 4 // Leo
		} else {
			divSign = 3 // Cancer
		}
	} else {
		if part == 0 {
			divSign = 3 // Cancer
		} else {
			divSign = 4 // Leo
		}
	}
	return VargaPosition{SignIndex: divSign, LongitudeInDivision: math.Mod(rem, 15.0) * 2.0}
}

// --- D3 Drekkana ---
type RuleD3 struct{}

func (r *RuleD3) Name() string  { return "D3" }
func (r *RuleD3) Division() int { return 3 }
func (r *RuleD3) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	part := int(math.Floor(rem / 10.0))

	var divSign int
	if part == 0 {
		divSign = baseSign
	} else if part == 1 {
		divSign = baseSign + 4
	} else {
		divSign = baseSign + 8
	}
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, 10.0) * 3.0}
}

// --- D4 Chaturthamsha ---
type RuleD4 struct{}

func (r *RuleD4) Name() string  { return "D4" }
func (r *RuleD4) Division() int { return 4 }
func (r *RuleD4) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	part := int(math.Floor(rem / 7.5))

	divSign := baseSign + (part * 3)
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, 7.5) * 4.0}
}

// --- D7 Saptamsha ---
type RuleD7 struct{}

func (r *RuleD7) Name() string  { return "D7" }
func (r *RuleD7) Division() int { return 7 }
func (r *RuleD7) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	isOddSign := baseSign%2 == 0
	partLen := 30.0 / 7.0
	part := int(math.Floor(rem / partLen))

	var divSign int
	if isOddSign {
		divSign = baseSign + part
	} else {
		divSign = baseSign + 6 + part
	}
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 7.0}
}

// --- D9 Navamsha ---
type RuleD9 struct{}

func (r *RuleD9) Name() string  { return "D9" }
func (r *RuleD9) Division() int { return 9 }
func (r *RuleD9) Calculate(lon float64) VargaPosition {
	partLen := 30.0 / 9.0
	totalParts := int(math.Floor(lon / partLen))
	rem := math.Mod(lon, partLen)

	divSign := totalParts % 12
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: rem * 9.0}
}

// --- D10 Dashamsha ---
type RuleD10 struct{}

func (r *RuleD10) Name() string  { return "D10" }
func (r *RuleD10) Division() int { return 10 }
func (r *RuleD10) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	isOddSign := baseSign%2 == 0
	part := int(math.Floor(rem / 3.0))

	var divSign int
	if isOddSign {
		divSign = baseSign + part
	} else {
		divSign = baseSign + 8 + part
	}
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, 3.0) * 10.0}
}

// --- D12 Dwadashamsha ---
type RuleD12 struct{}

func (r *RuleD12) Name() string  { return "D12" }
func (r *RuleD12) Division() int { return 12 }
func (r *RuleD12) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	part := int(math.Floor(rem / 2.5))

	divSign := baseSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, 2.5) * 12.0}
}

// --- D16 Shodashamsha ---
type RuleD16 struct{}

func (r *RuleD16) Name() string  { return "D16" }
func (r *RuleD16) Division() int { return 16 }
func (r *RuleD16) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	partLen := 30.0 / 16.0
	part := int(math.Floor(rem / partLen))

	modality := baseSign % 3
	var startSign int
	if modality == 0 { // Movable
		startSign = 0 // Aries
	} else if modality == 1 { // Fixed
		startSign = 4 // Leo
	} else { // Dual
		startSign = 8 // Sag
	}

	divSign := startSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 16.0}
}

// --- D20 Vimshamsha ---
type RuleD20 struct{}

func (r *RuleD20) Name() string  { return "D20" }
func (r *RuleD20) Division() int { return 20 }
func (r *RuleD20) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	partLen := 30.0 / 20.0
	part := int(math.Floor(rem / partLen))

	modality := baseSign % 3
	var startSign int
	if modality == 0 { // Movable
		startSign = 0 // Aries
	} else if modality == 1 { // Fixed
		startSign = 8 // Sag
	} else { // Dual
		startSign = 4 // Leo
	}

	divSign := startSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 20.0}
}

// --- D24 Chaturvimshamsha ---
type RuleD24 struct{}

func (r *RuleD24) Name() string  { return "D24" }
func (r *RuleD24) Division() int { return 24 }
func (r *RuleD24) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	isOddSign := baseSign%2 == 0
	partLen := 30.0 / 24.0
	part := int(math.Floor(rem / partLen))

	var startSign int
	if isOddSign {
		startSign = 4 // Leo
	} else {
		startSign = 3 // Cancer
	}

	divSign := startSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 24.0}
}

// --- D27 Saptavimshamsha ---
type RuleD27 struct{}

func (r *RuleD27) Name() string  { return "D27" }
func (r *RuleD27) Division() int { return 27 }
func (r *RuleD27) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	partLen := 30.0 / 27.0
	part := int(math.Floor(rem / partLen))

	element := baseSign % 4
	var startSign int
	if element == 0 { // Fire
		startSign = 0 // Aries
	} else if element == 1 { // Earth
		startSign = 3 // Cancer
	} else if element == 2 { // Air
		startSign = 6 // Libra
	} else { // Water
		startSign = 9 // Capricorn
	}

	divSign := startSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 27.0}
}

// --- D30 Trimshamsha ---
type RuleD30 struct{}

func (r *RuleD30) Name() string  { return "D30" }
func (r *RuleD30) Division() int { return 30 }
func (r *RuleD30) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	isOddSign := baseSign%2 == 0

	var divSign int
	// Unequal degrees
	if isOddSign {
		if rem < 5.0 {
			divSign = 0 // Aries (Mars)
		} else if rem < 10.0 {
			divSign = 10 // Aquarius (Saturn)
		} else if rem < 18.0 {
			divSign = 8 // Sagittarius (Jupiter)
		} else if rem < 25.0 {
			divSign = 2 // Gemini (Mercury)
		} else {
			divSign = 6 // Libra (Venus)
		}
	} else {
		if rem < 5.0 {
			divSign = 1 // Taurus (Venus)
		} else if rem < 12.0 {
			divSign = 5 // Virgo (Mercury)
		} else if rem < 20.0 {
			divSign = 11 // Pisces (Jupiter)
		} else if rem < 25.0 {
			divSign = 9 // Capricorn (Saturn)
		} else {
			divSign = 7 // Scorpio (Mars)
		}
	}
	// For D30, the "LongitudeInDivision" isn't perfectly mapped linearly since lengths are unequal,
	// but standard practice is often just passing 0 or scaling linearly. We scale it.
	// For example, if it's in 5 degree bound, scale to 30.
	var partStart, partLen float64
	if isOddSign {
		if rem < 5.0 {
			partStart = 0.0
			partLen = 5.0
		} else if rem < 10.0 {
			partStart = 5.0
			partLen = 5.0
		} else if rem < 18.0 {
			partStart = 10.0
			partLen = 8.0
		} else if rem < 25.0 {
			partStart = 18.0
			partLen = 7.0
		} else {
			partStart = 25.0
			partLen = 5.0
		}
	} else {
		if rem < 5.0 {
			partStart = 0.0
			partLen = 5.0
		} else if rem < 12.0 {
			partStart = 5.0
			partLen = 7.0
		} else if rem < 20.0 {
			partStart = 12.0
			partLen = 8.0
		} else if rem < 25.0 {
			partStart = 20.0
			partLen = 5.0
		} else {
			partStart = 25.0
			partLen = 5.0
		}
	}
	valInDivision := ((rem - partStart) / partLen) * 30.0
	return VargaPosition{SignIndex: divSign, LongitudeInDivision: valInDivision}
}

// --- D40 Khavedamsha ---
type RuleD40 struct{}

func (r *RuleD40) Name() string  { return "D40" }
func (r *RuleD40) Division() int { return 40 }
func (r *RuleD40) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	isOddSign := baseSign%2 == 0
	partLen := 30.0 / 40.0
	part := int(math.Floor(rem / partLen))

	var startSign int
	if isOddSign {
		startSign = 0 // Aries
	} else {
		startSign = 6 // Libra
	}

	divSign := startSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 40.0}
}

// --- D45 Akshavedamsha ---
type RuleD45 struct{}

func (r *RuleD45) Name() string  { return "D45" }
func (r *RuleD45) Division() int { return 45 }
func (r *RuleD45) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	partLen := 30.0 / 45.0
	part := int(math.Floor(rem / partLen))

	modality := baseSign % 3
	var startSign int
	if modality == 0 { // Movable
		startSign = 0 // Aries
	} else if modality == 1 { // Fixed
		startSign = 4 // Leo
	} else { // Dual
		startSign = 8 // Sag
	}

	divSign := startSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 45.0}
}

// --- D60 Shashtiamsha ---
type RuleD60 struct{}

func (r *RuleD60) Name() string  { return "D60" }
func (r *RuleD60) Division() int { return 60 }
func (r *RuleD60) Calculate(lon float64) VargaPosition {
	baseSign := int(math.Floor(lon / 30.0))
	rem := math.Mod(lon, 30.0)
	partLen := 30.0 / 60.0
	part := int(math.Floor(rem / partLen))

	divSign := baseSign + part
	return VargaPosition{SignIndex: wrapSign(divSign), LongitudeInDivision: math.Mod(rem, partLen) * 60.0}
}
