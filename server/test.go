package main

import (
	"fmt"
	"math"

	"yelp/server/utils"
)

//BusiInfo alias name
type BusiInfo = utils.BusiInfo

//CellRange range of a cell
type CellRange = utils.CellRange

func main() {
	filename := "/home/xihua/go/src/yelp/data/vegas.json"
	data, err := utils.LoadJSONFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data), err)

	r := CellRange{X0: 180, X1: -180, Y0: 90, Y1: -90}

	for _, d := range data {
		lon := d["longitude"].(float64)
		lat := d["latitude"].(float64)
		r.X0 = math.Min(r.X0, lon)
		r.X1 = math.Max(r.X1, lon)
		r.Y0 = math.Min(r.Y0, lat)
		r.Y1 = math.Max(r.Y1, lat)
	}

	fmt.Println(r)
	tree := utils.BuildQuadTree(&data, nil, &r, 10000)
	fmt.Println(tree.Count)
}
