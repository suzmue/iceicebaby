package main

import (
	"fmt"
	"strings"
)

var printEachLattice bool = true // To print lattice, change false to true

func main() {
	// To change the lattice for which you are finding the partition function
	// 1. Change the first argument to be the number of rows
	// 2. Change the second argument to be the number of columns
	// 3. Change the third argument to be an array where the ith element is true if
	//    a path originates in xi and false otherwise
	// 4. Change the fourth argument to be an array where the jth element is true if
	//    a path exits at column j and false otherwise.
	weights := findPartitionFunction(3, 3, []bool{false, false, true}, []bool{false, false, true})
	if len(weights) == 0 {
		fmt.Println("There were no lattices that satisfy the entered constraints.")
		return
	}
	weights = simplify(weights)
	fmt.Printf("partition function = %s", polynomialToString(weights[0]))
	for i := 1; i < len(weights); i++ {
		fmt.Printf(" + %s", polynomialToString(weights[i]))
	}
	fmt.Println()
}

const (
	in  = false
	out = true
)

type vertex struct {
	// Direction of edges
	North bool
	East  bool
	South bool
	West  bool
}

// validate checks whether lattice is a valid lattice with the correct
// number of rows and columns
func validate(lattice [][]vertex, rows, columns int, inputs, outputs []bool) bool {
	if len(lattice) != rows {
		return false
	}
	for _, row := range lattice {
		if len(row) != columns {
			return false
		}
	}

	if len(inputs) != rows || len(outputs) != columns {
		return false
	}

	numIn := 0
	for i := 0; i < len(inputs); i++ {
		if inputs[i] {
			numIn++
		}
	}

	numOut := 0
	for i := 0; i < len(outputs); i++ {
		if outputs[i] {
			numOut++
		}
	}

	if numIn != numOut {
		return false
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			// TODO:
			// - Those that dont have paths go down and right
			// - there have to be a path from a to b
			// -
			if i > 0 {
				if lattice[i-1][j].North == lattice[i][j].South {
					return false
				}
			}

			if j > 0 {
				if lattice[i][j-1].West == lattice[i][j].East {
					return false
				}
			}
		}
	}

	// Check that each have two in and two out
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			var w = func(x bool) int {
				if x {
					return 1
				}
				return 0
			}
			v := lattice[i][j]
			count := w(v.North) + w(v.East) + w(v.South) + w(v.West)
			if count != 2 {
				return false
			}

			// Check that none go
			//		^
			//   <		<
			// 		^
			if v.North && v.West && !v.East && !v.South {
				return false
			}
		}
	}

	// Check inputs go correct direction
	for i := 0; i < rows; i++ {
		v := lattice[i][0]
		if inputs[i] == v.East {
			return false
		}
	}

	// Check output go correct direction
	for j := 0; j < columns; j++ {
		v := lattice[rows-1][j]
		if outputs[j] != v.North {
			return false
		}
	}
	// Check left go correct direction
	for i := 0; i < rows; i++ {
		v := lattice[i][columns-1]
		if v.West {
			return false
		}
	}

	// check bottom direction
	for j := 0; j < columns; j++ {
		v := lattice[0][j]
		if !v.South {
			return false
		}
	}

	// Check that there are paths
	checkPaths(lattice, rows, columns, inputs, outputs)

	return true
}

func checkPaths(lattice [][]vertex, rows, columns int, inputs, outputs []bool) bool {
	type st struct {
		s int
		t int
	}
	var sts []st = make([]st, 0, rows+columns)
	i := 0
	j := 0
	for i < len(inputs) && j < len(outputs) {
		if !inputs[i] {
			i++
			continue
		}
		if !outputs[j] {
			j++
			continue
		}

		if inputs[i] && outputs[j] {
			sts = append(sts, st{s: i, t: j})
			i++
			j++
		}
	}

	// Check that the up and left path from s leads to t
	for _, nodes := range sts {
		currI := nodes.s
		currJ := 0
		for currI < rows {
			if lattice[currI][currJ].North {
				currI++
			} else if lattice[currI][currJ].West {
				currJ++
			} else {
				// Bad state
				return false
			}
		}
		if currJ != nodes.t {
			return false
		}
	}
	return false
}

// calculateLatticeWeight calculates the total weight of the lattice, where the weight
// of each vertex is determined by weightFn.
func calculateLatticeWeight(lattice [][]vertex, rows, columns int, inputs, outputs []bool, weightFn func(vertex, int, int) []int) []int {
	weight := make([]int, rows+1)
	weight[0] = 0
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			newWeight := weightFn((lattice[i][j]), i, rows)
			for idx := 0; idx < rows+1; idx++ {
				if idx == 0 && newWeight[idx] > 0 {
					weight[idx] = 1
				} else {
					weight[idx] += newWeight[idx]
				}
			}
		}
	}
	return weight
}

var boltzmannWeight = func(v vertex, row, n int) []int {
	var weights []int = make([]int, n+1)
	// Horizontal
	if v.East == in && v.West == out {
		weights[row+1] = 1
	} else if v.West == out && v.South == in { // Left turn
		weights[row+1] = 1
	} else if v.North == out && v.East == in {
		weights[0] = 1
	}
	return weights
}

