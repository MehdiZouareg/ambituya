# Nom de l'application
APP_NAME=ambituya

# Version de l'application (facultatif, tu peux ajouter une gestion de version)
VERSION=0.4.0

# Définir les variables d'environnement pour le build cross-platform
BINARY_UNIX=$(APP_NAME)
BINARY_WINDOWS=$(APP_NAME)-$(VERSION).exe

# Définir les commandes de build
.PHONY: all clean build build-windows build-linux

# Build par défaut pour ta machine (par ex: Linux/Mac)
all: build

# Nettoyage des builds précédents
clean:
	rm -f $(BINARY_UNIX) $(BINARY_WINDOWS)

# Build pour la plateforme actuelle (Unix/Mac)
build:
	GOOS=linux GOARCH=amd64 go build -o ./build/$(BINARY_UNIX) -ldflags="-X main.version=$(VERSION)" ./cmd/ambituya

# Build pour Windows 64-bit
build-windows:
	GOOS=windows GOARCH=amd64 go build -o ./build/$(BINARY_WINDOWS) -ldflags="-H windowsgui -X main.version=$(VERSION)" ./cmd/ambituya

# Build pour Windows 64-bit avec terminal
build-windows-debug:
	GOOS=windows GOARCH=amd64 go build -o ./build/$(BINARY_WINDOWS) -ldflags="-X main.version=$(VERSION)" ./cmd/ambituya

# Build pour Linux 64-bit
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/$(BINARY_UNIX) -ldflags="-X main.version=$(VERSION)" ./cmd/ambituya

run: 
	go run ./cmd/ambituya/main.go

# Build pour toutes les plateformes
build-all: build-windows build-linux

# Afficher la version de l'application
version:
	@echo $(VERSION)
