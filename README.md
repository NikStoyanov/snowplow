# Image recognition API and Snowplow tracker

## Description
Uses the Inception5h convolutional neural network to classify and image and
Snowplow micro to track new images, faulty extensions and the top 5 probability
labels. The service is written in Go 1.12 and uses Tensorflow 1.12.

## Deployment
Requires golang1.12 and a *nix environment which supports Tensorflow 1.12.

To setup tensorflow for golang run:
```
make setup-tf
```

To setup the convolutional neural network Inception5h run:
```
make setup-cnn
```

Enter the image-recognition directory and run the program:
```
cd image-recognition
make
```

## Tracking
The primary purpose of tracking with Snowplow is to detect the top 5
probability labels which result from an image classification. By combining this
with the IP address of the client a geographical estimate can be made of the
interests of users of the service.

Snowplow micro is also used to track the creation of a new image recognition
process. The tracking is done initially to store the type of valid and invalid
extensions. The information can later be used to optimize the usage for certain
extensions and measure the faulty requests. If the number of faulty requests is
determined to be too high then this backend service will be made accesible only
through a frontend which can use JavaScript to vet the extension type, thereby,
reducing the cost of running the server.

## Usage
To classify and image pass a POST request with:
```
curl localhost:8090/recognize -F 'file=@./cat.jpg'
```

To track the stored events use httpie to provide a nice layout and run:
```
http localhost:9090/micro/good
```

## Docker compose bug
First attempt used the existing [source code](https://github.com/NikStoyanov/image-recognition) for the image
recognition which is developed using Docker. When configured the code would run
with Snowplow micro, however, an error message of *-1* was returned for the
status code. This points to
https://github.com/snowplow/snowplow-golang-tracker/blob/50e2f03c3f5c2afb6b3b252d49ab75af987e1b52/tracker/emitter.go#L358
not returning a valid error code.
