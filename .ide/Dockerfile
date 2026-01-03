ARG ALPINE_VERSION=3.23
ARG VERSION=4.107.0
FROM golang:alpine${ALPINE_VERSION} AS infra

FROM docker.cnb.cool/viscropst/coder-alpine/base:${ALPINE_VERSION}-${VERSION} as ide
COPY --from=infra --link /usr/local/go /usr/local/

# install vscode and extension
RUN code-server --install-extension ms-ceintl.vscode-language-pack-zh-hans@1.101.0 &&\
    code-server --install-extension tamasfe.even-better-toml@0.21.2 &&\
    code-server --install-extension golang.go@0.50.0 &&\
    code-server --install-extension tencent-cloud.coding-copilot@3.6.9917994 &&\
    echo done

COPY .ide/languagepacks.json /root/.local/share/code-server/languagepacks.json
COPY .ide/extensions.json /root/.local/share/code-server/extensions/extensions.json
COPY .ide/user.argv.json /root/.local/share/code-server/User/argv.json

ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8