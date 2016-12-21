FROM alpine
RUN apk -q update && apk -q add ca-certificates
ADD geoserver /
CMD ["/geoserver"]
