# Makefile for Modern Mermaid Desktop

.PHONY: dev build build-frontend install clean

# Helper target to setup icons
setup-icons:
	mkdir -p desktop/build
	cp public/icon-512.png desktop/build/appicon.png

# Default target - Start development mode
dev: setup-icons
	cd desktop && wails dev

# Build the desktop application for production
build: setup-icons
	cd desktop && wails build -debug

# Build the frontend (called by Wails)
# This command:
# 1. Compiles TypeScript
# 2. Builds Vite project
# 3. Cleans old assets
# 4. Copies new assets to desktop/assets
build-frontend:
	pnpm tsc -b
	VITE_IS_DESKTOP=true pnpm vite build
	rm -rf desktop/assets
	cp -r dist desktop/assets

# Install all dependencies
install:
	pnpm install
	cd desktop && go mod tidy

# Clean build artifacts
clean:
	rm -rf dist
	rm -rf desktop/assets
	rm -rf desktop/build/bin

