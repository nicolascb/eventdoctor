package eventdoctor

import (
	"errors"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

var (
	validate    *validator.Validate
	semverRegex = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
)

var (
	ErrNoOwnerForTopic = errors.New("deve existir um produtor owner para cada tópico")
)

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("semver", validateSemver)
}

func validateSemver(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return semverRegex.MatchString(value)
}

func LoadSpecFromReader(f io.Reader) (EventDoctorSpec, error) {
	var cfg EventDoctorSpec
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return EventDoctorSpec{}, err
	}

	return cfg, nil
}
