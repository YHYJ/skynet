TARGET       := skynet
INSTALL_PATH := /usr/local/bin
ATTRIBUTION  := root

.PHONY: all tidy build install clean
all: build

help:
	@echo "usage: make [OPTIONS]"
	@echo "    help        Show this message"
	@echo "    tidy        Update project module dependencies"
	@echo "    build       Compile and generate executable file"
	@echo "    install     Install executable file"
	@echo "    clean       Clean build process files"

tidy:
	@go mod tidy
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully tidied up dependencies"

build:
	@go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/skynet/function.buildTime=`date +%s` -X github.com/yhyj/skynet/function.buildBy=$(USER)" -o $(TARGET)
	@echo -en "\x1b[32m[✔]\x1b[0m Successfully generated \x1b[32;1m$(TARGET)\x1b[0m"

install:
	@install --mode=755 --owner=$(ATTRIBUTION) --group=$(ATTRIBUTION) $(TARGET) $(INSTALL_PATH)/$(TARGET)
	@echo -e "\r\x1b[K\x1b[0m\x1b[32m[✔]\x1b[0m Successfully installed \x1b[32m$(TARGET)\x1b[0m"

clean:
	@rm -f $(TARGET)     && echo -e "    - Removed \x1b[32m$(TARGET)\x1b[0m"
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully cleaned process files"
