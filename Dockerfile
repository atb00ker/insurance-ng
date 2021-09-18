FROM golang:1.16.6-alpine3.14 AS BASE

WORKDIR /go/src/app

FROM BASE AS BUILD

# Get dependencies
RUN apk add --update nodejs npm
COPY package.json .
RUN npm install --include=dev
COPY /src/server/ ./src/server/
COPY go.mod go.sum .babelrc webpack.config.js .env ./
RUN go install -v ./...
COPY /src/views/ ./src/views/

# Build
RUN npm run react-build
RUN go build -o insuranceng ./src/server/

FROM BASE

COPY --from=BUILD /go/src/app/dist/ ./dist
COPY --from=BUILD /go/src/app/insuranceng ./insuranceng

ENTRYPOINT ["./insuranceng"]
# ENTRYPOINT ["yes"]
