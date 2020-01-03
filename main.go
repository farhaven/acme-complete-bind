package main

import (
	"log"
	"bytes"
	"os"
	"os/exec"
	"fmt"

	"9fans.net/go/acme"
)

const _WinID = 97

func main() {
	log.Println("Here we go")
	w, err := acme.Open(_WinID, nil)
	if err != nil {
		log.Fatalln("can't open window", _WinID, ":", err)
	}

	events := w.EventChan()

	for e := range events {
		if e.C1 != 'K' {
			// Not a "keyboard" event
			continue
		}
		if e.C2 != 'I' {
			// Not an "input" event
			continue
		}


		if !bytes.Equal(e.Text, []byte{0x0f}) {
			continue
		}
		// ^O entered

		log.Printf("Text entered: %#+v", e)

		// Remove entered character from text, then run completion
		w.Addr("#%d,#%d", e.Q0, e.Q1)
		w.Write("data", []byte(""))

		// Run `L comp -e` in "From ACME" mode by setting the $winid appropriately
		lpath, err := exec.LookPath("L")
		if err != nil {
			log.Println("Can't find `L` command:", err)
			continue
		}
		lenv := os.Environ()
		lenv = append(lenv, fmt.Sprintf("winid=%d", _WinID))
		cmd := exec.Cmd{
			Env: lenv,
			Path: lpath,
			Args: []string{"L", "comp", "-e"},
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Failed to run L command:", err)
			log.Println("Output of L:", string(output))
			continue
		}
	}
}