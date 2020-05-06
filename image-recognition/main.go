package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/nikstoyanov/image-recognition/probability"
	"github.com/nikstoyanov/image-recognition/sntracker"
	"github.com/nikstoyanov/image-recognition/utils"
	"github.com/rs/cors"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// ClassifyResult stores the results from a classification.
// The store is given for a file name and an array of labels.
type ClassifyResult struct {
	Filename string                    `json:"filename"`
	Labels   []probability.LabelResult `json:"labels"`
}

var (
	graphModel   *tf.Graph
	sessionModel *tf.Session
	labels       []string
	tracker      *sp.Tracker
)

func main() {
	// Initiate tracking
	tracker = sntracker.SetupTracker()

	// Initiate Inception model
	if err := LoadModel(); err != nil {
		log.Fatal(err)
	}

	r := httprouter.New()

	r.POST("/recognize", recognizeHandler)
	handler := cors.Default().Handler(r)

	fmt.Println("Listening on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", handler))
}

// loadModel initiates the Inception CNN and the labels
func LoadModel() error {
	// Load inception model
	model, err := ioutil.ReadFile("./model/tensorflow_inception_graph.pb")

	if err != nil {
		return err
	}

	graphModel = tf.NewGraph()
	if err := graphModel.Import(model, ""); err != nil {
		return err
	}

	sessionModel, err = tf.NewSession(graphModel, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load labels
	labelsFile, err := os.Open("./model/imagenet_comp_graph_label_strings.txt")

	if err != nil {
		return err
	}

	// Close label file
	defer func() {
		err := labelsFile.Close()

		if err != nil {
			panic(err)
		}
	}()

	// Append labels
	scanner := bufio.NewScanner(labelsFile)
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// recognizeHandler receive a POST request from the frontend and passes it to the utilities
// module to scale the image and pass it to the Inception CNN
func recognizeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read image
	imageFile, header, err := r.FormFile("file")

	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		err := imageFile.Close()

		if err != nil {
			panic(err)
		}
	}()

	// Check the type of the uploaded file
	imageName := strings.Split(header.Filename, ".")

	// Record the request for a new image
	sntracker.TrackNewImage(tracker, imageName[0])

	if imageName[1] != "jpg" && imageName[1] != "png" {
		utils.ResponseError(w, "The file must be JPG or PNG!", http.StatusBadGateway)
		sntracker.TrackWrongExt(tracker, imageName[1])
		return
	}

	var imageBuffer bytes.Buffer

	// Copy image data to a buffer
	_, err = io.Copy(&imageBuffer, imageFile)
	if err != nil {
		panic(err)
	}

	// Make tensor
	tensor, err := utils.MakeTensorFromImage(&imageBuffer, imageName[:1][0])
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Run inference
	output, err := sessionModel.Run(
		map[tf.Output]*tf.Tensor{
			graphModel.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graphModel.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		utils.ResponseError(w, "Could not run inference", http.StatusInternalServerError)
		return
	}

	// Return best labels
	bestLabels := probability.FindBestLabels(output[0].Value().([][]float32)[0], labels)
	sntracker.TrackLabels(tracker, bestLabels)

	utils.ResponseJSON(w, ClassifyResult{
		Filename: header.Filename,
		Labels:   bestLabels,
	})
}
