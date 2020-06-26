build: clean
	GOOS=js GOARCH=wasm go build -o main.wasm
zip: clean build
	gzip main.wasm
	mv main.wasm.gz main.wasm
tiny: clean
	tinygo build -o wasm.wasm -target wasm ./main.go
node:
	GOOS=js GOARCH=wasm go run -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec" .
wasm_js:
	cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
clean:
	-rm main.wasm
