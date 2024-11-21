package v0

func AmfNull() []byte {
	return []byte{AMF_NULL}
}

func AmfUndefined() []byte {
	return []byte{AMF_UNDEFINED}
}
