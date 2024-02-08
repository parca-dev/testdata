package main

import (
	"path/filepath"
	"runtime"
)

func New() *CI {
	return &CI{
		Base: base().WithDirectory("/src", source()).WithWorkdir("/src"),
	}
}

type CI struct {
	// Base is the base container for the CI contains required dependencies, nix, devbox and the source code.
	Base *Container
}

// Validate runs the build, format, and generate commands and checks if there are any changes in the source code except the out directory.
func (m *CI) Validate() *Container {
	return m.Base.
		WithFocus().
		WithExec([]string{"devbox", "run", "build", "format", "generate"}).
		WithExec([]string{"git", "diff", "--exit-code", ":!out/"})
}

// base returns a ubuntu container with nix and devbox installed.
func base() *Container {
	// base container

	base := dag.Container().From("ubuntu:latest")

	// install base dependencies

	base = base.WithExec([]string{"apt-get", "update"})
	base = base.WithExec([]string{"apt-get", "install", "-y", "curl", "git"})

	// Mount the cache volume to the container

	base = base.WithMountedCache("/nix", dag.CacheVolume("parca-dev-testdata-nix"))
	base = base.WithMountedCache("/etc/nix", dag.CacheVolume("parca-dev-testdata-nix-etc"))

	// install nix

	base = base.WithExec([]string{"sh", "-c", "[ -f /etc/nix/group ] && cp /etc/nix/group /etc/group; exit 0"})
	base = base.WithExec([]string{"sh", "-c", "[ -f /etc/nix/passwd ] && cp /etc/nix/passwd /etc/passwd; exit 0"})
	base = base.WithExec([]string{"sh", "-c", "[ -f /etc/nix/shadow ] && cp /etc/nix/shadow /etc/shadow; exit 0"})
	base = base.WithExec([]string{"sh", "-c", "[ ! -f /nix/receipt.json ] && curl --proto =https --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install linux --extra-conf \"sandbox = false\" --init none --no-confirm; exit 0"})
	base = base.WithExec([]string{"cp", "/etc/group", "/etc/nix/group"})
	base = base.WithExec([]string{"cp", "/etc/passwd", "/etc/nix/passwd"})
	base = base.WithExec([]string{"cp", "/etc/shadow", "/etc/nix/shadow"})
	base = base.WithExec([]string{"sed", "-i", "s/auto-allocate-uids = true/auto-allocate-uids = false/g", "/etc/nix/nix.conf"})
	base = base.WithEnvVariable("PATH", "$PATH:/nix/var/nix/profiles/default/bin", ContainerWithEnvVariableOpts{Expand: true})

	// install devbox

	base = base.WithExec([]string{"adduser", "--disabled-password", "devbox"})
	base = base.WithExec([]string{"addgroup", "devbox", "nixbld"})
	base = base.WithEnvVariable("FORCE", "1")
	base = base.WithExec([]string{"sh", "-c", "curl -fsSL https://get.jetpack.io/devbox | bash"})
	base = base.WithExec([]string{"sh", "-c", `echo 'eval "$(devbox global shellenv --recompute)"' >> ~/.bashrc`})
	base = base.WithExec([]string{"devbox", "version"})

	return base
}

// root returns the root directory of the project as an absolute path.
func root() string {
	// get location of current file
	_, current, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(current), "..")
}

// source returns a directory containing the source code of the project with the required files and directories.
func source() *Directory {
	return dag.Host().Directory(
		root(),
		HostDirectoryOpts{
			Include: []string{
				".git",
				"jitdump/**",
				"src/**",
				"tables/**",
				"vendored/**",
				"devbox.*",
				"go.*",
				"build.go",
			},
		},
	)
}
