# Define variables
GO_VERSION = 1.23.2
GO_TAR = go$(GO_VERSION).linux-amd64.tar.gz
GO_URL = https://go.dev/dl/$(GO_TAR)
INSTALL_DIR = /usr/local

# Define targets
all: download extract set-path

download:
	wget $(GO_URL)

extract:
	sudo tar -C $(INSTALL_DIR) -xzf $(GO_TAR)

set-path:
	@echo "Para añadir Go al PATH, ejecuta el siguiente comando:"
	@echo "export PATH=\$$PATH:$(INSTALL_DIR)/go/bin"

# Clean up downloaded file
clean:
	rm -f $(GO_TAR)
