package gen

import (
	"fmt"
	"strings"
)

// todo: make all private

// Builder is the wrapper over strings.Builder with a couple of syntax sugar methods
type builder struct {
	Builder strings.Builder
}

// Println writes formatted line to string builder
func (b *builder) Printf(format string, a ...interface{}) {
	fmt.Fprintf(&b.Builder, format, a...)
	b.CrLf()
}

// CrLf writes empty line
func (b *builder) CrLf() {
	fmt.Fprintf(&b.Builder, "\n")
}

// Bytes returns the result as the byte array
func (b *builder) Bytes() []byte {
	return []byte(b.Builder.String())
}
