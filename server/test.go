package main

import (
	"fmt"

	"yelp/server/utils"
)

//BusiInfo alias name
type BusiInfo = utils.BusiInfo
type QuadTreeNode = utils.QuadTreeNode
type UnitDB = utils.UnitDB
type GeoHashDB = utils.GeoHashDB

//CellRange range of a cell
type Cell = utils.Cell

func main() {
	filename := "/home/xihua/go/src/yelp/data/vegas.json"
	data, err := utils.LoadJSONFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data))
	// cell := *utils.NewCell(-115.6, -114.74, 35.9, 36.6)
	// tree := utils.NewQuadTreeNode(&cell)
	db := *utils.NewGeoHashDB()
	tree := utils.BuildQuadTreeFromData(data, 100, db)
	tree.Info()
	// for k, v := range db {
	// 	fmt.Println(k, len(*v))
	// }
	println(len(db))
	lat := 36.181177
	lon := -115.116080
	dist := 1400.0
	search_terms := "food"
	category := "*"
	db_indices := utils.FindAllDBInRange(tree, lat, lon, dist)
	candidates, distances := utils.FindALlBusiInRange(db, db_indices, lat, lon, dist, search_terms, category)

	for i, b := range candidates {
		fmt.Println(i, int(distances[i]), b["name"])
	}

	// crowd := *db["RBDDDDBAAAA"]
	// for _, p := range crowd {
	// 	fmt.Println(p)
	// }
	// var ar1 []string
	// var ar2 = []string{"i love you"}
	// ar3 := append(ar1, ar2...)
	// fmt.Println(ar1 == nil, ar2, ar3)

	//c := utils.NewCell(-115.7, -114, 35, 37)
	//strs := tree.FindOverlapLeafs(c)
	//fmt.Println(strs)
	// k := 30000
	// for i := 0; i < k; i++ {
	// 	tree.Insert(&data[i], 1000)
	// 	//	tree.Info()
	// }

	//r := CellRange{X0: -115.6, X1: -114.74, Y0: 35.9, Y1: 36.6}

	// for _, d := range data {
	// 	lon := d["longitude"].(float64)
	// 	lat := d["latitude"].(float64)
	// 	r.X0 = math.Min(r.X0, lon)
	// 	r.X1 = math.Max(r.X1, lon)
	// 	r.Y0 = math.Min(r.Y0, lat)
	// 	r.Y1 = math.Max(r.Y1, lat)
	// }

	// r.CalCellArea()
	// fmt.Println(r)
	// tree := utils.BuildQuadTree(&data, nil, &r, 1000, "T")
	// fmt.Println(tree.Count)
	// hash := make(map[string]*QuadTreeNode)
	// tree.BuildGeoHash(&hash)
	// for k, v := range hash {
	// 	fmt.Println(k, v.Count)
	// }
}
