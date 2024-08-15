package token

type TokenGenerator interface {
	Generate(payload interface{}) (string, error)
	Verify(token string) (interface{}, error)
	VerifyAndScan(token string, dest interface{}) error
}
