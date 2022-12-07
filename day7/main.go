package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type File struct {
	Name string
	Size int
}
type Dir struct {
	Name string
	Parent *Dir
	Dirs map[string]*Dir
	Files []File
}
func NewDir(name string, parent *Dir) *Dir {
	d :=  &Dir{Name: name, Parent: parent}
	d.Dirs = make(map[string]*Dir)
	return d
}

func (d *Dir) TotalSize() (total int) {
	for _, f := range d.Files {
		total += f.Size
	}
	for _, dir := range d.Dirs {
		total += dir.TotalSize()
	}
	return
}

func (d *Dir) TotalSizeOfSmallDirs(max int) (total int) {
	//fmt.Printf("DEBUG: %s has dirs %v\n", d.Name, d.Dirs)
	for _, dir := range d.Dirs {
		total += dir.TotalSizeOfSmallDirs(max)
		if dir.TotalSize() <= max {
			fmt.Printf("INFO: counting dir %s with size %d\n", dir.Name, dir.TotalSize())
			total += dir.TotalSize()
		}
	}
	return total
}

func (d *Dir) FlattenDirs() (flattened []*Dir) {
	for _, dir := range d.Dirs {
		flattened = append(flattened, dir)
		flattened = append(flattened, dir.FlattenDirs()...)
	}
	return flattened
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))
	root := &Dir{Name: "/"}
	root.Parent = root
	root.Dirs = make(map[string]*Dir)

	var pwd *Dir
	for _, line := range strings.Split(string(data), "\n") {
		fmt.Printf("INFO: %s\n", line)
		if line[0] == '$' {
			// command
			if line[2:4] == "cd" {
				// change dir
				args := line[5:]
				if args == ".." {
					pwd = pwd.Parent
				} else if args == "/" {
					pwd = root
				} else {
					// danger: what if it isn't there?
					if _, ok := pwd.Dirs[args]; !ok {
						fmt.Printf("WARN: can't find dir %s in %s!\n", args, pwd.Name)
					}
					pwd = pwd.Dirs[args]
				}
			}
		} else {
			// output
			tokens := strings.Split(line, " ")
			if matched, _ := regexp.MatchString(`^\d.*`, line); matched {
				// file
				size, _ := strconv.Atoi(tokens[0])
				pwd.Files = append(pwd.Files, File{Name: tokens[1], Size: size})
			} else if strings.HasPrefix(line, "dir") {
				pwd.Dirs[tokens[1]] = NewDir(tokens[1], pwd)
			}
		}
	}

	fmt.Printf("Part 1: Total size of directories no bigger than 100000: %d\n",
		root.TotalSizeOfSmallDirs(100000))

	unused := 70000000 - root.TotalSize()
	required := 30000000 - unused

	fmt.Printf("Part 2: Current unused space is %d, update requires %d more\n", unused, required)

	smallestToDelete := root
	for _, dir := range root.FlattenDirs() {
		if dir.TotalSize() < smallestToDelete.TotalSize() && dir.TotalSize() >= required {
			smallestToDelete = dir
		}
	}
	fmt.Printf("Selected %s with size %d\n", smallestToDelete.Name, smallestToDelete.TotalSize())
}
