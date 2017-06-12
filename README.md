# s3put

I know exists, I needed it anyway.

## Usage

It uses the standard s3cli authentication methods.

	s3put --bucket=<bucket> --region=<region> <file> <key>

You can use environment variables to provide secret key & access key:

	AWS_SECRET_KET=abc AWS_ACCESS_KEY=def s3put --bucket=<bucket> --region=<region> <file> <key>

Uses multipart upload to allow files bigger than 5GB.