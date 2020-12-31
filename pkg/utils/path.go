package utils

import (
	"path/filepath"
)

func ResolveRelativePath(ref, path string) string {
	base := filepath.Dir(ref)

	return ResolvePath(base, path)
}

func ResolvePath(base, path string) string {
	abs := path

	if !filepath.IsAbs(path) {
		abs = filepath.Join(base, path)
	}

	return filepath.Clean(abs)
}
