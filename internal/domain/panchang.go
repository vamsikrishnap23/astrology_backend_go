package domain

type PanchangResult struct {
	Date        string        `json:"date"`
	LocalTime   string        `json:"local_time"`
	Timezone    float64       `json:"timezone"`
	Sunrise     string        `json:"sunrise"`
	Sunset      string        `json:"sunset"`
	SolarNoon   string        `json:"solar_noon"`
	Moonrise    string        `json:"moonrise"`
	Moonset     string        `json:"moonset"`
	Vara        Vara          `json:"vara"`
	Tithi       Tithi         `json:"tithi"`
	Nakshatra   Nakshatra     `json:"nakshatra"`
	Yoga        Yoga          `json:"yoga"`
	Karana      Karana        `json:"karana"`
	RahuKalam   DailyPeriod   `json:"rahu_kalam"`
	Yamaganda   DailyPeriod   `json:"yamaganda"`
	Durmuhurtam []DailyPeriod `json:"durmuhurtam"`
	Varjyam     []DailyPeriod `json:"varjyam"`
}

type Vara struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	Ruler  string `json:"ruler"`
}

type Tithi struct {
	Number   int     `json:"number"`
	Name     string  `json:"name"`
	Paksha   string  `json:"paksha"`
	Progress float64 `json:"progress"`
	Start    string  `json:"start"`
	End      string  `json:"end"`
}

type Nakshatra struct {
	Number   int     `json:"number"`
	Name     string  `json:"name"`
	Pada     int     `json:"pada"`
	Progress float64 `json:"progress"`
	Start    string  `json:"start"`
	End      string  `json:"end"`
	Ruler    string  `json:"ruler,omitempty"`
}

type Yoga struct {
	Number   int     `json:"number"`
	Name     string  `json:"name"`
	Progress float64 `json:"progress"`
	Start    string  `json:"start"`
	End      string  `json:"end"`
}

type Karana struct {
	Number   int     `json:"number"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Progress float64 `json:"progress"`
	Start    string  `json:"start"`
	End      string  `json:"end"`
}

type DailyPeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
