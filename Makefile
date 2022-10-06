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
	../parca-agent/dist/eh-frame --executable vendored/libpython3.10.so.1.0 > tables/ours_libpython3.10.txt
	../parca-agent/dist/eh-frame --executable vendored/systemd > tables/ours_systemd.txt
	../parca-agent/dist/eh-frame --executable vendored/parca-agent > tables/ours_parca-agent.txt
	../parca-agent/dist/eh-frame --executable vendored/ruby > tables/ours_ruby.txt
	../parca-agent/dist/eh-frame --executable vendored/libruby > tables/ours_libruby.txt

all: validate
