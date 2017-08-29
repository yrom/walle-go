package walle

import (
	"os"
	"path/filepath"
	"fmt"
	"time"
)

type zipSections struct {
	beforeSigningBlock []byte
	signingBlock       []byte
	signingBlockOffset int64
	centraDir          []byte
	centralDirOffset   int64
	eocd               []byte
	eocdOffset         int64
}
type transform func(*zipSections) (*zipSections, error)

func (z *zipSections) writeTo(output string, transform transform) (err error) {
	f, err := os.Create(output)
	if err != nil {
		return
	}

	defer f.Close()

	newZip, err := transform(z)
	if err != nil {
		return
	}

	for _, s := range [][]byte{
			newZip.beforeSigningBlock,
			newZip.signingBlock,
			newZip.centraDir,
			newZip.eocd} {
		_, err := f.Write(s)
		if err != nil {
			return err
		}
	}
	return
}

var _debug bool

func GenerateChannelApk(out string, channels []string, extras map[string]string, input string, force bool, debug bool) {
	_debug = debug
	if len(input) == 0 {
		exit("Error: no input file specified!")
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		exitf("Error: no such file %s!", input)
	}

	if len(out) == 0 {
		out = filepath.Dir(input)
	} else {
		fi, err := os.Stat(out)
		if os.IsNotExist(err) || !fi.IsDir() {
			exitf("Error: output %s is neither exist nor a dir!", out)
		}
	}
	if len(channels) == 0 {
		exit("Error: no channel specified!")
	}
	//TODO: add new option for generating new channel from channelled apk
	if c, _ := readChannelInfo(input); len(c.Channel) != 0 {
		exitf("Error: file %s is registered a channel block %s", filepath.Base(input), c.String())
	}
	var start time.Time
	if _debug {
		start = time.Now()
	}

	fmt.Printf("Generating channels %s for %s into dir %s ...\n", channels, filepath.Base(input), out)
	z, err := newZipSections(input)
	if err != nil {
		exitf("Error occurred on parsing apk %s, %s\n", input, err)
	}
	name, ext := fileNameAndExt(input)
	for _, channel := range channels {
		output := filepath.Join(out, name+"-"+channel+ext)
		c := ChannelInfo{Channel: channel, Extras: extras}
		err = gen(c, z, output, force)
		if err != nil {
			exitf("Error occurred on generating channel %s, %s\n", channel, err)
		}
	}
	if _debug {
		println("Consume", time.Since(start).String())
	} else {
		println("Done!")
	}

}

func newZipSections(input string) (z zipSections, err error) {
	in, err := os.Open(input)
	if err != nil {
		return
	}
	defer in.Close()

	// read eocd
	eocd, eocdOffset, err := findEndOfCentralDirectoryRecord(in)
	if err != nil {
		return
	}
	centralDirOffset := getEocdCentralDirectoryOffset(eocd)
	centralDirSize := getEocdCentralDirectorySize(eocd)
	z.eocd = eocd
	z.eocdOffset = eocdOffset
	z.centralDirOffset = int64(centralDirOffset)

	// read signing block
	signingBlock, signingBlockOffset, err := findApkSigningBlock(in, centralDirOffset)
	if err != nil {
		return
	}
	z.signingBlock = signingBlock
	z.signingBlockOffset = signingBlockOffset
	// read bytes before signing block
	//TODO: waste too large memory
	if signingBlockOffset >= 64 * 1024 * 1024 {
		fmt.Print("Warning: maybe waste large memory on processing this apk! ")
		fmt.Println("Before APK Signing Block bytes size is", signingBlockOffset/1024/1024, "MB")
	}
	beforeSigningBlock := make([]byte, signingBlockOffset)
	n, err := in.ReadAt(beforeSigningBlock, 0)
	if err != nil {
		return
	}
	if int64(n) != signingBlockOffset {
		return z, fmt.Errorf("Read bytes count mismatched! Expect %d, but %d", signingBlockOffset, n)
	}
	z.beforeSigningBlock = beforeSigningBlock

	centralDir := make([]byte, centralDirSize)
	n, err = in.ReadAt(centralDir, int64(centralDirOffset))
	if uint32(n) != centralDirSize {
		return z, fmt.Errorf("Read bytes count mismatched! Expect %d, but %d", centralDirSize, n)
	}
	z.centraDir = centralDir
	if _debug {
		fmt.Printf("signingBlockOffset=%d, signingBlockLenth=%d\n"+
			"centralDirOffset=%d, centralDirSize=%d\n"+
			"eocdOffset=%d, eocdLenthe=%d\n",
			signingBlockOffset,
			len(signingBlock),
			centralDirOffset,
			centralDirSize,
			eocdOffset,
			len(eocd))
	}
	return
}

func gen(info ChannelInfo, sections zipSections, output string, force bool) (err error) {

	fi, err := os.Stat(output)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if fi != nil {
		if !force {
			return fmt.Errorf("file already exists %s.", output)
		}
		println("Force generating channel", info.Channel)
	}

	var s time.Time
	if _debug {
		s = time.Now()
	}
	err = sections.writeTo(output, newTransform(info))
	if _debug {
		fmt.Printf("    write %s consume %s", output, time.Since(s).String())
		fmt.Println()
	}
	return
}

func newTransform(info ChannelInfo) transform {
	return func(zip *zipSections) (*zipSections, error) {

		newBlock, diffSize, err := makeSigningBlockWithChannelInfo(info, zip.signingBlock)
		if err != nil {
			return nil, err
		}
		newzip := new(zipSections)
		newzip.beforeSigningBlock = zip.beforeSigningBlock
		newzip.signingBlock = newBlock
		newzip.signingBlockOffset = zip.signingBlockOffset
		newzip.centraDir = zip.centraDir
		newzip.centralDirOffset = zip.centralDirOffset
		newzip.eocdOffset = zip.eocdOffset
		newzip.eocd = makeEocd(zip.eocd, uint32(int64(diffSize)+zip.centralDirOffset))
		return newzip, nil
	}
}
