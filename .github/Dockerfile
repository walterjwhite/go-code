FROM archlinux:latest

RUN mkdir -p /usr/local/bin

RUN curl -sSL https://raw.githubusercontent.com/walterjwhite/app-registry/refs/heads/main/install/Linux/setup/bin/app-install -o /usr/local/bin/app-install && \
    chmod +x /usr/local/bin/app-install

RUN conf_log_level=2 /usr/local/bin/app-install install && \
    /usr/local/bin/app-install programming-languages app-framework

RUN git config --global user.name "github-actions[bot]" && \
    git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"
