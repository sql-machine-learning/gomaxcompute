package gomaxcompute

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	// Regexp syntax: https://github.com/google/re2/wiki/Syntax
	reDSN   = regexp.MustCompile(`^([a-zA-Z0-9_-]+):([=a-zA-Z0-9_-]+)@([:a-zA-Z0-9/_.-]+)\?([^/]+)$`)
	reQuery = regexp.MustCompile(`^([a-zA-Z0-9_-]+)=([a-zA-Z0-9_-]*)$`)
)

const HINT_PREFIX = "hint_"

type Config struct {
	AccessID   string
	AccessKey  string
	Project    string
	Endpoint   string
	QueryHints map[string]string
}

func ParseDSN(dsn string) (*Config, error) {
	sub := reDSN.FindStringSubmatch(dsn)
	if len(sub) != 5 {
		return nil, fmt.Errorf("dsn %s doesn't match access_id:access_key@url?curr_project=project&scheme=http|https", dsn)
	}
	id, key, endpointURL := sub[1], sub[2], sub[3]

	var schemeArgs []string
	var currProjArgs []string
	var ok bool
	queryHints := make(map[string]string)

	querys, err := url.ParseQuery(sub[4])
	if err != nil {
		return nil, err
	}

	if schemeArgs, ok = querys["scheme"]; !ok || len(schemeArgs) != 1 {
		return nil, fmt.Errorf("dsn %s should have one scheme argument", dsn)
	}
	if currProjArgs, ok = querys[currentProject]; !ok || len(currProjArgs) != 1 {
		return nil, fmt.Errorf("dsn %s should have one current_project argument", dsn)
	}

	for k, v := range querys {
		// The query args such as hints_odps.sql.mapper.split_size=16
		// would be converted to the maxcompute query hints: {"odps.sql.mapper.split_size": "16"}
		if strings.HasPrefix(k, "hint_") {
			queryHints[k[5:]] = v[0]
		}
	}

	if schemeArgs[0] != "http" && schemeArgs[0] != "https" {
		return nil, fmt.Errorf("dsn %s 's scheme is neither http nor https", dsn)
	}

	config := &Config{
		AccessID:   id,
		AccessKey:  key,
		Project:    currProjArgs[0],
		Endpoint:   schemeArgs[0] + "://" + endpointURL,
		QueryHints: queryHints}

	return config, nil
}

func (cfg *Config) FormatDSN() string {
	pair := strings.Split(cfg.Endpoint, "://")
	if len(pair) != 2 {
		return ""
	}
	scheme, endpointURL := pair[0], pair[1]
	dsnFormt := fmt.Sprintf("%s:%s@%s?curr_project=%s&scheme=%s",
		cfg.AccessID, cfg.AccessKey, endpointURL, cfg.Project, scheme)
	if len(cfg.QueryHints) != 0 {
		for k, v := range cfg.QueryHints {
			dsnFormt = fmt.Sprintf("%s&hint_%s=%v", dsnFormt, k, v)
		}
	}
	return dsnFormt
}
