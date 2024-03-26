package prism

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"slices"
	"sync"
	"time"
)

const (
	readTimeout  = 5 * time.Minute
	writeTimeout = 5 * time.Minute
)

type Client struct {
	config    ClientConfig
	conn      net.Conn
	msgCh     chan Message
	sendMutex sync.Mutex
}

type ClientConfig struct {
	IP   string
	Port string
	User string
	Pass string
}

func NewClient(config ClientConfig) *Client {
	return &Client{
		config: config,
		msgCh:  make(chan Message),
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", net.JoinHostPort(c.config.IP, c.config.Port))
	if err != nil {
		return err
	}

	c.conn = conn

	c.resetReadDeadline()
	c.resetWriteDeadline()

	c.listen()

	return c.login()
}

func (c *Client) Disconnect() {
	c.conn.Close()
	close(c.msgCh)
}

type SendOpts struct {
	ReplySubjects []Subject
}

type Response struct {
	Messages []Message
}

func (c *Client) Send(msg Message, opts *SendOpts) (*Response, error) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()

	if opts == nil {
		opts = &SendOpts{}
	}

	subjects := append([]Subject{}, opts.ReplySubjects...)

	var resp Response
	var msgErr error

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if len(subjects) == 0 {
				return
			}

			select {
			case <-time.After(5 * time.Second):
				msgErr = errors.New("timeout")
				return
			case m := <-c.msgCh:
				if i := slices.Index(subjects, m.Subject); i != -1 {
					resp.Messages = append(resp.Messages, m)
					subjects = append(subjects[:i], subjects[i+1:]...)
				}
				if slices.Contains(errorSubjects, m.Subject) {
					msgErr = NewErrorFromMessage(msg)
					return
				}
			}
		}
	}()

	_, err := c.conn.Write(msg.Encode())
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	wg.Wait()

	if msgErr != nil {
		return nil, fmt.Errorf("send response: %w", msgErr)
	}

	return &resp, nil
}

func (c *Client) login() error {
	cck := cck(32)

	login1Msg := Message{
		Subject: SubjectLogin1,
		Fields:  [][]byte{[]byte("1"), []byte(c.config.User), cck},
	}

	resp, err := c.Send(login1Msg, &SendOpts{
		ReplySubjects: []Subject{SubjectLogin1},
	})
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	passHash := resp.Messages[0].Fields[0]
	serverChallenge := resp.Messages[0].Fields[1]

	passwordHash := sha1.New()
	saltedPassword := sha1.New()
	challengeDigest := sha1.New()

	_, err = passwordHash.Write([]byte(c.config.Pass))
	if err != nil {
		return err
	}

	salted := append(passHash, SeparatorStart...)
	salted = append(salted, hex.EncodeToString(passwordHash.Sum(nil))...)

	_, err = saltedPassword.Write(salted)
	if err != nil {
		return err
	}

	_, err = challengeDigest.Write(
		bytes.Join(
			[][]byte{
				[]byte(c.config.User),
				cck,
				serverChallenge,
				[]byte(hex.EncodeToString(saltedPassword.Sum(nil))),
			},
			SeparatorField,
		),
	)
	if err != nil {
		return err
	}

	login2Msg := Message{
		Subject: SubjectLogin2,
		Fields:  [][]byte{[]byte(hex.EncodeToString(challengeDigest.Sum(nil)))},
	}

	_, err = c.Send(login2Msg, &SendOpts{
		ReplySubjects: []Subject{SubjectConnected},
	})
	if err != nil {
		return fmt.Errorf("login2: %w", err)
	}

	return nil
}

func (c *Client) listen() {
	scanner := bufio.NewScanner(io.TeeReader(c.conn, os.Stdout))
	scanner.Split(splitMessages)

	go func() {
		for scanner.Scan() {
			c.resetReadDeadline()
			msg, err := DecodeMessage(scanner.Bytes())
			if err != nil {
				fmt.Println("decode error:", err)
				continue
			}

			go func() {
				c.msgCh <- *msg
			}()
		}
	}()
}

func (c *Client) resetReadDeadline() error {
	return c.conn.SetReadDeadline(time.Now().Add(readTimeout))
}

func (c *Client) resetWriteDeadline() error {
	return c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
}

func (c *Client) C() <-chan Message {
	return c.msgCh
}

func splitMessages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	msg, _, found := bytes.Cut(data, SeparatorEnd)
	if !found {
		return 0, nil, nil
	}

	advance = len(msg) + len(SeparatorEnd)

	return advance, msg, nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func cck(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
