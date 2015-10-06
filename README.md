# s3zipper
Microservice that Servers Streaming Zip file from S3 Securely

## Read the blog here
[Original Blog Post](http://engineroom.teamwork.com/how-to-securely-provide-a-zip-download-of-a-s3-file-bundle/)

## Sample redis console call:

```
set zip:hello '[{"S3Path":"photos/images/000/000/001/original/feb-14-prints-charming-cal-1440x900.png", "FileName":"sip1.jpg", "Folder":"event_photos"}, {"S3Path":"photos/images/000/000/002/original/feb-14-love-of-my-life-cal-1600x1200.jpg", "FileName":"sip2.jpg", "Folder":"event_photos"}]' EX 10
```
