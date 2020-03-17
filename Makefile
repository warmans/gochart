.PHONY: demo.build
demo.build:
	cd demo && GOOS=js GOARCH=wasm go build -o main.wasm

.PHONY: demo.serve
demo.serve:
	cd demo && goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'