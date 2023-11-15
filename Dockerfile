FROM ubuntu:20.04

COPY ./bankdownloader /opt/app/
WORKDIR /opt/app
ENTRYPOINT ["/opt/app/bank-downloaders"]
