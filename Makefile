EH_FRAME_BIN = ../dist/eh-frame

all: lint build

lint:
	clang-format -i src/*.cpp
	go fmt src/*.go

build:
	gcc src/basic-cpp.cpp -o out/x86/basic-cpp -g
	gcc src/basic-cpp-plt.cpp -o out/x86/basic-cpp-plt
	gcc src/basic-cpp.cpp -o out/x86/basic-cpp-no-fp -fomit-frame-pointer
	gcc src/basic-cpp.cpp -o out/x86/basic-cpp-no-fp-with-debuginfo -fomit-frame-pointer -g
	gcc src/basic-cpp-plt.cpp -o out/x86/basic-cpp-plt-pie -pie -fPIE
	gcc src/basic-cpp-plt.cpp -o out/x86/basic-cpp-plt-hardened -pie -fPIE -fstack-protector-all -D_FORTIFY_SOURCE=2 -Wl,-z,now -Wl,-z,relro -O2
	gcc src/basic-cpp-multithreaded-no-fp.cpp -o out/x86/basic-cpp-multithreaded-no-fp -fomit-frame-pointer
	# The JIT code has frame pointers.
	gcc src/basic-cpp-jit.cpp -o out/x86/basic-cpp-jit -g
	gcc src/basic-cpp-jit.cpp -o out/x86/basic-cpp-jit-no-fp -fomit-frame-pointer -g
	# Go code.
	go build -o out/x86/basic-go src/main.go

	gcc src/basic-cpp.cpp -fno-PIE -no-pie -o out/arm64/basic-cpp -g
	gcc src/basic-cpp-plt.cpp -fno-PIE -no-pie -o out/arm64/basic-cpp-plt
	gcc src/basic-cpp.cpp -fno-PIE -no-pie -o out/arm64/basic-cpp-no-fp -fomit-frame-pointer
	gcc src/basic-cpp.cpp -fno-PIE -no-pie -o out/arm64/basic-cpp-no-fp-with-debuginfo -fomit-frame-pointer -g
	gcc src/basic-cpp-plt.cpp -o out/arm64/basic-cpp-plt-pie -pie -fPIE
	gcc src/basic-cpp-plt.cpp -o out/arm64/basic-cpp-plt-hardened -pie -fPIE -fstack-protector-all -D_FORTIFY_SOURCE=2 -Wl,-z,now -Wl,-z,relro -O2

	# The JIT code has frame pointers but isn't for Arm64
	# TODO(sylfrena): Port JIT toy binaries for arm64
	# gcc src/basic-cpp-jit.cpp -o out/arm64/basic-cpp-jit -g
	# gcc src/basic-cpp-jit.cpp -o out/arm64/basic-cpp-jit-no-fp -fomit-frame-pointer -g

	# Go code.
	go build -o out/arm64/basic-go src/main.go

validate:
	$(EH_FRAME_BIN) --executable out/x86/basic-cpp > tables/x86/ours_basic-cpp.txt
	$(EH_FRAME_BIN) --executable out/x86/basic-cpp-plt > tables/x86/ours_basic-cpp-plt.txt
	$(EH_FRAME_BIN) --executable out/x86/basic-cpp-no-fp > tables/x86/ours_basic-cpp-no-fp.txt

	$(EH_FRAME_BIN) --executable vendored/x86/libc.so.6 > tables/x86/ours_libc_so_6.txt
	$(EH_FRAME_BIN) --executable vendored/x86/libpython3.10.so.1.0 > tables/x86/ours_libpython3.10.txt
	$(EH_FRAME_BIN) --executable vendored/x86/systemd > tables/x86/ours_systemd.txt
	$(EH_FRAME_BIN) --executable vendored/x86/parca-agent > tables/x86/ours_parca-agent.txt
	$(EH_FRAME_BIN) --executable vendored/x86/ruby > tables/x86/ours_ruby.txt
	$(EH_FRAME_BIN) --executable vendored/x86/libruby > tables/x86/ours_libruby.txt
	$(EH_FRAME_BIN) --executable vendored/x86/redpanda > tables/x86/ours_redpanda.txt

	$(EH_FRAME_BIN) --executable out/arm64/basic-cpp > tables/arm64/ours_basic-cpp.txt
	$(EH_FRAME_BIN) --executable out/arm64/basic-cpp-plt > tables/arm64/ours_basic-cpp-plt.txt
	$(EH_FRAME_BIN) --executable out/arm64/basic-cpp-no-fp > tables/arm64/ours_basic-cpp-no-fp.txt

	# $(EH_FRAME_BIN) --executable vendored/arm64/libc.so.6 > tables/arm64/ours_libc_so_6.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/libpython3.10.so.1.0 > tables/arm64/ours_libpython3.10.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/systemd > tables/arm64/ours_systemd.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/parca-agent > tables/arm64/ours_parca-agent.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/ruby > tables/arm64/ours_ruby.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/libruby > tables/arm64/ours_libruby.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/redpanda > tables/arm64/ours_redpanda.txt


validate-compact:
	$(EH_FRAME_BIN) --executable out/x86/basic-cpp --compact > compact_tables/x86/ours_basic-cpp.txt
	$(EH_FRAME_BIN) --executable out/x86/basic-cpp-plt --compact > compact_tables/x86/ours_basic-cpp-plt.txt
	$(EH_FRAME_BIN) --executable out/x86/basic-cpp-no-fp --compact > compact_tables/x86/ours_basic-cpp-no-fp.txt

	$(EH_FRAME_BIN) --executable vendored/x86/libc.so.6 --compact > compact_tables/x86/ours_libc_so_6.txt
	$(EH_FRAME_BIN) --executable vendored/x86/libpython3.10.so.1.0 --compact > compact_tables/x86/ours_libpython3.10.txt
	$(EH_FRAME_BIN) --executable vendored/x86/systemd --compact  > compact_tables/x86/ours_systemd.txt
	$(EH_FRAME_BIN) --executable vendored/x86/parca-agent --compact > compact_tables/x86/ours_parca-agent.txt
	$(EH_FRAME_BIN) --executable vendored/x86/ruby --compact  > compact_tables/x86/ours_ruby.txt
	$(EH_FRAME_BIN) --executable vendored/x86/libruby --compact  > compact_tables/x86/ours_libruby.txt
	$(EH_FRAME_BIN) --executable vendored/x86/redpanda --compact > compact_tables/x86/ours_redpanda.txt

	$(EH_FRAME_BIN) --executable out/arm64/basic-cpp --compact > compact_tables/arm64/ours_basic-cpp.txt
	$(EH_FRAME_BIN) --executable out/arm64/basic-cpp-plt --compact > compact_tables/arm64/ours_basic-cpp-plt.txt
	$(EH_FRAME_BIN) --executable out/arm64/basic-cpp-no-fp --compact > compact_tables/arm64/ours_basic-cpp-no-fp.txt
	$(EH_FRAME_BIN) --executable vendored/arm64/libc.so.6 --compact > compact_tables/arm64/ours_libc_so_6.txt

	# $(EH_FRAME_BIN) --executable vendored/arm64/libpython3.10.so.1.0 --compact > compact_tables/arm64/ours_libpython3.10.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/systemd --compact  > compact_tables/arm64/ours_systemd.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/parca-agent --compact > compact_tables/arm64/ours_parca-agent.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/ruby --compact  > compact_tables/arm64/ours_ruby.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/libruby --compact  > compact_tables/arm64/ours_libruby.txt
	# $(EH_FRAME_BIN) --executable vendored/arm64/redpanda --compact > compact_tables/arm64/ours_redpanda.txt

all: validate
