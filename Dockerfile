FROM scratch
COPY ./dist/bank-downloaders /opt/app/
WORKDIR /opt/app
ENTRYPOINT ["/opt/app/bank-downloaders"]
