package main

type CI struct{}

func (m *CI) devbox(src *Directory) *Container {
	// TODO(kakkoyun): Use local container.
	return dag.Container().
		From("jetpackio/devbox-root-user:latest").
		WithMountedDirectory("/code", src).
		WithWorkdir("/code").
		WithExec([]string{"devbox", "run", "--", "echo", "'Installed Packages.'"})
}

func (m *CI) build(src *Directory) *Container {
	return m.devbox(src).
		WithExec([]string{"devbox", "run", "build"})
}

func (m *CI) generate(src *Directory) *Container {
	return m.build(src).
		WithExec([]string{"devbox", "run", "generate"})
}

func (m *CI) validate(src *Directory) *Container {
	return m.generate(src).
		WithExec([]string{"git", "diff", "--exit-code", ":!out/"})
}

func (m *CI) Run(src *Directory) *Container {
	return m.validate(src)
}
