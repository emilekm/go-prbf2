package prism

type ServerVersion string

const ServerVersion1 ServerVersion = "1"

type Login1Request struct {
	ServerVersion      ServerVersion
	Username           string
	ClientChallengeKey []byte
}

type Login1Response struct {
	Hash            []byte
	ServerChallenge []byte
}

type Login2Request struct {
	ChallengeDigest string
}
