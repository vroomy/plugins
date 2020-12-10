package plugins

import (
	"fmt"
	"os"
	"testing"

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

	var (
		sym Symbol
		ok  bool
	)

	if sym, ok = p.Lookup("main.main"); !ok {
		t.Fatal("fn of main.main not found")
	}

	fn := sym.AsEmptyFunc()
	if fn == nil {
		return
	}

	fn()

	if sym, ok = p.Lookup("main.BigInt"); !ok {
		t.Fatal("fn of main.main not found")
	}

	bigIntFn := sym.AsInterfaceFunc()
	if bigIntFn == nil {
		return
	}

	fmt.Println("Value?", bigIntFn())
}
