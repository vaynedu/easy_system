package third_part

import (
	"errors"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// HTTPClient 全局共享的HTTP客户端单例
	HTTPClient *resty.Client

	// once 确保HTTPClient只初始化一次
	once sync.Once
)

// init 包初始化函数，确保HTTPClient在包加载时初始化
func init() {
	InitHTTPClient()
}

// InitHTTPClient 初始化全局HTTP客户端
func InitHTTPClient() {
	once.Do(func() {
		// 创建HTTP传输层配置
		transport := NewTransport()

		// 创建HTTP客户端
		HTTPClient = resty.New().
			SetTransport(transport).
			SetTimeout(10 * time.Second).        // 设置全局超时时间
			SetRetryCount(3).                    // 设置重试次数
			SetRetryWaitTime(1 * time.Second).   // 设置重试等待时间
			SetRetryMaxWaitTime(5 * time.Second) // 设置最大重试等待时间
	})
}

// GetHTTPClient 获取全局HTTP客户端实例
func GetHTTPClient() *resty.Client {
	InitHTTPClient() // 确保客户端已初始化
	return HTTPClient
}

// NewRestyClient 创建新的HTTP客户端实例（用于特殊需求场景）
func NewRestyClient() *resty.Client {
	return resty.NewWithClient(newClient())
}

func newClient() *http.Client {
	return &http.Client{
		Transport: NewTransport(),
		Timeout:   10 * time.Second,
	}
}

// NewTransport 创建HTTP传输层配置
func NewTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// DialContext 决定了HTTP客户端如何拨号并建立底层的TCP连接
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,  // TCP连接超时时间
			KeepAlive: 30 * time.Second, // TCP连接保持时间
			DualStack: true,             // 是否启用IPv4和IPv6双栈支持
		}).DialContext,
		ForceAttemptHTTP2:     false,                      // 不强制尝试HTTP/2
		MaxConnsPerHost:       runtime.GOMAXPROCS(0) * 64, // 每个主机的最大连接数
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) * 64, // 每个主机的最大空闲连接数
		IdleConnTimeout:       30 * time.Second,           // 空闲连接超时时间
		TLSHandshakeTimeout:   3 * time.Second,            // TLS握手超时时间
		ExpectContinueTimeout: 1 * time.Second,            // 100-continue握手超时时间
		DisableKeepAlives:     false,                      // 启用长连接
	}
}

// RequestTrace 定义请求跟踪信息的结构化数据
type RequestTrace struct {
	DNSLookup      time.Duration // DNS查询时间
	ConnTime       time.Duration // 连接建立时间(含TCP+TLS)
	TCPConnTime    time.Duration // TCP连接时间
	TLSHandshake   time.Duration // TLS握手时间
	ServerTime     time.Duration // 服务器处理时间
	ResponseTime   time.Duration // 响应总时间
	TotalTime      time.Duration // 请求总耗时
	ResponseSize   int           // 响应体大小
	IsConnReused   bool          // 连接是否复用
	IsConnWasIdle  bool          // 连接是否为空闲连接
	ConnIdleTime   time.Duration // 连接空闲时间
	RequestAttempt int           // 请求尝试次数
	RemoteAddr     string        // 远程服务器地址
}

// GetTraceInfo 从响应中提取跟踪信息并返回结构化数据
func GetTraceInfo(resp *resty.Response) (*RequestTrace, error) {
	if resp == nil || resp.Request == nil {
		return nil, errors.New("resp or request is nil")
	}

	traceInfo := resp.Request.TraceInfo()
	return &RequestTrace{
		DNSLookup:      traceInfo.DNSLookup,
		ConnTime:       traceInfo.ConnTime,
		TCPConnTime:    traceInfo.TCPConnTime,
		TLSHandshake:   traceInfo.TLSHandshake,
		ServerTime:     traceInfo.ServerTime,
		ResponseTime:   traceInfo.ResponseTime,
		TotalTime:      traceInfo.TotalTime,
		ResponseSize:   len(resp.Body()),
		IsConnReused:   traceInfo.IsConnReused,
		IsConnWasIdle:  traceInfo.IsConnWasIdle,
		ConnIdleTime:   traceInfo.ConnIdleTime,
		RequestAttempt: traceInfo.RequestAttempt,
		RemoteAddr:     traceInfo.RemoteAddr.String(),
	}, nil
}
