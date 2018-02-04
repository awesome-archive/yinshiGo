#!/bin/bash
rm -rf dist/yinshiGo
for arch in "386" "amd64"; do
	os=linux
	echo "Start building $os-$arch."
	export GOARCH=$arch
	export GOOS=$os
	go build -o dist/yinshiGo -ldflags="-s -w" github.com/popu125/yinshiGo
	upx dist/yinshiGo
	echo "Packing $os-$arch."
	tar -C dist/ -czf $os-$arch-withdb.tar.gz $(ls dist)
	tar -C dist/ -czf $os-$arch.tar.gz yinshiGo
	rm dist/yinshiGo

	os=windows
	echo "Start building $os-$arch."
	export GOARCH=$arch
	export GOOS=$os
	go build -o dist/yinshiGo.exe -ldflags="-s -w" github.com/popu125/yinshiGo
	upx dist/yinshiGo.exe
	echo "Packing $os-$arch."
	tar -C dist/ -czf $os-$arch-withdb.tar.gz $(ls dist)
	tar -C dist/ -czf $os-$arch.tar.gz yinshiGo.exe
	rm dist/yinshiGo.exe
done
echo "Build Done."