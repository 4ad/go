package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
)

var Linksparc64 = obj.LinkArch{
	ByteOrder: binary.BigEndian,
	Name:      "sparc64",
	Thechar:   'u',
	Minlc:     4,
	Ptrsize:   8,
	Regsize:   8,
}
