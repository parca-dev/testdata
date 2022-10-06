lint:
	clang-format -i src/*
build:
	gcc src/basic-cpp.cpp -o out/basic-cpp
	gcc src/basic-cpp-plt.cpp -o out/basic-cpp-plt
	gcc src/basic-cpp.cpp -o out/basic-cpp-no-fp -fomit-frame-pointer

validate: build
	../parca-agent/dist/eh-frame --executable out/basic-cpp > tables/ours_basic-cpp.txt
	../parca-agent/dist/eh-frame --executable out/basic-cpp-plt > tables/ours_basic-cpp-plt.txt
	../parca-agent/dist/eh-frame --executable out/basic-cpp-no-fp > tables/ours_basic-cpp-no-fp.txt
	../parca-agent/dist/eh-frame --executable vendored/libc.so.6 > tables/ours_libc_so_6.txt

all: validate
