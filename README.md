# Ice Ice Baby

## Installation instructions

* **Step 1.** If you haven't done so already, install [Go](https://golang.org)
* **Step 2.** Clone this repository:
```
$ git clone https://github.com/suzmue/iceicebaby.git
```
* **Step 3.** Open `main.go`. 

In the function `main` you can input the number of rows, columns, and which rows and columns will serve as inputs and outputs to the `findPartitionFunction` function.

If you want to print out each lattice, set `var printEachLattice bool = true`.
* **Step 4.** You can run this program by typing:

```
go run main.go
```

If you want to direct the output to a file, run
```
go run main.go > file.txt
```
