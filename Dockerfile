FROM ubuntu:20.04

RUN curl -fsSL https://get.jetpack.io/devbox | bash

WORKDIR /code

COPY devbox.json devbox.json
COPY devbox.lock devbox.lock

RUN devbox run -- echo "Installed Packages."
RUN devbox shellenv --init-hook >> ~/.profile

COPY ./bankdownloader /code

ENTRYPOINT ["/code/bankdownloader"]