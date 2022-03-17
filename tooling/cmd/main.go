// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	var b strings.Builder

	// Scan standard input for log data per line.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		// Convert the JSON to a map for processing.
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			fmt.Println(err)
		}

		// Having a traceid present in the logs.
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["traceid"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// Build out the know portions of the log in the order
		// I want them in.
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["ts"],
			m["level"],
			traceID,
			m["caller"],
			m["msg"],
		))

		// Add the rest of the keys ignoring the ones we already
		// added for the log.
		for k, v := range m {
			switch k {
			case "service", "ts", "level", "traceid", "caller", "msg":
				continue
			}

			// It's nice to see the key[value] in this format
			// especially since map ordering is random.
			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		// Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
