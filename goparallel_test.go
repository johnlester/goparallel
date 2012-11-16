package goparallel

import (
	"testing"
	"math/big"
	"fmt"
	"runtime"
)

const (
	bigsToFactor = 100
	// startingBigString = "1234567891011121314151617181920"
	startingBigString = "248163264"
)

var (
	bigOne = big.NewInt(1)
	bigTwo = big.NewInt(2)
)


/////////////////////////////////////////////////////////////
// goParallel tests
/////////////////////////////////////////////////////////////

func TestGoParallel(t *testing.T) {

	//Create big to start
	start := new(big.Int)
	_, err := fmt.Sscan(startingBigString, start)	
	if err != nil {
		t.Errorf("error scanning big string:", err)
	}

	//Create inputs
	inputs := make([]InputItem, bigsToFactor)
	for i := 0; i < bigsToFactor; i++ {
		nextBig := big.NewInt(int64(i))
		nextBig = nextBig.Add(nextBig, start)
		inputs[i] = InputItem{nextBig}
	}

	//Do computation in parallel
	runtime.GOMAXPROCS(runtime.NumCPU())
	var outputs []OutputItem
	outputs = DoParallel(inputs, 5)
	
	//Tests
	if len(inputs) != len(outputs) {
		t.Errorf("inputs (%v) should have same len as outputs (%v)", len(inputs), len(outputs))
	}

}




/////////////////////////////////////////////////////////////
// bif factor tests
/////////////////////////////////////////////////////////////

func TestGoParallel_Sqrt(t *testing.T) {
	four := big.NewInt(4)
	two := big.NewInt(2)
	if sqrtBig(four).Cmp(two) != 0 {
		t.Errorf("mo money, mo problems")
	}
	
	twentyFive := big.NewInt(25)
	five := big.NewInt(5)
	if sqrtBig(twentyFive).Cmp(five) != 0 {
		t.Errorf("mo money, mo problems")
	}
	
	oneOhOne := big.NewInt(101)
	ten := big.NewInt(10)
	if sqrtBig(oneOhOne).Cmp(ten) != 0 {
		t.Errorf("mo money, mo problems")
	}

	z := big.NewInt(0)
	m := big.NewInt(0)
	z, m = z.DivMod(oneOhOne, ten, m)
	if m.Cmp(bigOne) != 0 {
		t.Errorf("mo money, mo problems")
	}

}

