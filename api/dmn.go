package main

import (
	`fmt`
	`net/http`
	`strings`
	`github.com/jscherff/dmnsdk/model`
)
// =============================================================================
// Get List Parameters
//
//	Name			Description
//	----			-----------
//	decisionDefinitionId	Filter by decision definition id.
//	decisionDefinitionIdIn	Filter by decision definition ids.
//	name			Filter by decision definition name.
//	nameLike		Filter by decision definition name substring.
//				parameter is a substring of.
//	deploymentId		Filter by the deployment the id belongs to.
//	key			Filter by decision definition key.
//	keyLike			Filter by decision definition key substring.
//	category		Filter by decision definition category.
//	categoryLike		Filter by decision definition category substring.
//	version			Filter by decision definition version.
//	latestVersion		Only include those decision definitions that
//				are latest versions. Value may only be true,
//				as false is the default behavior.
//	resourceName		Filter by decision definition resource.
//	resourceNameLike	Filter by decision definition resource substring.
//	versionTag		Filter by the version tag.
//	sortBy			Sort the results by a given criterion. Valid
//				values are category, key, id, name, version,
//				deploymentId and tenantId. Must be used in
//				conjunction with the sortOrder parameter.
//	sortOrder		Sort the results in a given order, asc or desc.
//	firstResult		Index of the first result to return.
//	maxResults		Maximum number of results to return.
//
// Get Count Parameters
//
//	Name			Description
//	----			-----------
//	decisionDefinitionId	Filter by decision definition id.
//	decisionDefinitionIdIn	Filter by decision definition ids.
//	name			Filter by decision definition name.
//	nameLike		Filter by decision definition name substring.
//				parameter is a substring of.
//	deploymentId		Filter by the deployment the id belongs to.
//	key			Filter by decision definition key.
//	keyLike			Filter by decision definition key substring.
//	category		Filter by decision definition category.
//	categoryLike		Filter by decision definition category substring.
//	version			Filter by decision definition version.
//	latestVersion		Only include those decision definitions that
//				are latest versions. Value may only be true,
//				as false is the default behavior.
//	resourceName		Filter by decision definition resource.
//	resourceNameLike	Filter by decision definition resource substring.
//	versionTag		Filter by the version tag.
// ----------------------------------------------------------------------------- 

type Endpoint string

const (
	epDefinitionPrefix string	= `/engine-rest/decision-definition`
	epDefinitionList Endpoint	= ``
	epDefinitionCount Endpoint	= `/count`
	epDefinitionInfoById Endpoint	= `/%s`
	epDefinitionInfoByKey Endpoint	= `/key/%s`
	epDefinitionById Endpoint	= `/%s`
	epDefinitionByKey Endpoint	= `/key/%s`
	epDmnXmlById Endpoint		= `/%s/xml`
	epDmnXmlByKey Endpoint		= `/key/%s/xml`
)

func (this Endpoint) String() (string) {
	return string(this)
}

func (this Endpoint) With(param string) (string) {
	return fmt.Sprintf(string(this), param)
}

type DmnApi interface {
	DefinitionList() (*model.DefinitionList, error)
}

type dmnApi struct {
	Server string
}

func NewDmnApi(server string) (DmnApi) {
	return &dmnApi{server + epDefinitionPrefix}
}

func (this *dmnApi) GetDefinitionList() (*model.DefinitionList, error) {

	dl := &model.DefinitionList{}
	url := this.Server + epDefinitionList.String()

	if err := dl.LoadUrl(url); err != nil {
		return nil, err
	}

	return dl, nil
}

func (this *dmnApi) GetDefinitionById(id string) (*model.Definition, error) {

	d := &model.Definition{}
	url := this.Server + epDefinitionById.With(id)

	if err := d.LoadUrl(url); err != nil {
		return nil, err
	}

	return d, nil
}

func (this *dmnApi) GetDefinitionById(id string) (*model.Definition, error) {

	d := &model.Definition{}
	url := this.Server + epDefinitionById.With(id)

	if err := d.LoadUrl(url); err != nil {
		return nil, err
	}

	return d, nil
}





func main() {
}

// withParam substitutes a path parameter with the provided value.
func withParam(endpoint, param, value string) (string) {
}

// loadJson unmarshals JSON from an io.Reader into an object.
func loadJson(t interface{}, r io.Reader) (error) {
	return json.NewDecoder(r).Decode(&t)
}

// loadXml unmarshals XML from an io.Reader into an object.
func loadXml(t interface{}, r io.Reader) (error) {
	return xml.NewDecoder(r).Decode(&t)
}

// loadJsonFromUrl unmarshals JSON from a URL into an object.
func loadJsonFromUrl(t interface{}, u string) (error) {

	if resp, err := http.Get(u); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		return loadJson(t, resp.Body)
	}
}

// loadJsonFromFile unmarshals JSON from a file into an object.
func loadJsonFromFile(t interface{}, f string) (error) {

        if fh, err := os.Open(f); err != nil {
                return err
	} else {
                defer fh.Close()
		return loadJson(t, fh)
        }
}
