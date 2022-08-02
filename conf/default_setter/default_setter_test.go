package default_setter_test

import (
	"testing"

	"github.com/iotexproject/Bumblebee/conf/default_setter"
	"github.com/iotexproject/Bumblebee/x/ptrx"
	. "github.com/onsi/gomega"
)

func TestStruct(t *testing.T) {
	type A struct {
		A int
		B float32
		C *string
		d string
	}
	dft := A{1, 2, ptrx.String("abc"), "def"}
	tar := A{}
	NewWithT(t).Expect(default_setter.Set(dft, &tar)).To(BeNil())
	NewWithT(t).Expect(dft).To(Equal(tar))
}
