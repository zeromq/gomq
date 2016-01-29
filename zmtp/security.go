package zmtp

type SecurityMechanismType string

const (
	NullSecurityMechanismType  SecurityMechanismType = "NULL"
	PlainSecurityMechanismType SecurityMechanismType = "PLAIN"
	CurveSecurityMechanismType SecurityMechanismType = "CURVE"
)

type SecurityMechanism interface {
	Type() SecurityMechanismType
	Handshake() error
	Encrypt([]byte) []byte
}
