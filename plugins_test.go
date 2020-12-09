package plugins

import (
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/pkujhd/goloader"
)

var (
	testPlugins *Plugins
	testDir     = "./test_data"
)

func testInit() (p *Plugins, err error) {
	if err = os.Mkdir(testDir, 0744); err != nil {
		return
	}

	return New(testDir)
}

func testTeardown() (err error) {
	return os.RemoveAll(testDir)
}

func TestPlugin_init(t *testing.T) {
	var (
		p   Plugin
		err error
	)

	p.filename = "./foo.o"

	syms := make(symbols)
	if err = goloader.RegSymbol(syms); err != nil {
		err = fmt.Errorf("error registering symbol: %v", err)
		return
	}

	if err = p.init(syms); err != nil {
		t.Fatal(err)
	}

	var sym Symbol
	if sym, err = p.Lookup("github.com/vroomy/plugins/testing/foo.init.0"); err != nil {
		t.Fatal(err)
	}

	runFunc := *(*func())(unsafe.Pointer(&sym))
	runFunc()

}
