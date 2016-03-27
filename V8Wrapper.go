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
	worker.Load("lodash.js", fmt.Sprintf(`
					    %s
					    var _ = this._;
					    `, string(lodashLib)))
}

//StartV8Worker ...
func StartV8Worker() {
	V8Worker = v8worker.New(func(msg string) {
		output := fmt.Sprintf("Asynchronous Message %s", msg)
		fmt.Println(output)
	}, func(msg string) string {
		output := fmt.Sprintf("got message sync %s", msg)
		fmt.Println(output)
		return msg
	})
	loadLodash(V8Worker)

	V8Worker.Load("transform.js", `$recvSync(function(msg) {
                    var obj = JSON.parse(msg);
                    var output = {
                      name : obj.firstname + " " + obj.lastname,
                      age : obj.age
                    };
					          return JSON.stringify(output);
					        });`)
}

//CleanUpV8Worker ...
func CleanUpV8Worker() {
	fmt.Println("Cleaning UP")
	V8Worker.TerminateExecution()
}
