package main

// "github.com/klauspost/compress/zstd"
import (
	"fmt"
	"os"
	"strings"
	"time"

	crcOs "github.com/crc-org/crc/v2/pkg/os"
	"github.com/klauspost/compress/zstd"
)

func main() {
	fmt.Println("hello")
	if err := doSomething(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("done")
	os.Exit(0)

}

func doSomething() error {
	outputfile := strings.ReplaceAll(time.Now().String(), " ", "")
	uncompressedPath := fmt.Sprintf("output-%v.raw", outputfile)
	f, err := os.Open("6edf958c5594bebd8050b58b8317f3f4c97c1f534a2757c3db5bfa33bdc64824.raw.zst")
	if err != nil {
		return err
	}
	defer f.Close()

	dstFile, err := os.OpenFile(uncompressedPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	zstdReader, err := zstd.NewReader(f)
	if err != nil {
		return err
	}
	defer zstdReader.Close()

	_, err = crcOs.CopySparse(dstFile, zstdReader)
	if err != nil {
		return err
	}

	return nil
}
