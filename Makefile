EH_FRAME_BIN = ../dist/eh-frame

all: lint build

lint:
	clang-format -i src/*.cpp
	go fmt src/*.go

build:
	gcc src/basic-cpp.cpp -o out/basic-cpp -g
	gcc src/basic-cpp-plt.cpp -o out/basic-cpp-plt
	gcc src/basic-cpp.cpp -o out/basic-cpp-no-fp -fomit-frame-pointer
	gcc src/basic-cpp.cpp -o out/basic-cpp-no-fp-with-debuginfo -fomit-frame-pointer -g
	gcc src/basic-cpp-plt.cpp -o out/basic-cpp-plt-pie -pie -fPIE
	gcc src/basic-cpp-plt.cpp -o out/basic-cpp-plt-hardened -pie -fPIE -fstack-protector-all -D_FORTIFY_SOURCE=2 -Wl,-z,now -Wl,-z,relro -O2
	# The JIT code has frame pointers.
	gcc src/basic-cpp-jit.cpp -o out/basic-cpp-jit -g
	gcc src/basic-cpp-jit.cpp -o out/basic-cpp-jit-no-fp -fomit-frame-pointer -g
	# Go code.
	go build -o out/basic-go src/main.go

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

validate-compact:
	$(EH_FRAME_BIN) --executable out/basic-cpp --compact > compact_tables/ours_basic-cpp.txt
	$(EH_FRAME_BIN) --executable out/basic-cpp-plt --compact > compact_tables/ours_basic-cpp-plt.txt
	$(EH_FRAME_BIN) --executable out/basic-cpp-no-fp --compact > compact_tables/ours_basic-cpp-no-fp.txt

	$(EH_FRAME_BIN) --executable vendored/libc.so.6 --compact > compact_tables/ours_libc_so_6.txt
	$(EH_FRAME_BIN) --executable vendored/libpython3.10.so.1.0 --compact > compact_tables/ours_libpython3.10.txt
	$(EH_FRAME_BIN) --executable vendored/systemd --compact  > compact_tables/ours_systemd.txt
	$(EH_FRAME_BIN) --executable vendored/parca-agent --compact > compact_tables/ours_parca-agent.txt
	$(EH_FRAME_BIN) --executable vendored/ruby --compact  > compact_tables/ours_ruby.txt
	$(EH_FRAME_BIN) --executable vendored/libruby --compact  > compact_tables/ours_libruby.txt
	$(EH_FRAME_BIN) --executable vendored/redpanda --compact > compact_tables/ours_redpanda.txt

all: validate
