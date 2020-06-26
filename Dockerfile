FROM alpine
ADD travel-app /travel-app
ENTRYPOINT [ "/travel-app" ]