package config

import (
	"encoding/json"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func write(cfg Config) error {

	//Marshalling the config struct to bytes
	j, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	//Call the getConfigFilePath to get the config path name
	writeFilename, err := getConfigFilePath()
	if err != nil {
		return err
	}

	//Write the byte to the config file
	file, err := os.Create(writeFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(j)

	return nil
}

func getConfigFilePath() (string, error) {
	homeName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := homeName + "/" + configFileName

	return filePath, nil
}

func Read() (Config, error) {
	fileName, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	fileData, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}

	defer fileData.Close()

	data, err := io.ReadAll(fileData)
	if err != nil {
		return Config{}, err
	}

	var configFile Config
	if err = json.Unmarshal(data, &configFile); err != nil {
		return Config{}, err
	}

	return configFile, nil

}

func (cfg Config) SetUser(username string) error {
	cfg.Current_user_name = username
	return write(cfg)
}
