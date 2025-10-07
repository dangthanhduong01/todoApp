#!/usr/bin/env python3
from PIL import Image, ImageDraw
import os

# Create a 128x128 icon
size = 128
image = Image.new('RGBA', (size, size), (0, 0, 0, 0))
draw = ImageDraw.Draw(image)

# Background circle
draw.ellipse([4, 4, size-4, size-4], fill='#4A90E2', outline='#2E5C8A', width=4)

# Todo list background (white rectangle)
draw.rectangle([28, 24, 100, 96], fill='white', outline='#E0E0E0', width=2)

# Header line
draw.line([36, 36, 92, 36], fill='#333333', width=3)

# First todo item (completed - green checkbox)
draw.rectangle([36, 46, 44, 54], fill='#4CAF50', outline='#2E7D32', width=1)
# Checkmark
draw.line([38, 50, 40, 52], fill='white', width=2)
draw.line([40, 52, 42, 48], fill='white', width=2)
# Text line
draw.line([48, 50, 84, 50], fill='#333333', width=2)

# Second todo item (uncompleted)
draw.rectangle([36, 60, 44, 68], fill='white', outline='#666666', width=1)
draw.line([48, 64, 84, 64], fill='#333333', width=2)

# Third todo item (uncompleted)
draw.rectangle([36, 74, 44, 82], fill='white', outline='#666666', width=1)
draw.line([48, 78, 80, 78], fill='#333333', width=2)

# Plus button (add new item)
draw.ellipse([74, 82, 86, 94], fill='#FF5722', outline='#D32F2F', width=1)
draw.line([77, 88, 83, 88], fill='white', width=2)
draw.line([80, 85, 80, 91], fill='white', width=2)

# Save the image
image.save('todoapp.png')
print("Icon created: todoapp.png")
