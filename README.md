# misgo

## Motivation

Although some Go modules are open-sourced, they use custom domain names for their module paths. `GOSUMDB` ensures that the provided sum is untampered, but when running `go get foo` for the first time, it appears there's still an implicit trust that the code fetched is identical to the source code hosted on GitHub/GitLab and hasn't been altered. I hope there’s a way to verify this.

## How it works?

- find the repository: `document.querySelector('.UnitMeta-repo a').getAttribute('href')` in `https://pkg.go.dev/foo@bar`
  - https://pkg.go.dev/about#source-links
  - https://github.com/golang/gddo/wiki/Source-Code-Links
- checksum: https://github.com/golang/go/blob/807e01db4840e25e4d98911b28a8fa54244b8dfa/src/cmd/go/internal/modfetch/cache.go#L429
- gomodsum: https://github.com/golang/go/blob/807e01db4840e25e4d98911b28a8fa54244b8dfa/src/cmd/go/internal/modfetch/fetch.go#L647-L652
- gosumdb: https://github.com/ProjectSerenity/firefly/blob/0effba12f4ea172166e098e955c0f5ecca29932f/tools/gomodproxy/gosumdb.go

## TODO

- [x] parse go.mod
- [x] parse go.sum
- [ ] recursively parse dependencies
- [ ] deal with pseudo version[^1]
- [ ] enhance fetchers: `https://github.com/FiloSottile/edwards25519/archive/<tag or commit>.zip` , see nixpkgs' fetchers https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/fetchgithub/default.nix
- [ ] create an evil example

## Ref

- https://esc.sh/blog/setting-up-a-git-http-server-with-nginx/

[^1]: https://news.ycombinator.com/item?id=37333545
