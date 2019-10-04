package probability

import "sort"

// LabelResult provides the propability each label
// represents the image.
type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

// ByProbability is a sorted list of probability and label after classification.
// Because of lack of generics the methods below need to be overloaded.
type ByProbability []LabelResult

func (a ByProbability) Len() int {
	return len(a)
}

func (a ByProbability) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByProbability) Less(i, j int) bool {
	return a[i].Probability > a[j].Probability
}

func FindBestLabels(probabilities []float32, labels []string) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Sort by probability
	sort.Sort(ByProbability(resultLabels))
	// Return top 5 labels
	return resultLabels[:5]
}
