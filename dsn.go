package gomaxcompute

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Regexp syntax: https://github.com/google/re2/wiki/Syntax
	reDSN   = regexp.MustCompile(`^([a-zA-Z0-9_-]+):([=a-zA-Z0-9_-]+)@([:a-zA-Z0-9/_.-]+)\?([^/]+)$`)
	reQuery = regexp.MustCompile(`^([a-zA-Z0-9_-]+)=([a-zA-Z0-9_-]*)$`)
)

type Config struct {
	AccessID  string
	AccessKey string
	Project   string
	Endpoint  string
}

func ParseDSN(dsn string) (*Config, error) {
	sub := reDSN.FindStringSubmatch(dsn)
	if len(sub) != 5 {
		return nil, fmt.Errorf("dsn %s doesn't match access_id:access_key@url?curr_project=project&scheme=http|https", dsn)
	}
	id, key, endpointURL := sub[1], sub[2], sub[3]

	query := make(map[string]string)
	for _, s := range strings.Split(sub[4], "&") {
		pair := reQuery.FindStringSubmatch(s)
		if len(pair) != 3 {
			return nil, fmt.Errorf("dsn %s doesn't match access_id:access_key@url?curr_project=project&scheme=http|https", dsn)
		}
		if pair[1] != "scheme" && pair[1] != currentProject {
			return nil, fmt.Errorf("dsn %s 's query is neither scheme or %s", dsn, currentProject)
		}
		query[pair[1]] = pair[2]
	}
	if _, ok := query[currentProject]; !ok {
		return nil, fmt.Errorf("dsn %s doesn't have curr_project", dsn)
	}
	if _, ok := query["scheme"]; !ok {
		return nil, fmt.Errorf("dsn %s doesn't have scheme", dsn)
	}
	if query["scheme"] != "http" && query["scheme"] != "https" {
		return nil, fmt.Errorf("dsn %s 's scheme is neither http nor https", dsn)
	}

	config := &Config{
		AccessID:  id,
		AccessKey: key,
		Project:   query[currentProject],
		Endpoint:  query["scheme"] + "://" + endpointURL}

	return config, nil
}

func (cfg *Config) FormatDSN() string {
	pair := strings.Split(cfg.Endpoint, "://")
	if len(pair) != 2 {
		return ""
	}
	scheme, endpointURL := pair[0], pair[1]
	return fmt.Sprintf("%s:%s@%s?curr_project=%s&scheme=%s",
		cfg.AccessID, cfg.AccessKey, endpointURL, cfg.Project, scheme)
}
