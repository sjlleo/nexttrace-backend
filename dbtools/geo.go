package dbtools

type Geo struct {
	Country string
	Prov    string
	City    string
	Lat     float64
	Lon     float64
}

func AddGeo(g *Geo) {
	db := GetDB()
	db.Model(g).Create(g)
}
