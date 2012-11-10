package goparallel

import (
//	"math"
//	"fmt"
//	"math/rand"
//	"log"
//	"runtime"
)

var (
)

const (
)

type InputItem interface{}
type OutputItem interface{}
type Task func(input InputItem) OutputItem

func DoParallel(task Task, inputs []InputItem, maxCores int) []OutputItem {
	
}
