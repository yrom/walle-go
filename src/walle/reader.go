package walle

import (
	"fmt"
)

func PrintChannel(files []string) {
	processAllFiles(files,
		func(c ChannelInfo) string {
			return "channel=" + c.Channel
		})
}

func PrintRaw(files []string) {
	processAllFiles(files,
		func(c ChannelInfo) string {
			return c.String()
		})
}

// Iterate over files with block consumer function
func processAllFiles(files []string, process func(ChannelInfo) string) {
	for _, file := range files {
		if !isRegularFile(file) {
			fmt.Printf("%s is not a regular file!\n", file)
			continue
		}
		info, err := readChannelInfo(file)
		if err != nil {
			fmt.Printf("Error occured on reading file %s, %s\n", file, err)
			continue
		}
		result := process(info)
		fmt.Printf("%s : %s\n", file, result)
	}
}
