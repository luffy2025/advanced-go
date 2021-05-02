package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 同时启动 2 个 http server，并在 5 秒后关闭
func main() {
	healthServer := newServer(withAddr("localhost:8081"),
		withHandle("/health", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(200)
		}))

	mainServer := newServer(withAddr("localhost:8080"),
		withHandle("/", func(writer http.ResponseWriter, request *http.Request) {
			_, err := writer.Write([]byte("welcome."))
			if err != nil {
				log.Printf("%+v\n", errors.Wrap(err, "http handle error. path: /"))
			}
		}), withHandle("/hello", func(writer http.ResponseWriter, request *http.Request) {
			_, err := writer.Write([]byte("hello gopher."))
			if err != nil {
				log.Printf("%+v\n", errors.Wrap(err, "http handle error. path: /hello"))
			}
		}))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <- time.After(5 * time.Second):
			cancel()
		}
	}()

	err := run(ctx, healthServer, mainServer)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(ctx context.Context, servers ...*server) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, svr := range servers {
		svr := svr
		g.Go(func() error {
			<-ctx.Done()
			return svr.stop(ctx)
		})
		g.Go(func() error {
			return svr.start()
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sg := <-c:
				return errors.Errorf("os signal: %s", sg)
			}
		}
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

type server struct {
	*http.Server

	addr     string
	handlers map[string]http.HandlerFunc
}

func newServer(opts ...option) *server {
	svr := &server{
		handlers: make(map[string]http.HandlerFunc),
	}

	for _, opt := range opts {
		opt(svr)
	}

	mux := http.NewServeMux()
	for p, h := range svr.handlers {
		mux.HandleFunc(p, h)
	}
	svr.Server = &http.Server{Addr: svr.addr, Handler: mux}

	return svr
}

func (s *server) start() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return errors.Wrapf(err, "http listen error. address:%s", s.addr)
	}

	err = s.Serve(lis)
	if err != nil {
		return errors.Wrapf(err, "http serve error. address:%s", s.addr)
	}
	return nil
}

func (s *server) stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

type option func(*server)

func withHandle(path string, handle http.HandlerFunc) option {
	return func(s *server) {
		if s.handlers == nil {
			s.handlers = make(map[string]http.HandlerFunc)
		}
		s.handlers[path] = handle
	}
}

func withAddr(addr string) option {
	return func(s *server) {
		s.addr = addr
	}
}
