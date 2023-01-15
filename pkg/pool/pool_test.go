package pool_test

import (
	"context"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	// Packages
	pool "github.com/mutablelogic/go-accessory/pkg/pool"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

const (
	MONGO_URL = "${MONGO_URL}"
)

func Test_Pool_001(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t))
	assert.NotNil(pool)
	assert.NoError(pool.Close())
}

func Test_Pool_002(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t))
	assert.NotNil(pool)
	defer pool.Close()

	conn := pool.Get()
	assert.NotNil(conn)
	t.Log(conn)
}

func Test_Pool_003(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t), pool.OptTrace(func(ctx context.Context, delta time.Duration, err error) {
		if err != nil {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", err)
		} else {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
		}
	}))
	assert.NotNil(pool)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			conn := pool.Get()
			assert.NotNil(conn)
			assert.NoError(conn.Ping(context.TODO()))
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			pool.Put(conn)
		}()
	}
	wg.Wait()
	assert.NoError(pool.Close())
	assert.Equal(0, pool.Size())
}

func Test_Pool_004(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t), pool.OptTrace(func(ctx context.Context, delta time.Duration, err error) {
		if err != nil {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", err)
		} else {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
		}
	}))
	assert.NotNil(pool)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			time.Sleep(time.Millisecond * 10)
			t.Log("Spawning", pool.Size(), "connections")
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				conn := pool.Get()
				assert.NotNil(conn)
				assert.NoError(conn.Ping(context.TODO()))
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				pool.Put(conn)
			}()
		}
	}
	wg.Wait()
	assert.NoError(pool.Close())
	assert.Equal(0, pool.Size())
}

func Test_Pool_005(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t), pool.OptTrace(func(ctx context.Context, delta time.Duration, err error) {
		if err != nil {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", err)
		} else {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
		}
	}), pool.OptMaxSize(10))
	assert.NotNil(pool)

	var conns []Conn
	for i := 0; i < 20; i++ {
		conn := pool.Get()
		if i < 10 {
			assert.NotNil(conn)
			conns = append(conns, conn)
		} else {
			assert.Nil(conn)
		}
	}
	for _, conn := range conns {
		pool.Put(conn)
	}
	assert.NoError(pool.Close())
	assert.Equal(0, pool.Size())
}

///////////////////////////////////////////////////////////////////////////////
// Return URL or skip test

func uri(t *testing.T) *url.URL {
	if uri := os.ExpandEnv(MONGO_URL); uri == "" {
		t.Skip("no MONGO_URL environment variable, skipping test")
	} else if uri, err := url.Parse(uri); err != nil {
		t.Skip("invalid MONGO_URL environment variable, skipping test")
	} else {
		return uri
	}
	return nil
}
