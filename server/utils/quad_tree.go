package utils

import (
	"fmt"
	"math"
)

//Cell range of a cell
type Cell struct {
	x0     float64
	x1     float64
	y0     float64
	y1     float64
	width  float64
	height float64
	area   float64
}

//NewCell build a new cell
func NewCell(x0, x1, y0, y1 float64) *Cell {
	c := Cell{x0: x0, x1: x1, y0: y0, y1: y1}
	r := 6371000.0
	phi := (c.y0 + c.y1) / 2 * math.Pi / 180
	c.height = r * math.Sin((c.y1-c.y0)/180*math.Pi)
	c.width = r * math.Cos(phi) * math.Sin((c.x1-c.x0)/180*math.Pi)
	c.area = c.height * c.width
	return &c
}

//IsContain check wether a business inside of a cell
func (c *Cell) IsContain(b *BusiInfo) bool {
	lon := (*b)["longitude"].(float64)
	lat := (*b)["latitude"].(float64)
	if lon >= c.x0 && lon < c.x1 && lat >= c.y0 && lat < c.y1 {
		return true
	}
	return false
}

//QuadTreeNode quad tree node
type QuadTreeNode struct {
	children [4]*QuadTreeNode // TR -> BR -> BL -> TL
	parent   *QuadTreeNode
	isLeaf   bool
	data     []BusiInfo
	cell     Cell
	count    int
	code     string
}

//Info print the information
func (t *QuadTreeNode) Info() {
	fmt.Println("count:", t.count, "isLeaf:", t.isLeaf, "code:", t.code, "Area: ", t.cell.width, t.cell.height)
}

//NewQuadTreeNode create a default node
func NewQuadTreeNode(c *Cell) *QuadTreeNode {
	return &QuadTreeNode{cell: *c, isLeaf: true, count: 0, parent: nil, code: "R", data: nil}
}

//Insert insert a data point
func (t *QuadTreeNode) Insert(b *BusiInfo, limit int) bool {
	if !t.cell.IsContain(b) {
		return false
	}
	if !t.isLeaf {
		for _, child := range t.children {
			if child.Insert(b, limit) {
				t.count++
				return true
			}
		}
		return false
	} else if t.count+1 <= limit {
		d := *b
		t.data = append(t.data, d)
		//fmt.Println("inerted to ", t.code)
		t.count++
		return true
	} else {
		// partition data
		t.isLeaf = false
		c := t.cell
		mx := (c.x0 + c.x1) / 2
		my := (c.y0 + c.y1) / 2
		var cells [4]*Cell
		cells[0] = NewCell(mx, c.x1, my, c.y1) //TR
		cells[1] = NewCell(mx, c.x1, c.y0, my) //BR
		cells[2] = NewCell(c.x0, mx, c.y0, my) //BL
		cells[3] = NewCell(c.x0, mx, my, c.y1) //TL
		for i := 0; i < 4; i++ {
			t.children[i] = NewQuadTreeNode(cells[i])
			t.children[i].code = t.code + string('A'+i)
		}
		//assign existing data
		for _, d := range t.data {
			for i := 0; i < 4; i++ {
				if t.children[i].Insert(&d, limit) {
					break
				}
			}
		}
		t.data = nil
		return t.Insert(b, limit)
	}
}

//BuildQuadTreeFromData build a quad tree from data
func BuildQuadTreeFromData(data []BusiInfo, limit int) *QuadTreeNode {
	//determine the range
	x0 := 180.0
	x1 := -180.0
	y0 := 90.0
	y1 := -90.0
	for _, d := range data {
		lon := d["longitude"].(float64)
		lat := d["latitude"].(float64)
		x0 = math.Min(x0, lon)
		x1 = math.Max(x1, lon)
		y0 = math.Min(y0, lat)
		y1 = math.Max(y1, lat)
	}
	fmt.Println(x0, x1, y0, y1)
	cx := (x0 + x1) / 2
	cy := (y0 + y1) / 2
	w := math.Max(x1-x0, y1-y0) * 1.1
	scale := math.Cos((y0 + y1) / 2 * math.Pi / 180)
	cell := NewCell(cx-w/2, cx+w/2, cy-w/2*scale, cy+w/2*scale) //Make as square as possible
	tree := NewQuadTreeNode(cell)
	for i := 0; i < len(data); i++ {
		if !tree.Insert(&data[i], limit) {
			fmt.Println(data[i])
		}
	}
	return tree
}

//BuildGeoHash Build the GeoHash table
func (t *QuadTreeNode) BuildGeoHash(hash *map[string]*QuadTreeNode) {
	if t.isLeaf {
		(*hash)[t.code] = t
	} else {
		for _, c := range t.children {
			c.BuildGeoHash(hash)
		}
	}
}

// 	else if t.isLeaf && t.count < limit {
// 		d := *b
// 		t.data = append(t.data, d)
// 		return true
// 	} else {
// 		// parition data

