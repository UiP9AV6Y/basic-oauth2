package utils

type StringReader string

func (r StringReader) Read(p []byte) (n int, err error) {
	for i, _ := range p {
		b := i % len(r)
		p[i] = r[b]
	}

	return len(p), nil
}
