# Contributing

Prerequisite:

- Docker
- Golang
- NPM

```bash
# start etcd
docker run --rm --name etcd-server \
    -p 2379:2379 \
    -p 2380:2380 \
    -e ALLOW_NONE_AUTHENTICATION=yes \
    -e ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 \
    bitnami/etcd
```

```bash
# run the server
go run .
```

```bash
# run the UI
npm start
```

```bash
# install pre-commit hook
cat > .git/hooks/pre-commit <<EOF
set -eux
go vet .
goimports -w .
npx prettier --write .
git diff --exit-code
go generate .
go install .
EOF
chmod +x .git/hooks/pre-commit
```

```bash
# build the binary
go generate .
go build .
```
