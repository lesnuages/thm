package vmware

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

func NewClient(ctx context.Context, conf *Config) (*vim25.Client, error) {
	u, err := soap.ParseURL(conf.GovcURL)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, fmt.Errorf("invalid URL: %s", conf.GovcURL)
	}

	u.User = url.UserPassword(conf.GovcUsername, conf.GovcPassword)

	if u.User == nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: conf.GovcInsecure,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Run(conf *Config, f func(context.Context, *vim25.Client) error) error {
	var err error
	var c *vim25.Client

	ctx := context.Background()
	c, err = NewClient(ctx, conf)
	if err == nil {
		return f(ctx, c)
	} else {
		log.Fatal("could not create client\n", err)
	}
	return err
}
