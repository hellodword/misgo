# misgo

## Motivation

Although some Go modules are open-sourced, they use custom domain names for their module paths. `GOSUMDB` ensures that the provided sum is untampered, but when running `go get foo` for the first time, it appears there's still an implicit trust that the code fetched is identical to the source code hosted on GitHub/GitLab and hasn't been altered. I hope there’s a way to verify this.
