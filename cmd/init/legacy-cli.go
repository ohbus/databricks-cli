package init

import (
	"fmt"

	"github.com/databricks/bricks/project"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

func loadCliProfiles() (profiles []Answer, err error) {
	file, err := homedir.Expand("~/.databrickscfg")
	if err != nil {
		return
	}
	gitConfig, err := ini.Load(file)
	if err != nil {
		return
	}
	for _, v := range gitConfig.Sections() {
		host, err := v.GetKey("host")
		if err != nil {
			// invalid profile
			continue
		}
		profiles = append(profiles, Answer{
			Value:   v.Name(),
			Details: fmt.Sprintf(`Connecting to "%s" workspace`, host),
			Callback: func(ans Answer, prj *project.Project, _ Results) {
				prj.Profile = ans.Value
			},
		})
	}
	return
}

func getConnectionProfile() (*Choice, error) {
	profiles, err := loadCliProfiles()
	if err != nil {
		return nil, err
	}
	// TODO: propmt for password and create ~/.databrickscfg
	return &Choice{
		key:     "profile",
		Label:   "Databricks CLI profile",
		Answers: profiles,
	}, err
}