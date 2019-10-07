# Image recognition API

[![Build Status](https://travis-ci.com/NikStoyanov/image-recognition.svg?branch=master)](https://travis-ci.com/NikStoyanov/image-recognition)

## Status
Uses Inception v5, further work will include a custom network for
fun :).

## Deployment
The CI builds the Docker image and deploys it to AWS ECR. Then AWS ECS takes the
new image and runs the container.

## CLI Use
`curl HOST:8090/recognize -F 'file=@./file.jpg'`

## Frontend
To use the React frontend go here (TBC).
