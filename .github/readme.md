### unbound-manager

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
### Q&A
What's the best way for me to make my own list?
- Open the repo after forking and cloning it. Go ahead and change the `urls` struct, replacing the urls there with the lists you wish to use, and then just run the file using the command `go run main.go`.

What's the best way to add my own exclusions?
- Simply open the exclusion file, add a domain, and submit a pull request; if your pull request is merged, the domain will be excluded the next time the list is updated.

Is the list updated on a regular basis?
- We strive to update the list every 24 hours, but this cannot be guaranteed, and if it is not updated for any reason please let us know.

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
