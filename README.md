# ontlogin sample container

## what is
  This repository has been built to make it easy for anyone to try out ONTLOGIN.  
  If you have a Docker environment, you can try out the simplest implementation of QR code authentication in a few steps.  

## how to run

  running build.sh.  
  What we are doing in build.sh is as follows.  
  at the 1st, build backend script.  
  2nd, build frontend script.  
  use golang (1.17) container and node lts container(now version is 16).  

```bash
docker run --rm -w /go/src/app -v $PWD/backend:/go/src/app golang:1.17 bash -c "go build"
docker run --rm -w /opt/app -v $PWD/frontend:/opt/app node:lts-stretch bash -c "npm install ; npm run build"
```
  next step, run docker-compose.  

```bash
docker-compose up -d
```

  With the default settings, you can try ONT ID login using QR code by visiting http://localhost:8080.

## about frontend/backend
### frontend
  custom vue demo script in https://github.com/ontology-tech/ontlogin-sdk-js (commit 7874a0ac21ac31ff9bb8e23cd90b12cbec0a16a4)  
  patch ontlogin.js using (patch-package)[https://www.npmjs.com/package/patch-package]. (please show ./patches)  

### backend 
  using a modified of https://github.com/ontology-tech/ontlogin-sample-go  
