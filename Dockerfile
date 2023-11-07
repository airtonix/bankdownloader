FROM scratch
COPY ./bankdownloader /opt/app/
WORKDIR /opt/app
ENTRYPOINT ["/opt/app/bank-downloaders"]
