package api

import (
	`fmt`
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
	epDefinitionXmlById Endpoint	= `/%s/xml`
	epDefinitionXmlByKey Endpoint	= `/key/%s/xml`
	epDefinitionById Endpoint	= `/%s`
	epDefinitionByKey Endpoint	= `/key/%s`
)

func (this Endpoint) String() (string) {
	return string(this)
}

func (this Endpoint) With(param string) (string) {
	return fmt.Sprintf(string(this), param)
}

type DmnApi interface {
	GetDefinitionList() (*model.DefinitionList, error)
	GetDefinitionInfoById(string) (*model.DefinitionInfo, error)
	GetDefinitionInfoByKey(string) (*model.DefinitionInfo, error)
	GetDefinitionXmlById(string) (*model.DefinitionInfo, error)
	GetDefinitionXmlByKey(string) (*model.DefinitionInfo, error)
	GetDefinitionById(string) (*model.Definition, error)
	GetDefinitionByKey(string) (*model.Definition, error)
}

type dmnApi struct {
	Server string
}

func NewDmnApi(server string) (DmnApi) {
	return &dmnApi{server + epDefinitionPrefix}
}

func (this *dmnApi) GetDefinitionList() (*model.DefinitionList, error) {
	url := this.Server + epDefinitionList.String()
	return model.NewDefinitionList(url)
}

func (this *dmnApi) GetDefinitionInfoById(id string) (*model.DefinitionInfo, error) {
	url := this.Server + epDefinitionInfoById.With(id)
	return model.NewDefinitionInfo(url)
}

func (this *dmnApi) GetDefinitionInfoByKey(key string) (*model.DefinitionInfo, error) {
	url := this.Server + epDefinitionInfoByKey.With(key)
	return model.NewDefinitionInfo(url)
}

func (this *dmnApi) GetDefinitionXmlById(id string) (*model.DefinitionInfo, error) {
	url := this.Server + epDefinitionXmlById.With(id)
	return model.NewDefinitionInfo(url)
}

func (this *dmnApi) GetDefinitionXmlByKey(key string) (*model.DefinitionInfo, error) {
	url := this.Server + epDefinitionXmlByKey.With(key)
	return model.NewDefinitionInfo(url)
}


func (this *dmnApi) GetDefinitionById(id string) (*model.Definition, error) {
	if di, err := this.GetDefinitionXmlById(id); err != nil {
		return nil, err
	} else {
		return model.NewDefinition(di.DmnXml)
	}
}

func (this *dmnApi) GetDefinitionByKey(key string) (*model.Definition, error) {
	if di, err := this.GetDefinitionXmlByKey(key); err != nil {
		return nil, err
	} else {
		println(di.DmnXml)
		return model.NewDefinition(di.DmnXml)
	}
}
