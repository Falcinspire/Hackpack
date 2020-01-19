package service

// import (
// 	"io/ioutil"
// 	"os"

// 	"github.com/falcinspire/hackpack/internal/app"
// )

// func Service() {
// 	dir, err := ioutil.TempDir("app", "project")
// 	if err != nil {
// 		panic(err)
// 	}
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}
// 	os.Chdir(dir)
// 	err = unzip("hackpack.zip", ".")
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = os.Remove("hackpack.zip")
// 	if err != nil {
// 		panic(err)
// 	}
// 	app.CompileHackpack()
// 	os.Chdir(wd)
// }
