package main

import (
	"flag"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CommandOption struct {
	flagSet   *flag.FlagSet
	delimiter string
	languages string
}

func main() {
	options := &CommandOption{}
	options.flagSet = flag.NewFlagSet("ocr2csv", flag.ExitOnError)
	options.flagSet.StringVar(&options.languages, "langs", "eng", "Comma-separated language codes for OCR")
	options.flagSet.StringVar(&options.delimiter, "delimiter", ",", "Delimiter to separate OCR outputs")
	options.flagSet.Parse(os.Args[1:])

	args := options.flagSet.Args()
	if len(args) < 1 {
		log.Fatal("Directory name is required as first argument")
	}
	dirPath := args[0]

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(files, func(i, j int) bool {
		nameI, err := files[i].Info()
		if err != nil {
			log.Fatal(err)
		}
		nameJ, err := files[j].Info()
		if err != nil {
			log.Fatal(err)
		}
		return nameI.Name() < nameJ.Name()
	})

	ocrClient := gosseract.NewClient()
	defer ocrClient.Close()

	langs := strings.Split(options.languages, ",")
	if err := ocrClient.SetLanguage(langs...); err != nil {
		log.Fatalf("Failed to set languages: %v", err)
	}

	ocrPrint := func(ocrClient *gosseract.Client, filePath string) error {
		ocrClient.SetImage(filePath)
		text, err := ocrClient.Text()
		fmt.Printf("%s", strings.ReplaceAll(strings.ReplaceAll(text, "\n", "\\n"), ",", "\\,"))
		return err
	}

	var lastPrefix string
	var fileIndex int
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		name := info.Name()
		if filepath.Ext(name) == ".png" {
			prefix := strings.Split(name, "-")[0]
			filePath := filepath.Join(dirPath, name)
			if fileIndex > 0 && prefix != lastPrefix {
				fmt.Printf("\n")
			}
			if fileIndex > 0 && prefix == lastPrefix {
				fmt.Printf("%s", options.delimiter)
			}
			if err := ocrPrint(ocrClient, filePath); err != nil {
				log.Fatal(err)
			}
			lastPrefix = prefix
			fileIndex++
		}
	}
}
