package prism

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	user = "superuser"
	pass = "VtVAMW5J"
)

func TestClient(t *testing.T) {
	c := NewClient("10.232.88.130", "4712")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Error("timeout")
				wg.Done()
				return
			case msg := <-c.C():
				wg.Done()
				println(string(msg.Fields[0]))
				t.Error("unexpected message")
				return
			}
		}
	}()

	err := c.Connect(user, pass)
	require.NoError(t, err)

	err = c.Send(Message{
		Subject: SubjectLogin1,
		Fields:  [][]byte{[]byte("1"), []byte(user), []byte(c.cck)},
	})
	require.NoError(t, err)

	wg.Wait()
}
