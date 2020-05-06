setup-cnn:
	mkdir -p ./image-recognition/model && wget "https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip" -O ./image-recognition/model/inception.zip && unzip ./image-recognition/model/inception.zip -d ./image-recognition/model && chmod -R 777 ./image-recognition/model
setup-tf:
	curl -L "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz" | tar -C "/usr/local" -xz
	ldconfig
