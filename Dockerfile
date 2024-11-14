# From alpine lates load the binary in dist folder
FROM alpine:latest
COPY dist/crispy-enigma_linux_386/crispy-enigma /crispy-enigma
CMD sh