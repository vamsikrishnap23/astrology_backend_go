package shadbala

// tribhagaBala is Tribhaga Bala (Raman art. 56-57; BPHS 27.12): the day (sunrise to sunset) and the
// night (sunset to the next sunrise) are each divided into three equal parts. In the three parts of the
// day Mercury, the Sun and Saturn get 60, in the three parts of the night the Moon, Venus and Mars; Jupiter
// always gets 60.
func tribhagaBala(planet string, day hinduDay, jd float64) float64 {
	if planet == "Jupiter" {
		return 60.0
	}
	var start, length float64
	var lords [3]string
	if jd < day.Sunset {
		start, length = day.Sunrise, day.Sunset-day.Sunrise
		lords = [3]string{"Mercury", "Sun", "Saturn"}
	} else {
		start, length = day.Sunset, day.NextSunrise-day.Sunset
		lords = [3]string{"Moon", "Venus", "Mars"}
	}
	part := int((jd - start) / (length / 3.0))
	if part > 2 {
		part = 2
	}
	if part < 0 {
		part = 0
	}
	if lords[part] == planet {
		return 60.0
	}
	return 0.0
}
