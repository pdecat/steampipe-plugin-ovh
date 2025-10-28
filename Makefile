STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
STEAMPIPE_PLUGIN_VERSION ?= latest

install:
	go build -o  $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/francois2metz/ovh@$(STEAMPIPE_PLUGIN_VERSION)/steampipe-plugin-ovh.plugin *.go
