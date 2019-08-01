package utils

import (
	"fmt"
	"math"
)

//CellRange range of a cell
type CellRange struct {
	X0     float64
	X1     float64
	Y0     float64
	Y1     float64
	width  float64
	height float64
	area   float64
}

func (c *CellRange) calCellArea() {
	r := 6371000.0
	phi := (c.Y0 + c.Y1) / 2 * math.Pi / 180
	c.height = r * math.Sin((c.Y1-c.Y0)/180*math.Pi)
	c.width = r * math.Cos(phi) * math.Sin((c.X1-c.X0)/180*math.Pi)
	c.area = c.height * c.width
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
	return &QuadTreeNode{}
}

//BuildQuadTree build a tree from json array
func BuildQuadTree(data *[]BusiInfo, parent *QuadTreeNode, r *CellRange, leafSizeLimit int, code string) *QuadTreeNode {
	sz := len(*data)
	root := NewQuadTreeNode()
	root.Count = len(*data)
	root.parent = parent
	root.code = code
	root.cellRange = *r
	root.cellRange.calCellArea()
	//fmt.Println("size = ", root.Count, ", range = ", root.cellRange, " leaf:", sz < leafSizeLimit, " code = ", root.code)
	if sz < leafSizeLimit {
		root.data = *data
		root.isLeaf = true
		fmt.Println("size = ", root.Count, ", area = (", int(root.cellRange.width), int(root.cellRange.height), ") leaf:", sz < leafSizeLimit, " code = ", root.code)
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
			for i := 0; i < 4; i++ {
				if d.IsInRange(&ranges[i]) {
					infos[i] = append(infos[i], d)
				}
			}
		}
		for i := 0; i < 4; i++ {
			root.Children[i] = BuildQuadTree(&infos[i], root, &ranges[i], leafSizeLimit, code+string('A'+i))
		}
	}
	return root
}

//Search all business in the range
func (t *QuadTreeNode) Search(r *CellRange) []BusiInfo {
	return nil
}
