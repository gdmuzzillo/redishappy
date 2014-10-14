package configuration

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/mdevilliers/redishappy/types"
)

type Configuration struct {
	Clusters  []types.Cluster
	HAProxy   types.HAProxy
	Sentinels []types.Sentinel
}

func LoadFromFile(filePath string) (*Configuration, error) {

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &Configuration{}, err
	}
	return ParseConfiguration(content)
}

func ParseConfiguration(configurationAsJson []byte) (*Configuration, error) {

	configuration := new(Configuration)
	err := json.Unmarshal(configurationAsJson, &configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func (c *Configuration) SanityCheckConfiguration(tests ...SanityCheck) (bool, []string) {

	errorlist := []string{}
	isSane := true

	for _, test := range tests {

		ok, err := test.Check(c)
		if !ok {
			errorlist = append(errorlist, err.Error())
			isSane = false
		}
	}

	return isSane, errorlist
}

func (config *Configuration) FindClusterByName(name string) (*types.Cluster, error) {

	for _, cluster := range config.Clusters {
		if cluster.Name == name {
			return &cluster, nil
		}
	}
	return &types.Cluster{}, errors.New("Cluster not found")
}
