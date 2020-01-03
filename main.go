// Command acme-complete-bind adds a key binding for ^O to ACME to run acme-lsp's `L comp -e` command without having to move the mouse from the text.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"9fans.net/go/acme"
)

func handleCompletionEvent(win *acme.Win, event *acme.Event) error {
	// Remove entered character from text, then run completion
	win.Addr("#%d,#%d", event.Q0, event.Q1)
	win.Write("data", []byte(""))

	// Run `L comp -e` in "From ACME" mode by setting the $winid appropriately
	lpath, err := exec.LookPath("L")
	if err != nil {
		return err
	}
	lenv := os.Environ()
	lenv = append(lenv, fmt.Sprintf("winid=%d", win.ID()))
	cmd := exec.Cmd{
		Env:  lenv,
		Path: lpath,
		Args: []string{"L", "comp", "-e"},
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("can't run completion: %w (output: %s)", err, string(output))
	}
	// Write completion output to error window
	if len(output) != 0 {
		win.Err(string(output))
	}

	return nil
}

// dontHandleWindow returns true if the window should not be handled. This happens if the name of the window
// - indicates that it is a utility window (file name starts with - or +)
// - indicates that it is a directory (file name ends with /)
func dontHandleWindow(name string) bool {
	base := path.Base(name)
	return strings.HasPrefix(base, "-") || strings.HasPrefix(base, "+") || strings.HasSuffix(name, "/")
}

func handleWindow(winId int) {
	win, err := acme.Open(winId, nil)
	if err != nil {
		log.Printf("can't open window %d: %s", winId, err)
		return
	}
	defer win.CloseFiles()

	tag := make([]byte, 1024)
	_, err = win.Read("tag", tag)
	if err != nil {
		log.Println("can't read window tag:", err)
	} else {
		name := string(bytes.SplitN(tag, []byte(" "), 2)[0])
		if dontHandleWindow(string(name)) {
			return
		}
	}

	events := win.EventChan()

	log.Println("handling events for window", winId)
	for e := range events {
		if e.C1 == 'F' && e.C2 == 'i' && dontHandleWindow(string(e.Text)) {
			log.Println("window got turned into an utility or directory window, going away")
			return
		}

		// Re-enable automatic menu
		win.Write("ctl", []byte("menu"))

		if e.C1 == 'K' && e.C2 == 'I' && bytes.Equal(e.Text, []byte{0x0f}) {
			// ^O entered by user with keyboard
			err := handleCompletionEvent(win, e)
			if err != nil {
				log.Println("failed to handle completion request:", err)
			}
			continue
		}

		if e.C2 == 'x' || e.C2 == 'X' || e.C2 == 'l' || e.C2 == 'L' {
			// Tell ACME to handle the event itself
			err = win.WriteEvent(e)
			if err != nil {
				log.Println("can't write event:", err)
			}
			continue
		}
	}

	log.Println("event channel closed, done handling window", winId)
}

func main() {
	// Start key handler for all existing windows
	acmeWindows, err := acme.Windows()
	if err != nil {
		log.Fatalln("can't get list of ACME windows:", err)
	}

	for _, win := range acmeWindows {
		if dontHandleWindow(win.Name) {
			continue
		}
		go handleWindow(win.ID)
	}

	acmeLog, err := acme.Log()
	if err != nil {
		log.Fatalln("Can't open ACME log")
	}
	defer acmeLog.Close()

	for {
		event, err := acmeLog.Read()
		if err != nil {
			log.Println("can't read event from ACME log:", err)
		}
		if event.Op != "new" {
			continue
		}
		go handleWindow(event.ID)
	}
}
