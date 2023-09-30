package simulation

import "fmt"

func CreateResultReport(res SimResult, algName string) {
	res.algName = algName
	fmt.Printf("Results for algorithm: %s\n\n%+v", algName, res)
}
