<!-- @format -->

# Product Image Upload Implementation Guide

## Overview

Your application now supports local filesystem image storage with a complete image management system. Users can upload images when creating products or upload them separately afterward.

## Key Components

### 1. **Storage Service** (`internal/storage/storage.go`)

- **Local filesystem storage** implementation
- Supports: JPG, JPEG, PNG, GIF, WEBP
- Max file size: 10MB (configurable)
- Unique filenames with timestamps and UUIDs
- Directory traversal protection

### 2. **Updated DTOs** (`internal/dto/product_dto.go`)

- `ProductResponse` - now includes `Images` array
- `ProductImageResponse` - image data structure
- `UploadProductImageRequest` - image upload payload

### 3. **Repository Layer** (`internal/repository/product_repo.go`)

- `CreateProductImage()` - save image to database
- `GetProductImages()` - fetch all product images (sorted by position)
- `DeleteProductImage()` - remove image record
- `UpdateProductImagePosition()` - reorder images
- `GetProductImage()` - fetch single image

### 4. **Service Layer** (`internal/service/product_service.go`)

- `GetProductWithImages()` - retrieves product with all images
- `CreateProductImage()` - handles image creation with validation
- `DeleteProductImage()` - deletes image and cleans up file
- `UpdateProductImagePosition()` - changes image order

### 5. **Handlers** (`internal/handlers/product_handlers.go`)

- Updated GetProduct to include images
- `UploadProductImage()` - multipart form handler
- `DeleteProductImage()` - removes image
- `UpdateProductImagePosition()` - reorders images

## API Endpoints

### Get Product with Images

```
GET /products?id={productId}
Response includes: Product details + Images array
```

### Upload Product Image

```
POST /products/{productId}/images
Content-Type: multipart/form-data

Form Parameters:
- image: File (required)
- position: int (optional, default: 0)

Response:
{
  "id": "image-uuid",
  "image_url": "http://localhost:8080/uploads/filename.jpg",
  "position": 0
}
```

### Delete Product Image

```
DELETE /images/{imageId}
Auth: Required (Vendor only)

Response:
{
  "message": "image deleted successfully"
}
```

### Update Image Position

```
PUT /images/{imageId}/position
Content-Type: application/json
Auth: Required (Vendor only)

Body:
{
  "position": 2
}

Response:
{
  "message": "image position updated successfully"
}
```

## Environment Variables

```
UPLOAD_DIR=./uploads              # Where files are stored (default: ./uploads)
BASE_URL=http://localhost:8080    # Base URL for serving images (default: http://localhost:8080)
```

## File Structure

```
./uploads/
├── 1704537600_a1b2c3d4.jpg
├── 1704537601_e5f6g7h8.png
└── 1704537602_i9j0k1l2.webp
```

## Usage Example (cURL)

### 1. Upload Product Image

```bash
curl -X POST "http://localhost:8080/products/product-id/images" \
  -H "Authorization: Bearer your-jwt-token" \
  -F "image=@/path/to/image.jpg" \
  -F "position=0"
```

### 2. Get Product with Images

```bash
curl "http://localhost:8080/products?id=product-id"
```

### 3. Delete Image

```bash
curl -X DELETE "http://localhost:8080/images/image-id" \
  -H "Authorization: Bearer your-jwt-token"
```

### 4. Reorder Image

```bash
curl -X PUT "http://localhost:8080/images/image-id/position" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"position": 2}'
```

## Security Features

✅ **Authorization**: Only vendors can upload/delete their own product images  
✅ **File Validation**: Only allowed image formats (JPG, PNG, GIF, WEBP)  
✅ **Size Limits**: Max 10MB per file  
✅ **Unique Names**: Prevents filename collisions  
✅ **Path Protection**: Prevents directory traversal attacks  
✅ **Cleanup**: Automatically removes files on deletion

## Database Schema

```sql
CREATE TABLE product_images (
    id CHAR(36) PRIMARY KEY,
    product_id CHAR(36) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_images_product
      FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);
```

## Future Enhancements

1. **S3 Migration**: Switch to AWS S3 by implementing new `Storage` interface
2. **Image Compression**: Automatically compress images before saving
3. **Thumbnails**: Generate thumbnail versions
4. **CDN**: Serve images through CloudFront/CloudFlare
5. **Batch Upload**: Upload multiple images at once
6. **Drag-and-drop Reordering**: Frontend integration

## Troubleshooting

**Images not saving?**

- Check `UPLOAD_DIR` has write permissions
- Verify disk space available
- Check file size doesn't exceed 10MB

**Images not serving?**

- Ensure `/uploads` route is accessible
- Check `BASE_URL` is correct
- Verify files exist in upload directory

**Upload fails with 403?**

- Ensure you're authenticated (JWT token)
- Verify you're a vendor
- Confirm product belongs to you
