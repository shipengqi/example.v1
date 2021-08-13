package consul

import (
	"fmt"
	naming2 "github.com/shipengqi/qim/library/naming"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/shipengqi/qim/library/log"
)

const (
	KeyProtocol  = "protocol"
	KeyHealthURL = "health_url"
)

type watcher struct {
	Service   string
	Callback  func([]naming2.Registration)
	WaitIndex uint64
	Quit      chan struct{}
}

type consul struct {
	sync.RWMutex

	cli      *api.Client
	watchers map[string]*watcher
}

func New(url string) (naming2.Interface, error) {
	conf := api.DefaultConfig()
	conf.Address = url
	cli, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &consul{
		cli: cli,
	}, nil
}

func (c *consul) Find(name string, tags ...string) ([]naming2.Registration, error) {
	services, _, err := c.load(name, 0, tags...)
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (c *consul) Register(s naming2.Registration) error {
	registration := &api.AgentServiceRegistration{
		ID:      s.ServiceID(),
		Name:    s.ServiceName(),
		Address: s.Address(),
		Port:    s.Port(),
		Tags:    s.GetTags(),
		Meta:    s.GetMeta(),
	}
	if registration.Meta == nil {
		registration.Meta = make(map[string]string)
	}
	registration.Meta[KeyProtocol] = s.GetProtocol()

	// consul health check
	healthURL := s.GetMeta()[KeyHealthURL]
	if len(healthURL) > 0 {
		check := &api.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s_normal", s.ServiceID()),
			HTTP:                           healthURL,
			Timeout:                        "1s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "20s",
		}
		registration.Check = check
	}

	err := c.cli.Agent().ServiceRegister(registration)
	return err
}

func (c *consul) Deregister(serviceID string) error {
	return c.cli.Agent().ServiceDeregister(serviceID)
}

func (c *consul) Subscribe(serviceName string, callback func([]naming2.Registration)) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.watchers[serviceName]; ok {
		return errors.New("serviceName has already been registered")
	}
	w := &watcher{
		Service:  serviceName,
		Callback: callback,
		Quit:     make(chan struct{}, 1),
	}
	c.watchers[serviceName] = w

	go c.watch(w)
	return nil
}

func (c *consul) Unsubscribe(serviceName string) error {
	c.Lock()
	defer c.Unlock()
	wh, ok := c.watchers[serviceName]

	delete(c.watchers, serviceName)
	if ok {
		close(wh.Quit)
	}
	return nil
}

// refresh service registration
func (c *consul) load(name string, waitIndex uint64, tags ...string) ([]naming2.Registration, *api.QueryMeta, error) {
	opts := &api.QueryOptions{
		UseCache:  true,
		MaxAge:    time.Minute, // MaxAge limits how old a cached value will be returned if UseCache is true.
		WaitIndex: waitIndex,
	}
	catalogServices, meta, err := c.cli.Catalog().ServiceMultipleTags(name, tags, opts)
	if err != nil {
		return nil, meta, err
	}

	services := make([]naming2.Registration, len(catalogServices))
	for i, s := range catalogServices {
		if s.Checks.AggregatedStatus() != api.HealthPassing {
			log.Debug().Msgf("load service: id:%s name:%s %s:%d Status:%s", s.ServiceID, s.ServiceName, s.ServiceAddress, s.ServicePort, s.Checks.AggregatedStatus())
			continue
		}
		dsvc := naming2.NewDefaultService(s.ServiceID, s.ServiceName)
		dsvc.SetAddress(s.ServiceAddress)
		dsvc.SetPort(s.ServicePort)
		dsvc.SetProtocol(s.ServiceMeta[KeyProtocol])
		dsvc.SetTags(s.ServiceTags)
		dsvc.SetMeta(s.ServiceMeta)
		services[i] = dsvc
	}
	log.Debug().Msgf("load service: %v, meta:%v", services, meta)
	return services, meta, nil
}

func (c *consul) watch(wh *watcher) {
	stopped := false

	var doWatch = func(service string, callback func([]naming2.Registration)) {
		services, meta, err := c.load(service, wh.WaitIndex) // <-- blocking until services has changed
		if err != nil {
			log.Warn().Err(err)
			return
		}
		select {
		case <-wh.Quit:
			stopped = true
			log.Info().Msgf("watch %s stopped", wh.Service)
			return
		default:
		}

		wh.WaitIndex = meta.LastIndex
		if callback != nil {
			callback(services)
		}
	}

	// build WaitIndex
	doWatch(wh.Service, nil)
	for !stopped {
		doWatch(wh.Service, wh.Callback)
	}
}
