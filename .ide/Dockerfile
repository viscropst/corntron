FROM golang:bookworm as build

RUN sed -i "s@http://deb.debian.org/debian@http://mirrors.cloud.tencent.com/debian@g" /etc/apt/sources.list.d/debian.sources && \
    apt-get update && \
    apt-get install -y --no-install-recommends git openssh-server openssh-client && \
    rm -rf /var/lib/apt/lists/*

# install vscode and extension
RUN curl -fsSL https://ghproxy.net/https://raw.githubusercontent.com/cdr/code-server/main/install.sh | sh &&\
    code-server --install-extension ms-ceintl.vscode-language-pack-zh-hans@1.101.0 &&\
    code-server --install-extension tamasfe.even-better-toml@0.21.0 &&\
    code-server --install-extension golang.go@0.50.0 &&\
    code-server --install-extension tencent-cloud.coding-copilot@3.6.9917994 &&\
    echo done

COPY .ide/languagepacks.json /root/.local/share/code-server/languagepacks.json
COPY .ide/extensions.json /root/.local/share/code-server/extensions/extensions.json
COPY .ide/user.argv.json /root/.local/share/code-server/User/argv.json

ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8