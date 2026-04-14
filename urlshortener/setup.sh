#!/bin/bash
# setup.sh — First-time setup script for the URL Shortener
#
# Run this once after cloning/downloading the project:
#   chmod +x setup.sh
#   ./setup.sh

set -e  # exit immediately if any command fails

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  URL Shortener — Initial Setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# ── Step 1: Check prerequisites ───────────────────────────────────────────
echo "🔍 Checking prerequisites..."

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Download it from https://go.dev/dl/"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Download it from https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker compose version &> /dev/null 2>&1; then
    echo "❌ Docker Compose v2 is not available. Make sure Docker Desktop is up to date."
    exit 1
fi

echo "✅ Go version:     $(go version)"
echo "✅ Docker version: $(docker --version)"
echo ""

# ── Step 2: Create .env from example ──────────────────────────────────────
if [ ! -f .env ]; then
    echo "📄 Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ .env created"
    echo ""
    echo "   ℹ️  Default config uses Docker Compose postgres credentials."
    echo "   ℹ️  Edit .env if you want to use a different PostgreSQL server."
    echo ""
else
    echo "✅ .env already exists — skipping"
fi

# ── Step 3: Download Go dependencies ──────────────────────────────────────
echo "📦 Downloading Go dependencies..."
go mod tidy
echo "✅ Dependencies ready (go.sum generated)"
echo ""

# ── Step 4: Offer to start Docker stack ───────────────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Setup complete! Next steps:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  Option A — Run with Docker (recommended):"
echo "    docker compose up --build"
echo ""
echo "  Option B — Run locally (needs local PostgreSQL):"
echo "    go run cmd/main.go"
echo ""
echo "  Once running, test it:"
echo "    curl -X POST http://localhost:3000/api/urls/ \\"
echo "      -H 'Content-Type: application/json' \\"
echo "      -d '{\"original_url\": \"https://google.com\", \"custom_code\": \"goog\"}'"
echo ""
echo "  pgAdmin (DB browser): http://localhost:5050"
echo "    Email:    admin@admin.com"
echo "    Password: admin"
echo ""

read -p "🚀 Start Docker stack now? (y/n): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🐳 Starting Docker Compose stack..."
    docker compose up --build -d
    echo ""
    echo "⏳ Waiting for services to be healthy..."
    sleep 5
    docker compose ps
    echo ""
    echo "✅ Stack is running!"
    echo "   API:     http://localhost:3000"
    echo "   pgAdmin: http://localhost:5050"
    echo ""
    echo "   View logs:  docker compose logs -f app"
    echo "   Stop:       docker compose down"
fi