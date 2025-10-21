package main

import "github.com/magefile/mage/mg"

type Templ mg.Namespace

func (p Templ) Install() error {
	return GoInstall("github.com/a-h/templ/cmd/templ@latest")
}
