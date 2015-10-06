# s3zipper
Microservice that Servers Streaming Zip file from S3 Securely

## Read the blog here
[Original Blog Post](http://engineroom.teamwork.com/how-to-securely-provide-a-zip-download-of-a-s3-file-bundle/)

## Sample redis console call:

```
set zip:hello '[{"S3Path":"path/to/image.png", "FileName":"image_name.jpg", "Folder":"images"}]' EX 10
```
