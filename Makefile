all: dist/glua.js dist/glua.min.js

dist/glua.js: main.go
	gopherjs build -o dist/glua.js

dist/glua.min.js: main.go
	gopherjs build -m -o dist/glua.min.js
