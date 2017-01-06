package lirc

/*
#cgo LDFLAGS: -llirc_client
#include <stdlib.h>
#include <lirc/lirc_client.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Callback type is function type that takes as input the code of the button that was pressed.
type Callback func(string)

var config *C.struct_lirc_config
var running bool = false

// Start the lirc server for the provided program name and config file path.
// prg specifies the name of the program registered with lirc,
// cfg is the path to the config file for the program,
// cb defines the callback to be run on events.
func Start(prg, cfg string, cb Callback) error {
	prg_name := C.CString(prg)
	cfg_path := C.CString(cfg)

	if ret := C.lirc_init(prg_name, 1); ret == -1 {
		return errors.New("could not init program with lirc")
	}

	if ret := C.lirc_readconfig(cfg_path, &config, nil); ret != 0 {
		return errors.New("could not read configuration")
	}

	C.free(unsafe.Pointer(prg_name))
	C.free(unsafe.Pointer(cfg_path))
	running = true
	go run(cb)
	return nil
}

// Stop listening for input.
func Stop() {
	running = false
}

func run(cb Callback) {
	var code *C.char // string code of input
	var c *C.char
	for running && C.lirc_nextcode(&code) == 0 {
		if code == nil {
			continue
		}
		for C.lirc_code2char(config, code, &c) == 0 {
			if c == nil {
				break
			}
			cb(C.GoString(c))
		}
		C.free(unsafe.Pointer(code))
	}
}
