package kv_test

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
	"github.com/hanzokv/go/v9"
	"github.com/hanzokv/go/v9/logging"
)

const (
	kvSecondaryPort = "6381"
)

const (
	ringShard1Port = "6390"
	ringShard2Port = "6391"
	ringShard3Port = "6392"
)

const (
	sentinelName       = "go-kv-test"
	sentinelMasterPort = "9121"
	sentinelSlave1Port = "9122"
	sentinelSlave2Port = "9123"
	sentinelPort1      = "26379"
	sentinelPort2      = "26380"
	sentinelPort3      = "26381"
)

var (
	kvPort = "6380"
	kvAddr = ":" + kvPort
)

var (
	kvStackPort = "6379"
	kvStackAddr = ":" + kvStackPort
)

var (
	sentinelAddrs = []string{":" + sentinelPort1, ":" + sentinelPort2, ":" + sentinelPort3}

	ringShard1, ringShard2, ringShard3             *kv.Client
	sentinelMaster, sentinelSlave1, sentinelSlave2 *kv.Client
	sentinel1, sentinel2, sentinel3                *kv.Client
)

var cluster = &clusterScenario{
	ports:   []string{"16600", "16601", "16602", "16603", "16604", "16605"},
	nodeIDs: make([]string, 6),
	clients: make(map[string]*kv.Client, 6),
}

// KV Software Cluster
var RECluster = false

// KV Community Edition Docker
var RCEDocker = false

// Notes version of kv we are executing tests against.
// This can be used before we change the bsm fork of ginkgo for one,
// which have support for label sets, so we can filter tests per kv version.
var KVVersion float64 = 8.4

func SkipBeforeKVVersion(version float64, msg string) {
	if KVVersion < version {
		Skip(fmt.Sprintf("(kv version < %f) %s", version, msg))
	}
}

func SkipAfterKVVersion(version float64, msg string) {
	if KVVersion > version {
		Skip(fmt.Sprintf("(kv version > %f) %s", version, msg))
	}
}

var _ = BeforeSuite(func() {
	addr := os.Getenv("REDIS_PORT")
	if addr != "" {
		kvPort = addr
		kvAddr = ":" + kvPort
	}
	var err error
	RECluster, _ = strconv.ParseBool(os.Getenv("RE_CLUSTER"))
	RCEDocker, _ = strconv.ParseBool(os.Getenv("RCE_DOCKER"))

	KVVersion, _ = strconv.ParseFloat(strings.Trim(os.Getenv("KV_VERSION"), "\""), 64)

	if KVVersion == 0 {
		KVVersion = 8.4
	}

	fmt.Printf("RECluster: %v\n", RECluster)
	fmt.Printf("RCEDocker: %v\n", RCEDocker)
	fmt.Printf("KV_VERSION: %.1f\n", KVVersion)
	fmt.Printf("CLIENT_LIBS_TEST_IMAGE: %v\n", os.Getenv("CLIENT_LIBS_TEST_IMAGE"))
	logging.Disable()

	if KVVersion < 7.0 || KVVersion > 9 {
		panic("incorrect or not supported kv version")
	}

	kvPort = kvStackPort
	kvAddr = kvStackAddr
	if !RECluster {
		ringShard1, err = connectTo(ringShard1Port)
		Expect(err).NotTo(HaveOccurred())

		ringShard2, err = connectTo(ringShard2Port)
		Expect(err).NotTo(HaveOccurred())

		ringShard3, err = connectTo(ringShard3Port)
		Expect(err).NotTo(HaveOccurred())

		sentinelMaster, err = connectTo(sentinelMasterPort)
		Expect(err).NotTo(HaveOccurred())

		sentinel1, err = startSentinel(sentinelPort1, sentinelName, sentinelMasterPort)
		Expect(err).NotTo(HaveOccurred())

		sentinel2, err = startSentinel(sentinelPort2, sentinelName, sentinelMasterPort)
		Expect(err).NotTo(HaveOccurred())

		sentinel3, err = startSentinel(sentinelPort3, sentinelName, sentinelMasterPort)
		Expect(err).NotTo(HaveOccurred())

		sentinelSlave1, err = connectTo(sentinelSlave1Port)
		Expect(err).NotTo(HaveOccurred())

		err = sentinelSlave1.SlaveOf(ctx, "127.0.0.1", sentinelMasterPort).Err()
		Expect(err).NotTo(HaveOccurred())

		sentinelSlave2, err = connectTo(sentinelSlave2Port)
		Expect(err).NotTo(HaveOccurred())

		err = sentinelSlave2.SlaveOf(ctx, "127.0.0.1", sentinelMasterPort).Err()
		Expect(err).NotTo(HaveOccurred())

		// populate cluster node information
		Expect(configureClusterTopology(ctx, cluster)).NotTo(HaveOccurred())
	}
})

var _ = AfterSuite(func() {
	if !RECluster {
		Expect(cluster.Close()).NotTo(HaveOccurred())
	}
})

func TestGinkgoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-kv")
}

//------------------------------------------------------------------------------

