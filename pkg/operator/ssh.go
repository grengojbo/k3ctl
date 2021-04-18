package operator

import (
	"github.com/appleboy/easyssh-proxy"
	// log "github.com/sirupsen/logrus"
)

// type SSHConfig struct {
// 	User    string
// 	Server  string
// 	Port    string
// 	KeyPath string
// }
type SSHOperator struct {
	Config *easyssh.MakeConfig
	// ProxyConfig SSHConfig
}

func MakeSsshConfig(config *easyssh.MakeConfig) {

	// ssh := &easyssh.MakeConfig{
	// 	User:   "appleboy",
	// 	Server: "example.com",
	// 	// Optional key or Password without either we try to contact your agent SOCKET
	// 	// Password: "password",
	// 	// Paste your source content of private key
	// 	// Key: `-----BEGIN RSA PRIVATE KEY-----
	// 	// MIIEpAIBAAKCAQEA4e2D/qPN08pzTac+a8ZmlP1ziJOXk45CynMPtva0rtK/RB26
	// 	// 7XC9wlRna4b3Ln8ew3q1ZcBjXwD4ppbTlmwAfQIaZTGJUgQbdsO9YA==
	// 	// -----END RSA PRIVATE KEY-----
	// 	// `,
	// 	KeyPath: "/Users/username/.ssh/id_rsa",
	// 	Port:    "22",
	// 	Timeout: 60 * time.Second,

	// 	// Parse PrivateKey With Passphrase
	// 	Passphrase: "1234",

	// 	// Optional fingerprint SHA256 verification
	// 	// Get Fingerprint: ssh.FingerprintSHA256(key)
	// 	// Fingerprint: "SHA256:mVPwvezndPv/ARoIadVY98vAC0g+P/5633yTC4d/wXE"

	// 	// Enable the use of insecure ciphers and key exchange methods.
	// 	// This enables the use of the the following insecure ciphers and key exchange methods:
	// 	// - aes128-cbc
	// 	// - aes192-cbc
	// 	// - aes256-cbc
	// 	// - 3des-cbc
	// 	// - diffie-hellman-group-exchange-sha256
	// 	// - diffie-hellman-group-exchange-sha1
	// 	// Those algorithms are insecure and may allow plaintext data to be recovered by an attacker.
	// 	// UseInsecureCipher: true,
	// }
	// log.Debugln("ssh: ", ssh)
	// return nil, nil
}
