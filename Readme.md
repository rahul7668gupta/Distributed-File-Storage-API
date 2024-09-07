# Distributed File Storage Server API Documentation

## Upload File
- **URL**: `/upload`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file`: The file to upload
- **Response**: 
  - Success: `{"file_id": "<unique_file_id>"}`
  - Error: `{"error": "<error_message>"}`

## Get Files
- **URL**: `/files`
- **Method**: `GET`
- **Response**: 
  - Success: `{"files": ["<file_id1>", "<file_id2>", ...]}`
  - Error: `{"error": "<error_message>"}`

## Download File
- **URL**: `/download/:id`
- **Method**: `GET`
- **Parameters**:
  - `id`: The unique file ID
- **Response**: 
  - Success: File download starts
  - Error: `{"error": "<error_message>"}`

## Sample Working Video
[[Upload, Download and Get Files Api sample working video]](https://jmp.sh/s/3fXEtGIVtcGC3LwMgMTb)