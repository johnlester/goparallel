package goparallel

import (
	"fmt"
	"runtime"
	"math/big"
	// "time"
	//	"math"
	//	"math/rand"
	//	"log"
)

var (
	// numCores = runtime.NumCPU()
)

const (
)

// type InputItem interface{}
// type OutputItem interface{}

type InputItem struct {
	n *big.Int
}

type Job struct {
    jobID int
	input InputItem
	resultsChan  chan<- Result
}	

type OutputItem struct {
	n *big.Int
	factors []*big.Int
}

type Result struct {
	jobID int
	output OutputItem
}



//Main function to do parallel computation
func DoParallel(inputs []InputItem, maxWorkers int) []OutputItem {
	
	//set num of workers
	numWorkers := maxWorkers
	
	fmt.Printf("Using %v of %v cores, with %v workers\n", runtime.GOMAXPROCS(0), runtime.NumCPU(), numWorkers)

	numOutputs := len(inputs)
	
	//Create channels
    jobsChan := make(chan Job, numWorkers)
    resultsChan := make(chan Result, len(inputs))
    doneChan := make(chan struct{}, numWorkers)

    go addJobs(jobsChan, inputs, resultsChan) // Executes in its own goroutine
    for i := 0; i < numWorkers; i++ {
        go doJobs(doneChan, jobsChan) // Each executes in its own goroutine
    }
    go awaitCompletion(doneChan, resultsChan, numWorkers) // Executes in its own goroutine
    finalOutputs := processResults(resultsChan, numOutputs)           // Blocks until the work is done
	return finalOutputs
}

func addJobs(jobsChan chan<- Job, inputs []InputItem, resultsChan chan<- Result) {
    for i, input := range inputs {
        jobsChan <- Job{i, input, resultsChan}
    }
    close(jobsChan)
}

func doJobs(doneChan chan<- struct{}, jobsChan <-chan Job) {
    for job := range jobsChan {
        job.Do()	//job.Do() doesn't need any parameters in this case
    }
    doneChan <- struct{}{}
}

func (job Job) Do() {
	output := OutputItem{job.input.n, factorsBig(job.input.n)}
	job.resultsChan <- Result{job.jobID, output}
}

func awaitCompletion(doneChan <-chan struct{}, resultsChan chan Result, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		<-doneChan
	}
	close(resultsChan)
}

func processResults(resultsChan <-chan Result, numOutputs int) []OutputItem {
	finalOutputs := make([]OutputItem, numOutputs)
	for result := range resultsChan {
		fmt.Printf("%v) %v: %v\n", result.jobID, result.output.n, result.output.factors)
		finalOutputs[result.jobID] = result.output
	}
	return finalOutputs
}


//////////////////////////////////////////////////
// Factoring functions
//////////////////////////////////////////////////

// Sqrt for Big numbers
func sqrtBig(n *big.Int) (x *big.Int) {
	switch n.Sign() {
	case -1:
		panic(-1)
	case 0:
		return big.NewInt(0)
	}

	var px, nx big.Int
	x = big.NewInt(0)
	x.SetBit(x, n.BitLen()/2+1, 1)
	for {
		nx.Rsh(nx.Add(x, nx.Div(n, x)), 1)
		if nx.Cmp(x) == 0 || nx.Cmp(&px) == 0 {
			break
		}
		px.Set(x)
		x.Set(&nx)
	}
	return
}

func factorsBig(n *big.Int) []*big.Int {
	factors := make([]*big.Int, 2, 10)
	factors[0] = big.NewInt(1)
	factors[1] = n
	zero := big.NewInt(0)
	one := big.NewInt(1)
	i := big.NewInt(2)
	sqrtN := sqrtBig(n)
	for i.Cmp(sqrtN) < 1 {
		z := big.NewInt(0)
		m := big.NewInt(0)
		div := big.NewInt(0)
		div = div.Add(div, i)
		z, m = z.DivMod(n, div, m)
		if m.Cmp(zero) == 0	{   //modulus equals zero
			factors = append(factors, div)
			factors = append(factors, z)
		}
		i = i.Add(i, one)
	}
	return factors
}

