package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arafatk/glot"
)

func main() {
	dimensions := 2
	// The dimensions supported by the plot
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)
	pointGroupName := "Simple Circles"
	style := "points"
	// points := [][]float64{{7, 3, 13, 5.6, 11.1}, {12, 13, 11, 1, 7}}
	// // Adding a point group
	// plot.AddPointGroup(pointGroupName, style, points)
	// // A plot type used to make points/ curves and customize and save them as an image.
	// plot.SetTitle("Example Plot")
	// // Optional: Setting the title of the plot
	// plot.SetXLabel("X-Axis")
	// plot.SetYLabel("Y-Axis")
	// // Optional: Setting label for X and Y axis
	// plot.SetXrange(-2, 18)
	// plot.SetYrange(-2, 18)
	// // Optional: Setting axis ranges
	// plot.SavePlot("2.png")
	// //fmt.Scanln() // wait for Enter Key

	jsonFile, err := os.Open("/home/xihua/code/datasets/yelp/dataset/business.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	reader := bufio.NewReader(jsonFile)
	var business map[string]interface{}
	longitudes := make([]float64, 0)
	latitudes := make([]float64, 0)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		//	fmt.Println(line)
		//zip := line[]
		json.Unmarshal([]byte(line), &business)
		//	zip := business["postal_code"]
		//	fmt.Println(zip)
		longitudes = append(longitudes, business["longitude"].(float64))
		latitudes = append(latitudes, business["latitude"].(float64))
	}
	fmt.Println(longitudes, latitudes)

	ps := [][]float64{longitudes, latitudes}
	plot.AddPointGroup(pointGroupName, style, ps)
	// A plot type used to make points/ curves and customize and save them as an image.
	plot.SetTitle("Example Plot")
	// Optional: Setting the title of the plot
	plot.SetXLabel("X-Axis")
	plot.SetYLabel("Y-Axis")
	// Optional: Setting label for X and Y axis
	plot.SetXrange(-180, 180)
	plot.SetYrange(-90, 90)
	// Optional: Setting axis ranges
	plot.SavePlot("3.png")
	var input string
	fmt.Scanln(&input)

}
