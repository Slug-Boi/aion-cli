generate:
	echo "Generating..."
	cd src/tailwind && npm run build
	cd src && go run main.go generate 