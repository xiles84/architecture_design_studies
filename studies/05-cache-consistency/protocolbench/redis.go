package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type redisConn struct {
	c net.Conn
	r *bufio.Reader
}
type redisStore struct {
	address string
	idle    chan *redisConn
}

func newRedis(address string) *redisStore {
	return &redisStore{address: address, idle: make(chan *redisConn, 24)}
}
func (r *redisStore) close() {
	for {
		select {
		case c := <-r.idle:
			c.c.Close()
		default:
			return
		}
	}
}
func readRESP(r *bufio.Reader) (any, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 3 {
		return nil, fmt.Errorf("short RESP")
	}
	body := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	switch body[0] {
	case '+':
		return body[1:], nil
	case '-':
		return nil, fmt.Errorf("redis: %s", body[1:])
	case ':':
		return strconv.ParseInt(body[1:], 10, 64)
	case '$':
		n, e := strconv.Atoi(body[1:])
		if e != nil {
			return nil, e
		}
		if n < 0 {
			return nil, nil
		}
		b := make([]byte, n+2)
		_, e = io.ReadFull(r, b)
		if e != nil {
			return nil, e
		}
		return string(b[:n]), nil
	}
	return nil, fmt.Errorf("unsupported RESP %q", body)
}
func (r *redisStore) do(ctx context.Context, args ...string) (any, error) {
	var c *redisConn
	select {
	case c = <-r.idle:
	default:
		raw, e := (&net.Dialer{}).DialContext(ctx, "tcp", r.address)
		if e != nil {
			return nil, e
		}
		c = &redisConn{raw, bufio.NewReader(raw)}
	}
	deadline := time.Now().Add(10 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	_ = c.c.SetDeadline(deadline)
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(args))
	for _, a := range args {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(a), a)
	}
	_, err := io.WriteString(c.c, b.String())
	var v any
	if err == nil {
		v, err = readRESP(c.r)
	}
	if err != nil {
		c.c.Close()
		return nil, err
	}
	select {
	case r.idle <- c:
	default:
		c.c.Close()
	}
	return v, nil
}
func key(k int64) string      { return fmt.Sprintf("cp:%d", k) }
func fenceKey(k int64) string { return key(k) + ":f" }
func encode(s state) string   { return fmt.Sprintf("%d|%d|%t", s.Generation, s.Value, s.Deleted) }
func decode(v string) (state, error) {
	p := strings.Split(v, "|")
	if len(p) != 3 {
		return state{}, fmt.Errorf("bad cache payload")
	}
	g, e := strconv.ParseInt(p[0], 10, 64)
	if e != nil {
		return state{}, e
	}
	x, e := strconv.ParseInt(p[1], 10, 64)
	if e != nil {
		return state{}, e
	}
	d, e := strconv.ParseBool(p[2])
	return state{g, x, d}, e
}
func (r *redisStore) get(ctx context.Context, k int64) (state, bool, error) {
	v, e := r.do(ctx, "GET", key(k))
	if e != nil || v == nil {
		return state{}, false, e
	}
	s, e := decode(v.(string))
	return s, true, e
}
func (r *redisStore) fence(ctx context.Context, k int64) (int64, error) {
	v, e := r.do(ctx, "GET", fenceKey(k))
	if e != nil {
		return 0, e
	}
	if v == nil {
		return 0, nil
	}
	return strconv.ParseInt(v.(string), 10, 64)
}
func (r *redisStore) invalidate(ctx context.Context, k int64) error {
	_, e := r.do(ctx, "EVAL", "local n=redis.call('INCR',KEYS[2]);redis.call('DEL',KEYS[1]);return n", "2", key(k), fenceKey(k))
	return e
}

const publishLua = `if tonumber(redis.call('GET',KEYS[2]) or '0')~=tonumber(ARGV[1]) then return 0 end
local old=redis.call('GET',KEYS[1]);if old and tonumber(string.match(old,'^(%d+)|'))>tonumber(ARGV[2]) then return 0 end
redis.call('SET',KEYS[1],ARGV[3],'PX',ARGV[4]);return 1`

func (r *redisStore) put(ctx context.Context, k int64, s state, f int64, ttl time.Duration) (bool, error) {
	v, e := r.do(ctx, "EVAL", publishLua, "2", key(k), fenceKey(k), fmt.Sprint(f), fmt.Sprint(s.Generation), encode(s), fmt.Sprint(ttl.Milliseconds()))
	if e != nil {
		return false, e
	}
	return v.(int64) == 1, nil
}