func kvOptions() *kv.Options {
	if RECluster {
		return &kv.Options{
			Addr: kvAddr,
			DB:   0,

			DialTimeout:           10 * time.Second,
			ReadTimeout:           30 * time.Second,
			WriteTimeout:          30 * time.Second,
			ContextTimeoutEnabled: true,

			MaxRetries: -1,

			PoolSize:        10,
			PoolTimeout:     30 * time.Second,
			ConnMaxIdleTime: time.Minute,
		}
	}
	return &kv.Options{
		Addr: kvAddr,
		DB:   0,

		DialTimeout:           10 * time.Second,
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		ContextTimeoutEnabled: true,

		MaxRetries: -1,

		PoolSize:        10,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: time.Minute,
	}
}

func kvClusterOptions() *kv.ClusterOptions {
	return &kv.ClusterOptions{
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,

		MaxRedirects: 8,

		PoolSize:        10,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: time.Minute,
	}
}

func kvRingOptions() *kv.RingOptions {
	return &kv.RingOptions{
		Addrs: map[string]string{
			"ringShardOne": ":" + ringShard1Port,
			"ringShardTwo": ":" + ringShard2Port,
		},

		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,

		MaxRetries: -1,

		PoolSize:        10,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: time.Minute,
	}
}

func performAsync(n int, cbs ...func(int)) *sync.WaitGroup {
	var wg sync.WaitGroup
	for _, cb := range cbs {
		wg.Add(n)
		// start from 1, so we can skip db 0 where such test is executed with
		// select db command
		for i := 1; i <= n; i++ {
			go func(cb func(int), i int) {
				defer GinkgoRecover()
				defer wg.Done()

				cb(i)
			}(cb, i)
		}
	}
	return &wg
}

func perform(n int, cbs ...func(int)) {
	wg := performAsync(n, cbs...)
	wg.Wait()
}

func eventually(fn func() error, timeout time.Duration) error {
	errCh := make(chan error, 1)
	done := make(chan struct{})
	exit := make(chan struct{})

	go func() {
		for {
			err := fn()
			if err == nil {
				close(done)
				return
			}

			select {
			case errCh <- err:
			default:
			}

			select {
			case <-exit:
				return
			case <-time.After(timeout / 100):
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		close(exit)
		select {
		case err := <-errCh:
			return err
		default:
			return fmt.Errorf("timeout after %s without an error", timeout)
		}
	}
}

func connectTo(port string) (*kv.Client, error) {
	client := kv.NewClient(&kv.Options{
		Addr:       ":" + port,
		MaxRetries: -1,
	})

	err := eventually(func() error {
		return client.Ping(ctx).Err()
	}, 30*time.Second)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func startSentinel(port, masterName, masterPort string) (*kv.Client, error) {
	client, err := connectTo(port)
	if err != nil {
		return nil, err
	}

	for _, cmd := range []*kv.StatusCmd{
		kv.NewStatusCmd(ctx, "SENTINEL", "MONITOR", masterName, "127.0.0.1", masterPort, "2"),
	} {
		client.Process(ctx, cmd)
		if err := cmd.Err(); err != nil && !strings.Contains(err.Error(), "ERR Duplicate master name.") {
			return nil, fmt.Errorf("%s failed: %w", cmd, err)
		}
	}

	return client, nil
}

//------------------------------------------------------------------------------

type badConnError string

func (e badConnError) Error() string   { return string(e) }
func (e badConnError) Timeout() bool   { return true }
func (e badConnError) Temporary() bool { return false }

type badConn struct {
	net.TCPConn

	readDelay, writeDelay time.Duration
	readErr, writeErr     error
}

var _ net.Conn = &badConn{}

func (cn *badConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (cn *badConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (cn *badConn) Read([]byte) (int, error) {
	if cn.readDelay != 0 {
		time.Sleep(cn.readDelay)
	}
	if cn.readErr != nil {
		return 0, cn.readErr
	}
	return 0, badConnError("bad connection")
}

func (cn *badConn) Write([]byte) (int, error) {
	if cn.writeDelay != 0 {
		time.Sleep(cn.writeDelay)
	}
	if cn.writeErr != nil {
		return 0, cn.writeErr
	}
	return 0, badConnError("bad connection")
}

//------------------------------------------------------------------------------

type hook struct {
	dialHook            func(hook kv.DialHook) kv.DialHook
	processHook         func(hook kv.ProcessHook) kv.ProcessHook
	processPipelineHook func(hook kv.ProcessPipelineHook) kv.ProcessPipelineHook
}

func (h *hook) DialHook(hook kv.DialHook) kv.DialHook {
	if h.dialHook != nil {
		return h.dialHook(hook)
	}
	return hook
}

func (h *hook) ProcessHook(hook kv.ProcessHook) kv.ProcessHook {
	if h.processHook != nil {
		return h.processHook(hook)
	}
	return hook
}

func (h *hook) ProcessPipelineHook(hook kv.ProcessPipelineHook) kv.ProcessPipelineHook {
	if h.processPipelineHook != nil {
		return h.processPipelineHook(hook)
	}
	return hook
}
