package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type recordHandler func([]string, *TCUFile)

//TCUFile structure to store a little data about the file.
type TCUFile struct {
	header          []string
	recordHandlers  map[string]recordHandler
	csvWriter       *csv.Writer
	reader          *bufio.Reader
	includeComments bool
}

//NewTCUFile init a new object
func NewTCUFile(comments bool) *TCUFile {
	rt := map[string]recordHandler{
		"Clock":     discardHandler,
		"Date":      discardHandler,
		"Time":      timeHandler,
		"Connector": connectorHandler,
		"Trigger":   triggerHandler,
		"Delay":     delayHandler,
		"Observer":  discardHandler,
		"Record":    discardHandler,
		"String":    discardHandler,
		"Timeout":   discardHandler,
	}

	return &TCUFile{
		header:          []string{"", "", "", "", "", "", "", ""},
		recordHandlers:  rt,
		includeComments: comments,
	}
}

var inputFile, outputFile *string
var comments *bool
var logger *TCUFile

func init() {
	inputFile = flag.String("i", "", "Input File")
	outputFile = flag.String("o", "", "Output File")
	comments = flag.Bool("c", false, "Include comments")
	flag.Parse()

	if *inputFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger = NewTCUFile(*comments)

}

func main() {

	//outputFile writer
	if *outputFile == "" {
		logger.csvWriter = csv.NewWriter(os.Stdout)
	} else {
		fo, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal("Error creating output file:", err)
		}
		defer func() {
			logger.csvWriter.Flush()
			fo.Close()
		}()
		logger.csvWriter = csv.NewWriter(fo)
	}

	//inputFile reader
	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal("Error creating input file:", err)
	}
	defer f.Close()
	logger.reader = bufio.NewReader(f)

	//iterate the file
	for {
		line, _, err := logger.reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}

		rec := strings.Fields(string(line))

		if len(rec) > 0 {

			if f, ok := logger.recordHandlers[rec[0]]; ok {
				f(rec, logger)
				continue
			}

			//TODO: handle converting hex to unixtime or timestring.
			logger.csvWriter.Write(rec)

		}
	}
}

func timeHandler(r []string, t *TCUFile) {
	/*
		T1: 3 ms, T2: 6 ms, T3: 12 ms, T4: 24 ms,
		T5: 60 ms, T6: 120 ms, T7: 300 ms, T8: 750 ms
	*/

	if t.includeComments {
		if r[0]+r[1] == "TimeSlice:" {
			record := []string{"#Time Slice:" + r[2]}
			t.csvWriter.Write(record)
		}
	}
}

func triggerHandler(r []string, t *TCUFile) {
	if t.includeComments {
		s := "#"
		for _, v := range r {
			if v == "," {
				v = " "
			}
			s = s + v
		}

		record := []string{s}
		t.csvWriter.Write(record)
	}
}

func delayHandler(r []string, t *TCUFile) {
	if t.includeComments {
		record := []string{"#Delay:" + r[2] + "%"}
		t.csvWriter.Write(record)
	}
}

func connectorHandler(r []string, t *TCUFile) {
	//$FSDEZ,$BFFHR,$BFBRM,$UNI4A,$INI4A,$ZPR,$VREF,SMETER
	line := r
	for i := 1; i <= 8; i++ {
		position := strings.Trim(line[1], ":")

		s, err := strconv.ParseInt(position, 10, 8)
		if err == nil {
			t.header[s-1] = line[2]
		}

		l, _, err := t.reader.ReadLine()
		if err != nil {
			log.Println("Error reading header info")
		}

		line = strings.Fields(string(l))
	}

	t.csvWriter.Write(t.header)
}

func dateHandler(r []string, t *TCUFile) {
	if t.includeComments {
		s := "#"
		for _, v := range r {
			s = s + v
		}
		record := []string{s}
		t.csvWriter.Write(record)
	}
}

func discardHandler(r []string, t *TCUFile) {
	//TODO: Should do something with this maybe....
}

func hexToInt() {
	/*s, e := strconv.ParseInt("570B395A", 16, 64)
	  if e != nil {
	  	fmt.Println(e)
	  }
	  fmt.Println(s)
	  t := time.Unix(s, 0)
	  t.Format("02/01/2006 15:04:05")
	  fmt.Println(t.Format("02/01/2006 15:04:05"))
	  fmt.Println(t)*/
}
