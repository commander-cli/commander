FROM golang:1.17

RUN mkdir -p $GOPATH/github.com/commander/commander-cli

# For some reason circleci will not mount the dir?
COPY . $GOPATH/github.com/commander/commander-cli/
WORKDIR $GOPATH/github.com/commander/commander-cli/

RUN curl https://s3.amazonaws.com/codeclimate/test-reporter/test-reporter-0.10.1-linux-amd64 --output test-reporter
RUN chmod +x test-reporter

RUN curl -L https://github.com/commander-cli/commander/releases/download/v2.4.0/commander-linux-amd64 --output /usr/bin/commander
RUN chmod +x /usr/bin/commander
