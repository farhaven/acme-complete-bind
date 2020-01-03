Command `acme-complete-bind` adds a key binding for ^O to ACME to run
acme-lsp's `L comp -e` command without having to move the mouse from the
text.

# Installation
* Install `acme-lsp` and a few language servers of your choice
* Run `GO111MODULE=on go get github.com/farhaven/acme-complete-bind@latest`.

After the installation, run the command `acme-complete-bind`. It does not take any parameters and only outputs log messages if something goes wrong.

ACME needs to be running before starting `acme-complete-bind`.

# Usage
In a text buffer, enter `^O`, that is, press the `O` key with the control key held. This should run `L comp -e`.

If instead of completions, you get a little picture of some guys head, you should make sure that `acme-complete-bind` is indeed running and that it can connect to your ACME instance via 9p. If that does not help, please open an issue in the bug tracker for this repository.