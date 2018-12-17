// @author Meghan#2032 <https://nektro.net/> 2018
//
// go run main.go "C:/uspto/"
//
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please pass a directory as a parameter")
		return
	}
	os.Mkdir(os.Args[1], os.ModePerm)
	folder := os.Args[1]
	if !strings.HasSuffix(folder, "/") && !strings.HasSuffix(folder, "\\") {
		folder += string(os.PathSeparator)
	}

	a := 0
	b := 0
	c := 0

	var wg sync.WaitGroup
	for i := 0; i < 10000000; i++ {
		wg.Add(100)
		for j := 0; j < 100; j++ {
			go saveFile(&wg, folder, i, a, b, c)
			a++
			if a == 99 {
				a = 0
				b++
			}
			if b == 999 {
				b = 0
				c++
			}
		}
		wg.Wait()
	}
}

func saveFile(wg *sync.WaitGroup, folder string, i int, a int, b int, c int) {
	x := padLeft(fmt.Sprintf("%d", a), 2, "0")
	y := padLeft(fmt.Sprintf("%d", b), 3, "0")
	z := padLeft(fmt.Sprintf("%d", c), 3, "0")
	url := fmt.Sprintf("http://pimg-fpiw.uspto.gov/fdd/%s/%s/%s/0.pdf", x, y, z)
	file := folder + padLeft(fmt.Sprintf("%d", i), 8, "0") + ".pdf"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile(file, body, os.ModePerm)
			fmt.Println("+ Saving:  " + url)
		} else {
			fmt.Println("- Skipped: " + url)
		}
	} else {
		fmt.Println("- Passed:  " + url)
	}
	wg.Done()
}

func padLeft(_this string, desiredLength int, pad string) string {
	if len(_this) >= desiredLength {
		return _this
	}
	right := min(len(pad), desiredLength-len(_this))
	padding := pad[0:right]
	return padLeft(padding+_this, desiredLength, pad)
}

func min(x int, y int) int {
	if x > y {
		return y
	}
	return x
}
