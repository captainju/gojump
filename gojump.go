package main

import (
	"flag"
	"bufio"
	"os"
	"strconv"
	"strings"
	"sort"
	"os/user"
	"fmt"
	"errors"
	"path/filepath"
)

const fileName string = ".gojump.txt"
const separator string = " <-> "

type Wabon struct {
	path   string
	weight int
}

type BunchWabon []Wabon

func (b BunchWabon) save() {
	sort.Sort(b)

	datafileName, err := dataFile()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(datafileName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, wabon := range b {
		writer.WriteString(strconv.Itoa(wabon.weight) + separator + wabon.path + "\n")
	}
	writer.Flush()
}

func (b *BunchWabon) load() error {
	datafileName, err := dataFile()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(datafileName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line, separator)
		w, e := strconv.Atoi(s[0])
		if e == nil {
			*b = append(*b, Wabon{path:s[1], weight:w})
		}
	}
	return nil
}

func (b BunchWabon) compl(keyword string, howMany int) []Wabon {
	var foundWabons BunchWabon
	lowerKeyword := strings.ToLower(keyword)
	for _, wabon := range b {
		if strings.Contains(strings.ToLower(wabon.path), lowerKeyword) {
			foundWabons = append(foundWabons, wabon)
		}
	}
	sort.Sort(foundWabons)
	if len(foundWabons) > howMany {
		foundWabons = foundWabons[0:howMany]
	}
	return foundWabons
}

func (b BunchWabon) findOne(keyword string) (Wabon, error) {
	var foundWabon Wabon
	lowerKeyword := strings.ToLower(keyword)
	for _, wabon := range b {
		if strings.Contains(strings.ToLower(wabon.path), lowerKeyword) {
			if foundWabon.weight < wabon.weight {
				foundWabon = wabon
			}
		}
	}
	if foundWabon.path != "" {
		return foundWabon, nil
	}
	return foundWabon, errors.New("not found")
}

func (b *BunchWabon) addPath(path string) {
	path = filepath.Clean(path)
	if path == "." {
		return
	}

	if exists, err := exists(path); err != nil || !exists {
		return
	}

	isPresent := false
	for num, wabon := range *b {
		if wabon.path == path {
			isPresent = true
			wabon.weight = wabon.weight + 1
			(*b)[num] = wabon
			break
		}
	}

	if !isPresent {
		*b = append(*b, Wabon{path:path, weight:1})
	}

	b.save()
}

func (b *BunchWabon) purge() {
	newBunch := BunchWabon{}
	for _, wabon := range *b {
		if ok, _ := exists(wabon.path); ok {
			newBunch = append(newBunch, wabon)
		}
	}
	b = &newBunch
	b.save()
}

func (b BunchWabon) Len() int {
	return len(b)
}
func (b BunchWabon) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b BunchWabon) Less(i, j int) bool {
	return b[i].weight > b[j].weight
}

func main() {
	addPtr := flag.Bool("add", false, "add a path to the index")
	listPtr := flag.Bool("list", false, "list all entries")
	complPtr := flag.Bool("compl", false, "completion")
	purgePtr := flag.Bool("purge", false, "purge non existing paths")
	flag.Parse()

	path := flag.Args()
	bunchOfWabons := BunchWabon{}
	err := bunchOfWabons.load()
	if err != nil {
		panic(err)
	}

	if *addPtr && len(path) != 0 {
		bunchOfWabons.addPath(path[0])
	} else if *listPtr {
		for _, wabon := range bunchOfWabons {
			fmt.Printf("%d : %s\n", wabon.weight, wabon.path)
		}
	} else if *purgePtr {
		bunchOfWabons.purge()
	} else if *complPtr && len(path) != 0 {
		foundWabons := bunchOfWabons.compl(path[0], 5)
		for _, wabon := range foundWabons {
			fmt.Println(wabon.path)
		}
	} else if len(path) != 0 {
		foundWabon, err := bunchOfWabons.findOne(path[0])
		if err != nil {
			os.Exit(1)
		}
		fmt.Print(foundWabon.path)
	}
}

func dataFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir + string(os.PathSeparator) + fileName, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
