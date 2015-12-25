go build *.go
cf uninstall-plugin wildcard_plugin
cf install-plugin -f wildcard_plugin
