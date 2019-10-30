# downloadmanager-versions-api

## What is it?
This project is a simple REST-based JSON API to allow on-client software to detect updates that have been made upstream

## How does a client use it?

Two endpoints are provided - one to receive updates of new versions from a build process and another for client software to "call back".

### /newversion
This is protected by an API Key (see Deployment section) and expects a POST request with a
JSON request body in the following format:
```json
{
  "event": "newversion",
  "buildId": 12345,
  "branch": "someBranch",
  "productName": "myProductName",
  "downloadUrl": "https://download-server.domain.com/path/to/download"
}
```

- `event` - must be the string "newversion"
- `buildId` - the numeric identifier for this build. A higher build number is assumed to be later.
- `branch` - the branch that this build is from
- `productName` - unique name for this software product, allows different software products to be queried
- `downloadUrl` - location that this software can be automatically downloaded from

The json object is defined in `lambdas/common/models.go`

An HTTP 201 response (Created) is returned with an empty response body.  If an error occurs, a text/plain response body is sent.

### /lookup
This is an open endpoint and expects a GET request with a JSON request body in the following format:
```json
{
  "branch": "someBranch",
  "productName": "myProductName",
  "alwaysShowMaster": false
}
```

- `branch` - look for versions from this branch. Set to `master` if branch is not relevant.
- `productName` - name of the software product to look for. Must match `productName` from the build process.
- `alwaysShowMaster` - if looking for a branch, also show the latest `master` branch build

The json object is defined in `lambdas/common/models.go`

It is assumed that a piece of client software will know what branch and productName it was built from and makes a request
at startup.  The endpoint returns a JSON array of the latest release for the provided branch and optionally for the master
branch as well, using the same record format as for the `/newversion` endpoint

## How do I deploy it?

The project deploys using AWS API Gateway and needs a few steps to build:

1. Clone and compile the code
    - You'll need Go 1.11 or higher to compile the code, as well as GNU `make` and `zip`.
    - A Docker image with the prerequsites can be used by running `docker run --rm -it guardianmultimedia/gobuild:1 /bin/sh`
    - Simply run `make test` to run the tests and `make deployables` to build Lambda-compatible zips in the lambdas/deployables directory.
    
2. Upload the Zip files
    - Create an S3 bucket to hold the code or use an existing one
    - The Cloudformation expects the zip files to live in a subpath `application/stack/stage/{zipfile}` where `application`,
    `stack` and `stage` are all parameters that you specify when you build the cloudformation.  Decide what you want them to
    be and upload the Zips to this path (`stage` must be either CODE or PROD).
    
3. Deploy the cloudformatiopn
    - This can be found in the `cloudformation/` directory of the source code. You'll need to provide the name of the bucket
    that you used in step 2 as well as the `application`, `stack` and `stage` parameters.
    - The cloudformation will fail to deploy if Lambda can't find the zips you uploaded in stage two. If this happens,
    check that you have correctly set the path in stage 2 to your values for `application`, `stack` and `stage`.
    
4. Create an API key
    - In the AWS Console, go to API Gateway and select "API Keys".
    - Create an auto-generated key
    - In "API Stage Association", select "Versions API" in the "Select API" box and select the relevant stage
    - NOW CLICK ADD. **It won't work** if you dont.
    
5. Find the deployment URL by clicking APIs -> Versions API -> Stages -> {PROD or CODE}.  In the panel you see, there should
be a blue box at the top with a section called "Invoke URL".  Use this as the base URL to talk to.

6. If you make a GET request as above you should see a 404 response (we have no data yet).

7. If you make a POST request as above you should see a 403 Forbidden response.  Put the API key that you created in step 4
into a request header called `x-api-key` and you should be able to set data

8. Once you have put a record into the table retry the GET request. You should be able to see the record for the most recent
build ID in your response.

9. You can now set up your client software and build process to allow automatic updates