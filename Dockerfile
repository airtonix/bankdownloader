FROM chromedp/headless-shell@sha256:b38a99ae6cf90a37f069de85fe45ececb041e9f0da03c58ca9ccb5018d845cd5

WORKDIR /code

COPY ./bankdownloader /code

ENTRYPOINT ["/code/bankdownloader"]