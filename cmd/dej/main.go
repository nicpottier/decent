package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/nicpottier/decent/parser"
	_ "github.com/nicpottier/decent/types"
)

func errorJSON(err error) map[string]string {
	return map[string]string{
		"type":  "error",
		"error": err.Error(),
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	out := json.NewEncoder(os.Stdout)

	for {
		mt, mb, err := parser.ReadNextToken(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			out.Encode(errorJSON(err))
			continue
		}

		m, err := parser.ParseMessage(mt, mb)
		if err != nil {
			out.Encode(errorJSON(err))
			continue
		}

		out.Encode(m)
	}
}
