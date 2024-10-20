package config

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	CompatiableOpenAIEnabled bool

	CompatiableOpenAIEndpoint string

	CompatiableOpenAIKey string

	CompatiableOpenAIModelName string

	DisableTelemetry bool

	PremiumCode string

	// Use for telemetry. We generate a random ID for each user.
	ClientID string
}

func xdgConfigDir() string {
	xdgConfigDir := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigDir == "" {
		xdgConfigDir = os.Getenv("HOME") + "/.config"
	}
	return xdgConfigDir
}

func configDir() string {
	d := xdgConfigDir() + "/so"
	_, err := os.Stat(d)

	// Check if dir is exists
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(d, 0755)
		if err != nil {
			log.Printf("Unable to make config dir `%s': %v", d, err)
			return ""
		}
	}

	return d
}

func read() (*Config, error) {
	fname := configDir() + "/config.json"
	buf, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = json.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func write(c *Config) error {
	fname := configDir() + "/config.json"
	buf, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(fname, buf, 0644)
}

func Get() *Config {
	c, err := read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read config file: %v", err)
	}

	return c
}

func Save(c *Config) error {
	return write(c)
}

func init() {
	// Ensure config file
	c, _ := read()
	if c == nil {
		c = &Config{}
	}

	if c.ClientID == "" {
		buf1 := make([]byte, 4)
		buf2 := bytes.NewBuffer(make([]byte, 8))
		for i := 0; i < 4; i++ {
			buf1[i] = byte(rand.Intn(256))
		}
		binary.Write(buf2, binary.LittleEndian, int64(time.Now().UnixNano()))
		c.ClientID = base64.URLEncoding.EncodeToString(append(buf1, buf2.Bytes()...))
	}

	write(c)
}
