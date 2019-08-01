package utils

import "fmt"

//BusiInfo alias name

//CellRange range of a cell
type CellRange struct {
	X0 float64
	X1 float64
	Y0 float64
	Y1 float64
}

//QuadTreeNode quad tree node
type QuadTreeNode struct {
	Children  [4]*QuadTreeNode // TR -> BR -> BL -> TL
	parent    *QuadTreeNode
	isLeaf    bool
	data      []BusiInfo
	cellRange CellRange
	Count     int
	code      string
}

//IsInRange check wether a business inside of a range
func (b *BusiInfo) IsInRange(r *CellRange) bool {
	lon := (*b)["longitude"].(float64)
	lat := (*b)["latitude"].(float64)
	//	fmt.Println(lon, lat)
	if lon >= r.X0 && lon < r.X1 && lat >= r.Y0 && lat < r.Y1 {
		return true
	}
	return false
}

//NewQuadTreeNode create a default node
func NewQuadTreeNode() *QuadTreeNode {
	return &QuadTreeNode{code: ""}
}

//BuildQuadTree build a tree from json array
func BuildQuadTree(data *[]BusiInfo, parent *QuadTreeNode, r *CellRange, leafSizeLimit int) *QuadTreeNode {
	var root *QuadTreeNode
	sz := len(*data)
	if sz == 0 || leafSizeLimit == 0 {
		return root
	}
	root = NewQuadTreeNode()
	root.Count = len(*data)
	root.parent = parent
	//root.code = parent.code
	fmt.Println("size = ", root.Count, ", range = ", *r, " leaf:", sz < leafSizeLimit, " code = ", root.code)
	if sz < leafSizeLimit {
		root.data = *data
		root.cellRange = *r
		root.isLeaf = true
	} else {
		root.isLeaf = false
		mx := (r.X0 + r.X1) / 2
		my := (r.Y0 + r.Y1) / 2
		var ranges [4]CellRange
		ranges[0] = CellRange{X0: mx, X1: r.X1, Y0: my, Y1: r.Y1} //TR
		ranges[1] = CellRange{X0: mx, X1: r.X1, Y0: r.Y0, Y1: my} //BR
		ranges[2] = CellRange{X0: r.X0, X1: mx, Y0: r.Y0, Y1: my} //BL
		ranges[3] = CellRange{X0: r.X0, X1: mx, Y0: my, Y1: r.Y1} //TL

		var infos [4][]BusiInfo
		for _, d := range *data {
			//fmt.Println(d)
			for i := 0; i < 4; i++ {
				if d.IsInRange(&ranges[i]) {
					infos[i] = append(infos[i], d)
				}
			}
		}
		for i := 0; i < 4; i++ {
			root.Children[i] = BuildQuadTree(&infos[i], root, &ranges[i], leafSizeLimit)
			//		root.Children[i].code = root.code + string('A'+i)
		}
	}
	//fmt.Println("size = ", root.Count, ", range = ", *r, " leaf:", sz < leafSizeLimit, " code = ", root.code)

	return root
}

// x0 := math.MaxFloat64 / 2
// X1 := -x0
// Y0 := x0
// Y1 := -x0
// for _, b := range *data {
// 	lon := b["longitude"].(float64)
// 	lat := b["latitude"].(float64)
// 	x0 = math.Min(x0, lon)
// 	X1 = math.Max(X1, lon)
// 	Y0 = math.Min(Y0, lat)
// 	Y1 = math.Max(Y1, lat)
// }
