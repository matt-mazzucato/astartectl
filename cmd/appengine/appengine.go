// Copyright © 2019 Ispirata Srl
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package appengine

import (
	"errors"

	"github.com/astarte-platform/astartectl/client"

	"github.com/astarte-platform/astartectl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppEngineCmd represents the appEngine command
var AppEngineCmd = &cobra.Command{
	Use:               "appengine",
	Short:             "Interact with AppEngine API",
	Long:              `Interact with AppEngine API.`,
	PersistentPreRunE: appEnginePersistentPreRunE,
}

var realm string
var appEngineJwt string
var realmManagementJwt string
var astarteAPIClient *client.Client

func init() {
	AppEngineCmd.PersistentFlags().StringP("realm-key", "k", "",
		"Path to realm private key used to generate JWT for authentication")
	AppEngineCmd.MarkPersistentFlagFilename("realm-key")
	AppEngineCmd.PersistentFlags().String("appengine-url", "",
		"AppEngine API base URL. Defaults to <astarte-url>/appengine.")
	viper.BindPFlag("appengine.url", AppEngineCmd.PersistentFlags().Lookup("appengine-url"))
	AppEngineCmd.PersistentFlags().String("realm-management-url", "",
		"Realm Management API base URL. Defaults to <astarte-url>/realmmanagement.")
	AppEngineCmd.PersistentFlags().StringP("realm-name", "r", "",
		"The name of the realm that will be queried")
}

func appEnginePersistentPreRunE(cmd *cobra.Command, args []string) error {
	appEngineURLOverride := viper.GetString("appengine.url")
	viper.BindPFlag("realm-management.url", cmd.Flags().Lookup("realm-management-url"))
	realmManagementURLOverride := viper.GetString("realm-management.url")
	astarteURL := viper.GetString("url")
	if appEngineURLOverride != "" {
		// Use explicit appengine-url
		var err error
		astarteAPIClient, err = client.NewClientWithIndividualURLs(appEngineURLOverride, "", "", realmManagementURLOverride, nil)
		if err != nil {
			return err
		}
	} else if astarteURL != "" {
		var err error
		astarteAPIClient, err = client.NewClient(astarteURL, nil)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Either astarte-url or appengine-url have to be specified")
	}

	viper.BindPFlag("realm.key", cmd.Flags().Lookup("realm-key"))
	appEngineKey := viper.GetString("realm.key")
	if appEngineKey == "" {
		return errors.New("realm-key is required")
	}

	viper.BindPFlag("realm.name", cmd.Flags().Lookup("realm-name"))
	realm = viper.GetString("realm.name")
	if realm == "" {
		return errors.New("realm is required")
	}

	var err error
	appEngineJwt, err = generateAppEngineJWT(appEngineKey)
	if err != nil {
		return err
	}
	realmManagementJwt, err = generateRealmManagementJWT(appEngineKey)
	if err != nil {
		return err
	}

	return nil
}

func generateAppEngineJWT(privateKey string) (jwtString string, err error) {
	return utils.GenerateAstarteJWTFromKeyFile(privateKey, utils.AppEngine, nil, 300)
}

func generateRealmManagementJWT(privateKey string) (jwtString string, err error) {
	return utils.GenerateAstarteJWTFromKeyFile(privateKey, utils.RealmManagement, nil, 300)
}
