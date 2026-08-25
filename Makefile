.PHONY: generate css gallery test dev all

all: generate css gallery test

# templ colle « // templ: version: X » au « package ui » de chaque fichier
# généré, et Go prend alors ces 63 commentaires pour la documentation du
# paquet. Une ligne vide les détache sans rien perdre.
generate:
	templ generate
	@find . -name '*_templ.go' -not -path './node_modules/*' -print0 \
	 | xargs -0 perl -0pi -e 's{(// templ: version: [^\n]*\n)(package )}{$$1\n$$2}'

css:
	npx tailwindcss -i assets/app.css -o dist/app.css --minify

gallery: generate css
	go run ./cmd/gallery

test:
	go test ./...

dev: generate css
	go run ./example
