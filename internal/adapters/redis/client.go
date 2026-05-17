package redis

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	addr string
}

func NewFromEnv() *Client {
	return &Client{addr: env("REDIS_ADDR", "localhost:6379")}
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.command(ctx, "PING")
	return err
}

func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl > 0 {
		_, err := c.command(ctx, "SETEX", key, strconv.Itoa(int(ttl.Seconds())), value)
		return err
	}
	_, err := c.command(ctx, "SET", key, value)
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	value, err := c.command(ctx, "GET", key)
	if errors.Is(err, ErrNil) {
		return "", ErrNil
	}
	return value, err
}

func (c *Client) Del(ctx context.Context, key string) error {
	_, err := c.command(ctx, "DEL", key)
	return err
}

var ErrNil = errors.New("redis nil")

func (c *Client) command(ctx context.Context, args ...string) (string, error) {
	dialer := net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write([]byte(encode(args...))); err != nil {
		return "", err
	}
	return readRESP(bufio.NewReader(conn))
}

func encode(args ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(args))
	for _, arg := range args {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(arg), arg)
	}
	return b.String()
}

func readRESP(r *bufio.Reader) (string, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return "", err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	switch prefix {
	case '+':
		return line, nil
	case ':':
		return line, nil
	case '-':
		return "", errors.New(line)
	case '$':
		n, err := strconv.Atoi(line)
		if err != nil {
			return "", err
		}
		if n < 0 {
			return "", ErrNil
		}
		buf := make([]byte, n+2)
		if _, err := r.Read(buf); err != nil {
			return "", err
		}
		return string(buf[:n]), nil
	default:
		return "", fmt.Errorf("unknown redis response prefix %q", prefix)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
