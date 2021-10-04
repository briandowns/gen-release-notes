# gen-release-notes

gen-release-notes is a simple utility to generate release-notes for K3s and RKE2.


## Installation

```
go get github.com/briandowns/gen-release-notes
```


### Examples

Generate release notes for k3s v1.21.5+k3s1

```sh
gen-release-notes -r k3s -m v1.21.5+k3s1 -p v1.21.4+k3s1
```

Or via Docker

```sh
docker run --rm -it briandowns/gen-release-notes:v0.2.0 gen-release-notes -r k3s -m v1.21.5+k3s1 -p v1.21.4+k3s1
```

## Contributions

* File Issue with details of the problem, feature request, etc.
* Submit a pull request and include details of what problem or feature the code is solving or implementing.
