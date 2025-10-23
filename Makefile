STEAMPIPE_INSTALL_DIR ?= ~/.steampipe

install:
	go build -o  $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/francois2metz/ovh@latest/steampipe-plugin-ovh.plugin *.go