// 	}

// 	if !t.isLeaf {
// 		for _, child := range t.children {
// 			if child.cell.IsContain(b) {
// 				return child.Insert(b, limit)
// 			}
// 			return false
// 		}
// 	} else {
// 		if t.count <= limit {
// 			d := *b
// 			t.data = append(t.data, d)
// 		} else { // partition four cells
// 			t.isLeaf = false
// 			c := t.cell
// 			mx := (c.x0 + c.x1) / 2
// 			my := (c.y0 + c.y1) / 2
// 			var cells [4] *Cell
// 			cells[0] = NewCell(mx, c.x1, my, c.y1) //TR
// 			cells[1] = NewCell(mx, c.x1, c.y0, my) //BR
// 			cells[2] = NewCell(c.x0, mx, c.y0, my) //BL
// 			cells[3] = NewCell(c.x0, mx, my, c.y1) //TL
// 			for i := 0; i < 4; i++ {
// 				t.children[i] = NewQuadTreeNode(cells[i])
// 			}
// 			for _, child := range t.children {
// 				if child.cell.IsContain(b) {
// 					child.Insert(b, limit)
// 				}
// 			}
// 		}
// 	}
// 	return true
// }

// // //BuildQuadTree build a tree from json array
// // func BuildQuadTree(data *[]BusiInfo, parent *QuadTreeNode, c *Cell, code string, leafSizeLimit int) *QuadTreeNode {
// // 	root := NewQuadTreeNode()
// // 	root.count = len(*data)
// // 	root.parent = parent
// // 	root.code = code
// // 	root.cell = *c
// // 	if root.count < leafSizeLimit {
// // 		root.data = *data
// // 		root.isLeaf = true
// // 	} else {
// // 		root.isLeaf = false
// // 		mx := (c.x0 + c.x1) / 2
// // 		my := (c.y0 + c.y1) / 2
// // 		var cells [4]Cell
// // 		cells[0] = *NewCell(mx, c.x1, my, c.y1) //TR
// // 		cells[1] = *NewCell(mx, c.x1, c.y0, my) //BR
// // 		cells[2] = *NewCell(c.x0, mx, c.y0, my) //BL
// // 		cells[3] = *NewCell(c.x0, mx, my, c.y1) //TL
// // 		//partition the data
// // 		for _, d := range *data {
// // 			for i := 0; i < 4; i++ {
// // 				if cells[i].IsContain(d) {
// // 					root.children[i].data
// // 				}
// // 			}
// // 		}

// // 	}

// 	// 	root.Cell = *r
// 	// 	root.Cell.CalCellArea()
// 	// 	//fmt.Println("size = ", root.Count, ", range = ", root.Cell, " leaf:", sz < leafSizeLimit, " code = ", root.code)
// 	// 	if sz < leafSizeLimit {
// 	// 		root.data = *data
// 	// 		root.isLeaf = true
// 	// 		fmt.Println("size = ", root.Count, ", area = (", int(root.Cell.width), int(root.Cell.height), ") leaf:", sz < leafSizeLimit, " code = ", root.code)
// 	// 	} else {
// 	// 		root.isLeaf = false
// 	// 		mx := (r.X0 + r.X1) / 2
// 	// 		my := (r.Y0 + r.Y1) / 2
// 	// 		var ranges [4]Cell
// 	// 		ranges[0] = Cell{X0: mx, X1: r.X1, Y0: my, Y1: r.Y1} //TR
// 	// 		ranges[1] = Cell{X0: mx, X1: r.X1, Y0: r.Y0, Y1: my} //BR
// 	// 		ranges[2] = Cell{X0: r.X0, X1: mx, Y0: r.Y0, Y1: my} //BL
// 	// 		ranges[3] = Cell{X0: r.X0, X1: mx, Y0: my, Y1: r.Y1} //TL

// 	// 		var infos [4][]BusiInfo
// 	// 		for _, d := range *data {
// 	// 			for i := 0; i < 4; i++ {
// 	// 				if d.IsInRange(&ranges[i]) {
// 	// 					infos[i] = append(infos[i], d)
// 	// 				}
// 			}
// 		}
// 		for i := 0; i < 4; i++ {
// 			root.Children[i] = BuildQuadTree(&infos[i], root, &ranges[i], leafSizeLimit, code+string('A'+i))
// 		}
// 	}
// 	return root
// }

// //BuildGeoHash Build the GeoHash table
// func (t *QuadTreeNode) BuildGeoHash(hash *map[string]*QuadTreeNode) {
// 	if t.isLeaf {
// 		(*hash)[t.code] = t
// 	} else {
// 		for _, c := range t.Children {
// 			c.BuildGeoHash(hash)
// 		}
// 	}
// }

// //Search all business in the range
// func (t *QuadTreeNode) Search(r *Cell) []BusiInfo {
// 	return nil
// }
