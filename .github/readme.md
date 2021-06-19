### unbound-manager

#### Note: We are arching this repo in favor of [content-blocker](https://github.com/complexorganizations/content-blocker) since Unbound is unable to grow to the standard we require.

---
### Installation
Lets first use `curl` and save the file in `/usr/local/bin/`
```
curl https://raw.githubusercontent.com/complexorganizations/unbound-manager/main/unbound-manager.sh --create-dirs -o /usr/local/bin/unbound-manager.sh
```
```
chmod +x /usr/local/bin/unbound-manager.sh
```
It's finally time to execute the script
```
bash /usr/local/bin/unbound-manager.sh
```
---
### Features
- Install, manage your own DNS
- DNSSEC Validation
- DNS Proxy
- Blocking based on DNS

---
### Variants
| Variants               |
| ---------------------  |
| [Host](https://raw.githubusercontent.com/complexorganizations/unbound-manager/main/configs/host) |

---
### Author
* Name: Prajwal Koirala
* Website: [prajwalkoirala.com](https://www.prajwalkoirala.com)

---
### Credits
Open Source Community

---
### License
Copyright © [Prajwal](https://github.com/prajwal-koirala)

This project is unlicensed.
