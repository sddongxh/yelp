package utils

import (
	"math"
	"sort"
	"strings"
)

//UnitDB at one server
type UnitDB []BusiInfo

//NewUnitDB generate a new unit db
func NewUnitDB() *UnitDB {
	db := make(UnitDB, 0)
	return &db
}

//GeoHashDB distributed database, GeoHash
type GeoHashDB map[string]*UnitDB

//NewGeoHashDB generate a new db
func NewGeoHashDB() *GeoHashDB {
	db := make(GeoHashDB)
	return &db
}

func distance(lat0, lon0, lat1, lon1 float64) float64 {
	r := 6371000.0
	phi := (lat0 + lat1) / 2 * math.Pi / 180
	alpha := (lon1 - lon0) / 180 * math.Pi
	beta := (lat1 - lat0) / 180 * math.Pi
	dx := r * math.Cos(phi) * alpha
	dy := r * beta
	return math.Sqrt(dx*dx + dy*dy)
}

//ByRating implements sort by rating
type ByRating []BusiInfo

func (a ByRating) Len() int           { return len(a) }
func (a ByRating) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRating) Less(i, j int) bool { return a[i]["stars"].(float64) > a[j]["stars"].(float64) }

//SearchBusiInRange find related businesses in unit DB
func (d *UnitDB) SearchBusiInRange(lat, lon, dist float64, searchTerms, category string) ([]BusiInfo, []float64) {
	res := make([]BusiInfo, 0)
	dists := make([]float64, 0)
	for _, b := range *d {
		lon1 := b["longitude"].(float64)
		lat1 := b["latitude"].(float64)
		v := distance(lat, lon, lat1, lon1)
		if v < dist {
			if b["categories"] == nil || b["name"] == nil || b["address"] == nil {
				continue
			}
			str := b["name"].(string) + " " + b["address"].(string) + " " + b["categories"].(string)
			str = strings.ToLower(str)
			//fmt.Println(str, strings.ToLower(searchTerms))
			if strings.Contains(str, strings.ToLower(searchTerms)) {
				res = append(res, b)
				dists = append(dists, v)
			}
		}
	}
	//sort.Sort(ByRating(res))
	// for i, b := range res {
	// 	fmt.Println(i, b["stars"])
	// }
	return res, dists
}

//FindALlBusiInRange implements merge sort
func FindALlBusiInRange(db GeoHashDB, indices []string, lat, lon, dist float64, searchTerms, category string) ([]BusiInfo, []float64) {
	result := make([]BusiInfo, 0)
	distances := make([]float64, 0)
	for i := 0; i < len(indices); i++ {
		ind := indices[i]
		bs, dists := db[ind].SearchBusiInRange(lat, lon, dist, searchTerms, category)
		result = append(result, bs...)
		distances = append(distances, dists...)
	}
	sort.Sort(ByRating(result))
	return result, distances
}
