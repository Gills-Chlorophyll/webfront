.DEFAULT_GOAL=run
BINARY_NAME=webfront
export S3_ACCESSKEY=11TEHKAX6YEG43D9RG5B
export S3_SECRETKEY=DFue3V5LKGmlpn7BShN3VZkXX5up4YyYHsAW2bIP
export S3_ENDPOINT=https://blr1.vultrobjects.com
export S3_REGION=IN
build:
	go build -o ./${BINARY_NAME}

run: build
	./${BINARY_NAME}