func polynomialToString(poly []int) string {
	str := ""
	if poly[0] == 0 {
		return "0"
	}

	single := false
	for i, x := range poly {
		if x == 0 {
			continue
		}
		if i == 0 {
			single = (x == 1)
			str += fmt.Sprintf("%d", x)
			continue
		}
		if x == 1 {
			str += fmt.Sprintf("x%d", i-1)
		} else {
			str += fmt.Sprintf("x%d^%d", i-1, x)
		}
	}

	if single && len(str) > 1 {
		str = strings.TrimPrefix(str, "1")
	}
	return str
}

func allVertexL(lattice [][]vertex, row, col, rows, columns int, inputs, outputs []bool) []vertex {
	all := []vertex{
		{North: true, South: true, East: false, West: false},
		{North: true, South: false, East: true, West: false},
		{North: false, South: true, East: false, West: true},
		{North: false, South: false, East: true, West: true},
		{North: false, South: true, East: true, West: false},
	}

	// make sure bottom edge lines up
	if row > 0 {
		new := make([]vertex, 0, len(all))
		for _, v := range all {
			if v.South != lattice[row-1][col].North {
				new = append(new, v)
			}
		}
		all = new
	}
	// make sure right edge lines up
	if col > 0 {
		new := make([]vertex, 0, len(all))
		for _, v := range all {
			if v.East != lattice[row][col-1].West {
				new = append(new, v)
			}
		}
		all = new
	}

	if col == columns-1 {
		new := make([]vertex, 0, len(all))
		for _, v := range all {
			if v.West == in {
				new = append(new, v)
			}
		}
		all = new
	}
	if row == 0 {
		new := make([]vertex, 0, len(all))
		for _, v := range all {
			if v.South == out {
				new = append(new, v)
			}
		}
		all = new
	}
	if col == 0 {
		new := make([]vertex, 0, len(all))
		for _, v := range all {
			if v.East != inputs[row] {
				new = append(new, v)
			}
		}
		all = new
	}
	if row == rows-1 {
		new := make([]vertex, 0, len(all))
		for _, v := range all {
			if v.North == outputs[col] {
				new = append(new, v)
			}
		}
		all = new
	}
	return all
}

func tryValidLattices(lattice [][]vertex, n int, rows, columns int, inputs, outputs []bool) [][]int {
	if n == rows*columns {
		// Done.
		if validate(lattice, rows, columns, inputs, outputs) {
			weight := calculateLatticeWeight(lattice, rows, columns, inputs, outputs, boltzmannWeight)
			if printEachLattice {
				printLattice(lattice, rows, columns)
				fmt.Printf("\nboltzmann-weight = %s\n\n\n", polynomialToString(weight))
			}
			return [][]int{weight}
		}
		return [][]int{}
	}
	row := n / columns
	col := n % columns
	allWeights := make([][]int, 0)
	for _, v := range allVertexL(lattice, row, col, rows, columns, inputs, outputs) {
		lattice[row][col] = v
		allWeights = append(allWeights, tryValidLattices(lattice, n+1, rows, columns, inputs, outputs)...)
	}
	return allWeights
}

func findPartitionFunction(rows, columns int, inputs, outputs []bool) [][]int {
	var lattice [][]vertex = make([][]vertex, rows, rows)
	for r := 0; r < rows; r++ {
		lattice[r] = make([]vertex, columns)
	}

	return tryValidLattices(lattice, 0, rows, columns, inputs, outputs)
}

func printRow(row []vertex, r, columns int) {
	for col := columns - 1; col >= 0; col-- {
		fmt.Printf("  ")
		if row[col].North == out {
			fmt.Printf("^")
		} else {
			fmt.Printf("v")
		}
	}
	fmt.Println()
	fmt.Printf(">")

	for col := columns - 1; col >= 0; col-- {
		if row[col].West == out {
			fmt.Printf("<")
		} else {
			fmt.Printf(">")
		}
		fmt.Printf("x")
		if row[col].East == in {
			fmt.Printf("<")
		} else {
			fmt.Printf(">")
		}

	}
	fmt.Printf(" x%d\n", r)

	for col := columns - 1; col >= 0; col-- {
		fmt.Printf("  ")
		if row[col].South == out {
			fmt.Printf("v")
		} else {
			fmt.Printf("^")
		}
	}
	fmt.Println()

}

func printLattice(lattice [][]vertex, rows, columns int) {
	for i := 0; i < columns; i++ {
		fmt.Printf("  %d", columns-(i+1))
	}
	fmt.Println()
	for i := rows - 1; i >= 0; i-- {
		printRow(lattice[i], i, columns)
	}
}

func intSliceEqual(a, b []int, start int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := start; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func simplify(poly [][]int) [][]int {
	output := make([][]int, 0)

	for i := 0; i < len(poly); i++ {
		output = append(output, poly[i])
		for j := i + 1; j < len(poly); j++ {
			if intSliceEqual(poly[i], poly[j], 1) {
				output[len(output)-1][0] += poly[j][0]
				poly = append(poly[:j], poly[j+1:]...)
				j--
			}
		}
	}
	return output
}
