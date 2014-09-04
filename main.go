package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	remove bool
)

type config struct {
	Interface string   `json:"interface,omitempty"`
	Proto     string   `json:"proto,omitempty"`
	Port      int      `json:"port,omitempty"`
	Allow     []string `json:"allow,omitempty"`
}

func init() {
	flag.BoolVar(&remove, "rm", false, "remove the rules from the system")
	flag.Parse()
}

func loadConfig(path string) ([]*config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c []*config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}

// iptables -A INPUT -i eth0 -p tcp --dport 8080 -j DROP
// iptables -I INPUT -i eth0 -s 127.0.0.1 -p tcp --dport 8080 -j ACCEPT
func process(c *config) error {
	port := fmt.Sprint(c.Port)
	// process the DROP rule
	if err := iptables("-A", "INPUT", "-i", c.Interface, "-p", c.Proto, "--dport", port, "-j", "DROP"); err != nil {
		return err
	}

	for _, a := range c.Allow {
		if err := iptables("-I", "INPUT", "-i", c.Interface, "-s", a, "-p", c.Proto, "--dport", port, "-j", "ACCEPT"); err != nil {
			return err
		}
	}

	return nil
}

func iptables(args ...string) error {
	// HACK but it works for what we want
	if remove {
		args[0] = "-D"
	}

	return exec.Command("iptables", args...).Run()
}

func main() {
	path := flag.Arg(0)
	if path == "" {
		log.Fatal("config path must be set")
	}

	configs, err := loadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range configs {
		if err := process(c); err != nil {
			log.Fatal(err)
		}
	}
}
