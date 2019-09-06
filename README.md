This is a simple http server that accepts an image, converts it to a different format, resizes/compresses it and sends it back.

## Run on CloudRun

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/kkty/image-server)

## Usage

```console
$ docker pull kkty/image-server
$ docker run -d -p 8080:8080 kkty/image-server
```

```console
$ curl "http://localhost:8080?quality=30&width=100&height=200" \
  -X POST \
  --data-binary '@original.png' \
  -H 'content-type:image/png' \
  -H 'accept:image/jpeg' \
  > compressed.jpg
```

- `content-type` / `accept` can either be `image/png`, `image/jpeg` or `image/gif`.
- `quality` can be an integer value ranging from 1 to 100.
  - It takes effect only when the output is in jpeg format, i.e. `accept` is set to `image/jpeg`.
- `width` and `height` set the size of the output image in pixels.
  - If `width` is set and `height` is not set, `height` will be such that the original aspect ratio will be kept, and vice versa.
  - If both of them are not set, the original size will be kept unchanged.

## Notes

- The port 8080 is used as default. To change the port to listen on, set `PORT` environment variable.
