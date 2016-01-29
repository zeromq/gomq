package zmtp

type SecurityNull struct{}

func NewSecurityNull() *SecurityNull {
	return &SecurityNull{}
}

func (s *SecurityNull) Type() SecurityMechanismType {
	return NullSecurityMechanismType
}

func (s *SecurityNull) Handshake() error {
	return nil
}

func (s *SecurityNull) Encrypt(data []byte) []byte {
	return data
}
