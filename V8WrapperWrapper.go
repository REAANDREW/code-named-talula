package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ry/v8worker"
)

var (
	//V8Worker ...
	V8Worker *v8worker.Worker
)

func loadLodash(worker *v8worker.Worker) {
	lodashLib, err := ioutil.ReadFile("./scripts/lodash.core.min.js")
	if err != nil {
		panic(err)
	}
	worker.Load("lodash.core.min.js", fmt.Sprintf(`
					    %s
					    var _ = this._;
					    `, string(lodashLib)))
}

func loadJsProcessor(worker *v8worker.Worker) {
	processorLib, err := ioutil.ReadFile("./scripts/processor.cnt.js")
	if err != nil {
		panic(err)
	}
	worker.Load("processor.cnt.js", string(processorLib))
}

//DiscardSendSync ...
func DiscardSendSync(msg string) string { return "" }

//StartV8Worker ...
func StartV8Worker() {
	V8Worker = v8worker.New(func(msg string) {}, DiscardSendSync)
	loadLodash(V8Worker)
	loadJsProcessor(V8Worker)
}

//CleanUpV8Worker ...
func CleanUpV8Worker() {
	fmt.Println("Cleaning UP")
	V8Worker.TerminateExecution()
}
