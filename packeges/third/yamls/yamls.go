package yamls

import (
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func Read(fpath string)  {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Printf("ReadFile: %v\n", err)
		return
	}
	cm := &corev1.ConfigMap{}
	err = yaml.Unmarshal(data, cm)
	if err != nil {
		fmt.Printf("Unmarshal: %v\n", err)
		return
	}
	fmt.Println(cm.Name)
	fmt.Printf("%+v\n", cm.Data)
}