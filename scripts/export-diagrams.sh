#!/bin/bash

set -e

DOCS_DIR="./docs"
OUTPUT_DIR="./diagrams"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Project Diagrams Export Tool${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if ! command -v mmdc &>/dev/null; then
  echo -e "${YELLOW}Warning: mermaid-cli (mmdc) is not installed${NC}"
  echo "Install with: npm install -g @mermaid-js/mermaid-cli"
  echo ""
  echo "Skipping diagram export..."
  exit 0
fi

mkdir -p "$OUTPUT_DIR"
echo -e "${GREEN}✓${NC} Created output directory: $OUTPUT_DIR"

export_diagram() {
  local input_file=$1
  local output_name=$2

  echo -e "${BLUE}Exporting${NC} $output_name..."

  if mmdc -i "$input_file" -o "$OUTPUT_DIR/${output_name}.png" -b transparent 2>/dev/null; then
    echo -e "  ${GREEN}✓${NC} PNG exported"
  else
    echo -e "  ${YELLOW}⚠${NC} PNG export failed"
  fi

  if mmdc -i "$input_file" -o "$OUTPUT_DIR/${output_name}.svg" -b transparent 2>/dev/null; then
    echo -e "  ${GREEN}✓${NC} SVG exported"
  else
    echo -e "  ${YELLOW}⚠${NC} SVG export failed"
  fi

}

echo ""
echo -e "${BLUE}Exporting diagrams...${NC}"
echo ""

if [ -f "$DOCS_DIR/ARCHITECTURE_DIAGRAM.md" ]; then
  export_diagram "$DOCS_DIR/ARCHITECTURE_DIAGRAM.md" "architecture"
fi

if [ -f "$DOCS_DIR/NETWORK_SERVICES_DIAGRAM.md" ]; then
  export_diagram "$DOCS_DIR/NETWORK_SERVICES_DIAGRAM.md" "network_services"
fi

if [ -f "$DOCS_DIR/ENTITY_RELATIONSHIP_DIAGRAM.md" ]; then
  export_diagram "$DOCS_DIR/ENTITY_RELATIONSHIP_DIAGRAM.md" "entity_relationship"
fi

if [ -f "$DOCS_DIR/SEQUENCE_DIAGRAMS.md" ]; then
  export_diagram "$DOCS_DIR/SEQUENCE_DIAGRAMS.md" "sequence_diagrams"
fi

echo ""
echo -e "${BLUE}Creating archive...${NC}"
ARCHIVE_NAME="diagrams_${TIMESTAMP}.tar.gz"
tar -czf "$OUTPUT_DIR/$ARCHIVE_NAME" -C "$OUTPUT_DIR" \
  --exclude="*.tar.gz" \
  $(ls "$OUTPUT_DIR" | grep -v ".tar.gz") 2>/dev/null || true

if [ -f "$OUTPUT_DIR/$ARCHIVE_NAME" ]; then
  echo -e "${GREEN}✓${NC} Archive created: $ARCHIVE_NAME"
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Export Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Output directory: $OUTPUT_DIR"
echo "Files exported:"
ls -lh "$OUTPUT_DIR" | grep -v "^total" | awk '{print "  - " $9 " (" $5 ")"}'
echo ""
echo -e "${BLUE}Tip:${NC} View diagrams in your browser or image viewer"
echo -e "${BLUE}Tip:${NC} Share the archive with team members"
echo ""
