FROM golang:bookworm as build
 
# install vscode and extension
RUN curl -fsSL https://code-server.dev/install.sh | sh &&\
    code-server --install-extension ms-ceintl.vscode-language-pack-zh-hans@1.101.0 &&\
    code-server --install-extension tamasfe.even-better-toml &&\
    code-server --install-extension golang.go &&\
    code-server --install-extension tencent-cloud.coding-copilot &&\
    echo done

RUN apt-get update && \
    apt-get install -y --no-install-recommends git openssh-server openssh-client && \
    rm -rf /var/lib/apt/lists/*

COPY .ide/languagepacks.json /root/.local/share/code-server/languagepacks.json
COPY .ide/extensions.json /root/.local/share/code-server/extensions/extensions.json
COPY .ide/user.argv.json /root/.local/share/code-server/User/argv.json

ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8