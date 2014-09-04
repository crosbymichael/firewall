firewall is a simple go app that adds and removes iptables rules blocking external access to ports
unless it matches a specific ip.

#### Apply the rules within the config
```bash
firewall config.json
```

#### Remove the rules within the config
```bash
firewall -rm config.json
```

#### Sample config
```json
[
    {
        "interface": "eth0",
        "proto": "tcp",
        "port": 8080,
        "allow": [
            "107.170.333.222"
        ]
    }
]
```
