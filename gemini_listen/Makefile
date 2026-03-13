BINARY_NAME=gemini-listen
SERVICE_NAME=gemini-listen.service

INSTALL_DIR=$(HOME)/.local/bin
SYSTEMD_USER_DIR=$(HOME)/.config/systemd/user

.PHONY: all build run clean install uninstall restart status docker-build

all: build

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

install: build
	mkdir -p $(INSTALL_DIR)
	mkdir -p $(SYSTEMD_USER_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	cp $(SERVICE_NAME) $(SYSTEMD_USER_DIR)/
	systemctl --user daemon-reload
	systemctl --user enable $(SERVICE_NAME)
	systemctl --user start $(SERVICE_NAME)

uninstall:
	systemctl --user stop $(SERVICE_NAME) || true
	systemctl --user disable $(SERVICE_NAME) || true
	rm -f $(SYSTEMD_USER_DIR)/$(SERVICE_NAME)
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	systemctl --user daemon-reload

restart:
	systemctl --user restart $(SERVICE_NAME)

status:
	systemctl --user status $(SERVICE_NAME)

docker-build:
	docker build -t $(BINARY_NAME) .
