package sshd

import (
	"crypto/rand"
	"crypto/rsa"
	"io"

	"golang.org/x/crypto/ssh"
)

// NewRandomSigner generates a random key of a desired bit length.
func NewRandomSigner(bits int) (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

// NewClientConfig creates a barebones ssh.ClientConfig to be used with ssh.Dial.
func NewClientConfig(name string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: name,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				return
			}),
		},
	}
}

// ConnectShell makes a barebones SSH client session, used for testing.
func ConnectShell(host string, name string, handler func(r io.Reader, w io.WriteCloser)) error {
	config := NewClientConfig(name)
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	in, err := session.StdinPipe()
	if err != nil {
		return err
	}

	out, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	/* FIXME: Do we want to request a PTY?
	err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{})
	if err != nil {
		return err
	}
	*/

	err = session.Shell()
	if err != nil {
		return err
	}

	handler(out, in)

	return nil
}
