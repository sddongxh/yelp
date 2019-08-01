package main

import (
	"fmt"

	"yelp/server/utils"
)

//BusiInfo alias name
type BusiInfo = utils.BusiInfo
type QuadTreeNode = utils.QuadTreeNode

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
	tree := utils.BuildQuadTreeFromData(data, 1000)
	tree.Info()

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
