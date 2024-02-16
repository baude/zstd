package main

// "github.com/klauspost/compress/zstd"
import (
	"compress/gzip"
	"fmt"
	"os"
	"syscall"
	"time"

	crcOs "github.com/crc-org/crc/v2/pkg/os"
	"github.com/klauspost/compress/zstd"
)

func main() {

	fmt.Println("hello")
	fmt.Println("doing zstd")
	if err := doZstd(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("doing gz")
	if err := doGz(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("done")
	os.Exit(0)

}

func doGz() error {
	now := time.Now()
	uncompressedPath := fmt.Sprintf("output-gz-%v-%v-%v-%v.raw", now.Day(), now.Hour(), now.Minute(), now.Second())
	f, err := os.Open("fedora-coreos-39.20240210.2.0-applehv.aarch64.raw.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	dstFile, err := os.OpenFile(uncompressedPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	_, err = crcOs.CopySparse(dstFile, gzReader)
	if err != nil {
		return err
	}

	if err := dstFile.Sync(); err != nil {
		return err
	}
	return sparseCheck(dstFile.Fd())
}

func sparseCheck(fd uintptr) error {
	info := syscall.Stat_t{}
	if err := syscall.Fstat(int(fd), &info); err != nil {

	}
	fmt.Printf("SparseSize: %d    ActualSize: %d    isSparse: %v\n", info.Blocks*512, info.Size, info.Blocks*512 < info.Size)
	return nil
}

func doZstd() error {
	now := time.Now()
	uncompressedPath := fmt.Sprintf("output-zstd-%v-%v-%v-%v.raw", now.Day(), now.Hour(), now.Minute(), now.Second())
	f, err := os.Open("fedora-coreos-39.20240210.2.0-applehv.aarch64.raw.zst")
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
	return sparseCheck(dstFile.Fd())
}
