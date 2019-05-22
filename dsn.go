package gomaxcompute

import (
	"fmt"
	"net/url"
)

type Config struct {
	AccessID  string
	AccessKey string
	Project   string
	Endpoint  string
	Others    map[string]string
}

func (cfg *Config) FormatDSN() string {
	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return ""
	}
	u.User = url.UserPassword(cfg.AccessID, cfg.AccessKey)
	q := u.Query()
	q.Set(currentProject, cfg.Project)
	for k, v := range cfg.Others {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// ParseDSN parses the DSN string to a Config, dsn: odpsURL.
func ParseDSN(dsn string) (*Config, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	username := u.User.Username()
	password, _ := u.User.Password()
	proj := u.Query().Get(currentProject)
	if username == "" || password == "" || proj == "" {
		return nil, fmt.Errorf("invalid odps url %v", u)
	}
	endpoint := (&url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}).String()
	ps := make(map[string]string)
	for k, vs := range u.Query() {
		if k != currentProject {
			if len(vs) == 0 {
				ps[k] = ""
			}
			ps[k] = vs[0]
		}
	}

	return &Config{
		AccessID:  u.User.Username(),
		AccessKey: password,
		Project:   proj,
		Endpoint:  endpoint,
		Others:    ps,
	}, nil
}
