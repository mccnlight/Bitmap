# Bitmap

## Overview

The **Bitmap Image Processor** is a command-line tool designed to read, manipulate, and save bitmap (BMP) image files. It supports various image processing operations such as mirroring, filtering, rotating, and cropping. The tool is written in Go and adheres to the BMP file format specifications, specifically supporting **24-bit uncompressed BMP files**.

## Features

### 1. Header Information

**Command:** `header`

**Description:** Displays detailed information about the BMP file header, including file type, size, dimensions, and more.

**Usage:**
```bash
./bitmap header <source_file>
```
**Example:**
```bash
./bitmap header sample.bmp
```

---

### 2. Image Manipulation

**Command:** `apply`

**Description:** Applies various image processing operations to the input BMP file and saves the result to an output file. Multiple operations can be combined and applied sequentially.

**Usage:**
```bash
./bitmap apply [options] <source_file> <output_file>
```

#### Options:

- **Mirror**: Mirrors the image **horizontally** or **vertically**.
  ```bash
  --mirror=<horizontal|vertical>
  ```

- **Filter**: Applies a filter to the image. Available filters include:
  - `blue`: Retains only the blue channel.
  - `red`: Retains only the red channel.
  - `green`: Retains only the green channel.
  - `grayscale`: Converts the image to grayscale.
  - `negative`: Applies a negative filter.
  - `pixelate`: Pixelates the image.
  - `blur`: Applies a blur effect.
  ```bash
  --filter=<blue|red|green|grayscale|negative|pixelate|blur>
  ```

- **Rotate**: Rotates the image by a specified angle (90, 180, 270 degrees, or their equivalents).
  ```bash
  --rotate=<right|left|90|-90|180|-180|270|-270>
  ```

- **Crop**: Crops the image based on the specified offset and dimensions.
  ```bash
  --crop=<offsetX-offsetY-width-height>
  ```

**Example:**
```bash
./bitmap apply --mirror=horizontal --rotate=right --filter=negative sample.bmp output.bmp
```

---

### 3. Help

**Description:** Displays usage instructions for the program or specific commands.

**Usage:**
```bash
./bitmap -h
./bitmap header --help
./bitmap apply --help
```

---

## Installation

### Clone the Repository:
```bash
git clone <repository_url>
cd bitmap
```

### Build the Project:
```bash
go build -o bitmap .
```

### Run the Program:
```bash
./bitmap <command> [arguments]
```

---

## Usage Examples

### Display Header Information:
```bash
./bitmap header sample.bmp
```

### Mirror an Image Horizontally:
```bash
./bitmap apply --mirror=horizontal sample.bmp mirrored.bmp
```

### Apply Multiple Filters:
```bash
./bitmap apply --filter=grayscale --filter=blur sample.bmp filtered.bmp
```

### Rotate and Crop an Image:
```bash
./bitmap apply --rotate=90 --crop=50-50-200-200 sample.bmp cropped.bmp
```

---

## Error Handling

The program will exit with a non-zero status code and display an error message if:

- The input file is not a valid **24-bit uncompressed BMP file**.
- Invalid arguments or options are provided.
- The file cannot be read or written.

---

## Supported File Format

The program **only supports 24-bit uncompressed BMP files**. If the input file does not meet these criteria, the program will exit with an error.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Author

**Your Name**  
Email: nurpeissov.a0818@gmail.com  
GitHub: [mccnlight](https://github.com/mccnlight)
