FROM scratch
COPY cf-apps /bin/cf-apps
COPY ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/bin/cf-apps"]
