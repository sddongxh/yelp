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

//Info print the information
func (c *Cell) Info() {
	fmt.Println("width:", c.width, "height:", c.height, "area:", c.area)
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

//IsOverlap check whether overlap
func (c *Cell) IsOverlap(d *Cell) bool {
	if c.x0 >= d.x1 || d.x0 >= c.x1 {
		return false
	}
	if c.y0 >= d.y1 || d.y0 >= c.y1 {
		return false
	}
	return true
}

//QuadTreeNode quad tree node
type QuadTreeNode struct {
	children [4]*QuadTreeNode // TR -> BR -> BL -> TL
	parent   *QuadTreeNode
	isLeaf   bool
	data     UnitDB
	cell     Cell
	count    int
	code     string
}

func (t QuadTreeNode) String() string {
	return fmt.Sprintf("W = %v, H = %v, total # of business = %v", int(t.cell.width), int(t.cell.height), t.count)
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
func (t *QuadTreeNode) Insert(b *BusiInfo, limit int, db GeoHashDB) bool {
	if !t.cell.IsContain(b) {
		return false
	}
	if !t.isLeaf {
		for _, child := range t.children {
			if child.Insert(b, limit, db) {
				t.count++
				return true
			}
		}
		return false
	} else if t.count+1 <= limit || t.cell.width < 100 {
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
			db[t.children[i].code] = &t.children[i].data
		}
		//assign existing data
		for _, d := range t.data {
			for i := 0; i < 4; i++ {
				if t.children[i].Insert(&d, limit, db) {
					break
				}
			}
		}
		t.data = nil
		delete(db, t.code)
		return t.Insert(b, limit, db)
	}
}

//IsOverlap check whether overlap with a region
func (t *QuadTreeNode) IsOverlap(c *Cell) bool {
	return t.cell.IsOverlap(c)
}

//FindOverlapLeafs find all overlapped leafs
func (t *QuadTreeNode) FindOverlapLeafs(c *Cell) []string {
	if !t.IsOverlap(c) {
		return nil
	} else if t.isLeaf {
		return []string{t.code}
	} else {
		var codes []string
		for i := 0; i < 4; i++ {
			codes = append(codes, t.children[i].FindOverlapLeafs(c)...)
		}
		return codes
	}
}

//BuildQuadTreeFromData build a quad tree from data
func BuildQuadTreeFromData(data []BusiInfo, limit int, db GeoHashDB) *QuadTreeNode {
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
	cx := (x0 + x1) / 2
	cy := (y0 + y1) / 2
	w := math.Max(x1-x0, y1-y0) * 1.1
	scale := math.Cos((y0 + y1) / 2 * math.Pi / 180)
	cell := NewCell(cx-w/2, cx+w/2, cy-w/2*scale, cy+w/2*scale) //Make as square as possible
	tree := NewQuadTreeNode(cell)
	db[tree.code] = &tree.data
	for i := 0; i < len(data); i++ {
		if !tree.Insert(&data[i], limit, db) {
			fmt.Println(data[i])
		}
	}
	return tree
}

//FindAllDBInRange implements indexing
func FindAllDBInRange(t *QuadTreeNode, lat, lon, dist float64) []string {
	r := 6371000.0
	a := dist / (r * math.Cos(lat/180*math.Pi)) / math.Pi * 180
	scale := math.Cos(lat / 180 * math.Pi)
	c := NewCell(lon-a, lon+a, lat-a*scale, lat+a*scale)
	return t.FindOverlapLeafs(c)
}
