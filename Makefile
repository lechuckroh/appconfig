BASE_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
REPORT_DIR=reports
TEST_DIR=.

GOCMD=go
GOCLEAN=$(GOCMD) clean
GOCOVER=$(GOCMD) tool cover
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

ENV_GOMOD_ON=GO111MODULE=on
ENV_TEST=MASTER_CONFIG_FILE=$(BASE_DIR)master-test.yml
GOBUILD_OPT=-mod=vendor -v
GOTEST_OPT=-mod=vendor -v

COVER_HTML=$(REPORT_DIR)/coverage.html
COVER_TXT=$(REPORT_DIR)/coverage.txt
COVER_XML=$(REPORT_DIR)/coverage.xml
JUNIT_REPORT=$(REPORT_DIR)/junit.xml

all: test

report-init:
	@mkdir -p $(REPORT_DIR)
install-cover-tool:
	@$(GOCMD) get github.com/t-yuki/gocover-cobertura
install-junit-tool:
	@$(GOCMD) get -u github.com/jstemmer/go-junit-report

# Test
test:
	@$(ENV_GOMOD_ON) $(ENV_TEST) $(GOTEST) $(GOTEST_OPT) -count=1 $(TEST_DIR)

cover: report-init
	@$(ENV_GOMOD_ON) $(GOTEST) $(GOTEST_OPT) -covermode=count -coverprofile=$(COVER_TXT) -coverpkg=$(TEST_DIR) $(TEST_DIR)
cover-html: cover
	@$(GOCOVER) -html=$(COVER_TXT) -o $(COVER_HTML)
cover-cout: cover
	@$(GOCOVER) -func=$(COVER_TXT)
cover-xml: install-cover-tool cover
	@gocover-cobertura < $(COVER_TXT) > $(COVER_XML)

cover-xml-junit: report-init install-junit-tool
	@make cover-xml | go-junit-report > $(JUNIT_REPORT)

# Clean
clean:
	@$(GOCLEAN)
	@rm -f $(COVER_TXT) $(COVER_HTML) $(COVER_XML) $(JUNIT_REPORT)

# Install dependencies to vendor/
vendor:
	@$(GOMOD) vendor
vendor-update:
	@$(GOGET) -u

# Test using docker
docker-test:
	@docker run --rm -v $(BASE_DIR):/usr/src/myapp -w /usr/src/myapp golang:latest make cover-xml-junit
