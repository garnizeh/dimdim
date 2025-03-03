package domain

import (
	"fmt"
	"net/url"
	"strings"
)

type Domain string

func (d Domain) IsDev() bool {
	return strings.HasPrefix(string(d), "localhost")
}

func (d Domain) URL(path ...string) string {
	if strings.TrimSpace(string(d)) == "" {
		// If this happens, it was a programming error because we expect that domain to be
		// programmer provided, not user provided.
		panic("invalid domain")
	}

	var (
		result string
		err    error
	)
	if d.IsDev() {
		result, err = url.JoinPath("http://"+string(d), path...)
	} else {
		result, err = url.JoinPath("https://"+string(d), path...)
	}

	if err != nil {
		// If this happens, it was a programming error because we expect that path to be
		// programmer provided, not user provided.
		panic(fmt.Sprintf("failed to return url from the domain %q: %v", d, err))
	}

	return result
}

func (d Domain) Domain() string {
	return string(d)
}
