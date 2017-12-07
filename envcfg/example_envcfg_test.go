package envcfg_test

import (
	"fmt"

	"github.com/vivareal/ligeiro/envcfg"
)

func ExampleLoad() {
	Config := envcfg.Load(envcfg.Map{
		"APPLICATION": "myappname",
		"LISTEN_PORT": ":8080",
	})

	fmt.Println(Config.Get("listenPort"))
	fmt.Println(Config.Get("logLevel"))
	// Output:
	// :8080
	// debug
}

func ExampleLoadBundled() {
	Config := envcfg.LoadBundled()

	fmt.Println(Config.Get("version"))
	fmt.Println(Config.Get("logLevel"))
	// Output:
	// detached
	// debug
}
