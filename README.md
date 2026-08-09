# Super Liquid Soccer

Single-binary localhost launcher built with Go.  
Also available on GitHub pages.

## Repository layout

```
.
├── main.go
├── go.mod
├── game/
│   ├── index.html
│   ├── index.js
│   ├── index.pck
│   └── index.wasm
└── .github/
    └── workflows/
        └── build.yml
```

## 1. Run locally

```bash
go run .
# → opens http://127.0.0.1:8080 in your default browser
```

## 2. Build for your own platform

```bash
go build -ldflags="-s -w" -trimpath -o super-liquid-soccer .
./super-liquid-soccer
```