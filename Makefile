EH_FRAME_BIN = ../dist/eh-frame

lint:
	clang-format -i src/*
build:
	gcc src/basic-cpp.cpp -o out/basic-cpp
	gcc src/basic-cpp-plt.cpp -o out/basic-cpp-plt
	gcc src/basic-cpp.cpp -o out/basic-cpp-no-fp -fomit-frame-pointer

validate:
	$(EH_FRAME_BIN) --executable out/basic-cpp > tables/ours_basic-cpp.txt
	$(EH_FRAME_BIN) --executable out/basic-cpp-plt > tables/ours_basic-cpp-plt.txt
	$(EH_FRAME_BIN) --executable out/basic-cpp-no-fp > tables/ours_basic-cpp-no-fp.txt

	$(EH_FRAME_BIN) --executable vendored/libc.so.6 > tables/ours_libc_so_6.txt
	$(EH_FRAME_BIN) --executable vendored/libpython3.10.so.1.0 > tables/ours_libpython3.10.txt
	$(EH_FRAME_BIN) --executable vendored/systemd > tables/ours_systemd.txt
	$(EH_FRAME_BIN) --executable vendored/parca-agent > tables/ours_parca-agent.txt
	$(EH_FRAME_BIN) --executable vendored/ruby > tables/ours_ruby.txt
	$(EH_FRAME_BIN) --executable vendored/libruby > tables/ours_libruby.txt
	$(EH_FRAME_BIN) --executable vendored/redpanda > tables/ours_redpanda.txt

all: validate
