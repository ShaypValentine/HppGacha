FROM golang:1.20-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh openssl sqlite nano openrc busybox-openrc

LABEL maintainer="ShaypValentine <shaycaith@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# RUN (crontab -l ; echo "1	*	*	*	*	echo " ) | crontab -
# RUN (crontab -l ; echo "0	1	*	*	*	 cp /app/persistent/hppgacha.db /app/persistent/hppgachaBackUp.db" ) | crontab -
# RUN (crontab -l ; echo "0	*/2*	*	* sqlite3 /app/persistent/hppgacha.db \"UPDATE users SET available_rolls = available_rolls + 1 WHERE available_rolls< 4\"" ) | crontab -
# RUN (crontab -l ; echo "0 11 */2 * * sqlite3 /app/persistent/hppgacha.db \"UPDATE shadow_portals SET available_shadow_rolls = available_shadow_rolls + 1 WHERE available_shadow_rolls < 2\"" ) | crontab -
# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 443
EXPOSE 8080
EXPOSE 8443

# Run the executable
CMD ["./main"]
