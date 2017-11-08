package service

import (
	"time"
	"math/rand"
	"strconv"
	"strings"
	"reflect"
)

import (
	"hackathon/battleship/schema/request"
)


func random(min, max int) int {
	max = max + 1
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func toKeyMap(x, y int) string {
	return strconv.Itoa(x) + "-" + strconv.Itoa(y)
}

func extractKeyMap(key_map string) []Vertex {
	key_arr := strings.Split(key_map, "-")
	var res []Vertex
	for i, item := range key_arr {
		if i%2 == 1 && i > 0 {
			x, _ := strconv.Atoi(key_arr[i-1])
			y, _ := strconv.Atoi(item)
			res = append(res, Vertex{X: x, Y: y})
		}
	}
	return res
}

func isEmpty(whole_ship request.RecognizedWholeShip) bool {
	return reflect.DeepEqual(whole_ship, request.RecognizedWholeShip{})
}

func isVertexEmpty(v Vertex) bool {
	return reflect.DeepEqual(v, Vertex{})
}

func getAllKeyMap(maps map[string][]Vertex) []string {
    keys := make([]string, 0, len(maps))
    for k := range maps {
        keys = append(keys, k)
    }
    return keys
}

func removeVertex(slice []Vertex, index int) []Vertex {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}

func checkItemInArray(array []string, item string) bool {
	for _, i := range array {
		if i == item {
			return true
		}
	}
	return false
}

func copyMap(m map[string]int) (res map[string]int) {
	res = make(map[string]int)
	for k, v := range m {
		res[k] = v
	}
	return
}
// func remove(slice *[]Gird, index int) {
// 	new_slice := append(slice[:index], slice[index+1:len(slice)]...)
// 	return new_slice
// }