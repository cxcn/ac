package ac

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

const (
	defaultDictPath    = "/ramdisk/vonng/dict.txt"
	defaultInputPath   = "/ramdisk/vonng/video_title.txt"
	defaultOutputPath  = "/ramdisk/vonng/vonng.txt"
	defaultProfilePath = "/ramdisk/vonng/profile.txt"
)

// args
var (
	dictPath    string // dict path
	inputPath   string // input path
	outputPath  string // output path
	profilePath string // profile path
)

/**************************************************************\
*                          Correctness                         *
\**************************************************************/
// Check run program once ,do prof and check sig
func CheckAC() {
	begin := time.Now()
	Run(inputPath, outputPath, dictPath)
	elapse := time.Now().Sub(begin)

	// Check correctness
	hasher := md5.New()
	body, err := ioutil.ReadFile(outputPath)
	if err != nil {
		panic(err)
	}
	hasher.Write(body)
	sig := hex.EncodeToString(hasher.Sum(nil))

	fmt.Printf("time: %s sig: %s\n", elapse, sig)
}

/**************************************************************\
*                          Performance                         *
\**************************************************************/
// Bench run program with n times. drop the max and min result
func Benchmark(n int) {
	// pprof
	if profilePath != "" {
		f, _ := os.Create(profilePath)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// bench
	bench := make([]int64, n)
	for i := 0; i < n; i++ {
		begin := time.Now()
		Run(inputPath, outputPath, dictPath)
		elapse := int64(time.Now().Sub(begin))
		fmt.Printf("Round %d: %s\n", i, time.Duration(elapse))
		bench[i] = elapse
	}

	// sort and drop min & max
	for i := 0; i < n; i++ {
		for j := i; j > 0 && bench[j] < bench[j-1]; bench[j], bench[j-1], j = bench[j-1], bench[j], j-1 {
		}
	}
	bench = bench[1 : len(bench)-1]

	// avg
	var sum int64
	for i := 0; i < len(bench); i++ {
		sum += bench[i]
	}
	avg := time.Duration(sum / int64(len(bench)))
	fmt.Printf("Avg: %s\n", avg)
}

/**************************************************************\
*                          Driver                              *
\**************************************************************/

func TestAC(t *testing.T) {
	flag.StringVar(&dictPath, "d", defaultDictPath, "dict filename")
	flag.StringVar(&inputPath, "i", defaultInputPath, "input filename")
	flag.StringVar(&outputPath, "o", defaultOutputPath, "output filename")
	flag.StringVar(&profilePath, "p", "", "profile filename")
	flag.Parse()

	CheckAC()
	Benchmark(10)
}
