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
//	decisionDmnId	Filter by decision definition id.
//	decisionDmnIdIn	Filter by decision definition ids.
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
//	decisionDmnId	Filter by decision definition id.
//	decisionDmnIdIn	Filter by decision definition ids.
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
	epDmnPrefix string	= `/engine-rest/decision-definition`
	epDmnList Endpoint	= ``
	epDmnMap Endpoint	= ``
	epDmnCount Endpoint	= `/count`
	epDmnInfoById Endpoint	= `/%s`
	epDmnInfoByKey Endpoint	= `/key/%s`
	epDmnXmlById Endpoint	= `/%s/xml`
	epDmnXmlByKey Endpoint	= `/key/%s/xml`
	epDmnById Endpoint	= `/%s`
	epDmnByKey Endpoint	= `/key/%s`
)

func (this Endpoint) String() (string) {
	return string(this)
}

func (this Endpoint) With(param string) (string) {
	return fmt.Sprintf(string(this), param)
}

type DmnApi interface {
	DmnList() (*model.DmnList, error)
	DmnMap() (model.DmnMap, error)
	DmnInfoById(id string) (*model.DmnInfo, error)
	DmnInfoByKey(key string) (*model.DmnInfo, error)
	DmnInfoByKeyVer(key string, ver int) (*model.DmnInfo, error)
	DmnXmlById(id string) (*model.DmnXml, error)
	DmnXmlByKey(key string) (*model.DmnXml, error)
	DmnXmlByKeyVer(key string, ver int) (*model.DmnXml, error)
	DmnById(id string) (*model.Dmn, error)
	DmnByKey(key string) (*model.Dmn, error)
	DmnByKeyVer(key string, ver int) (*model.Dmn, error)
}

type dmnApi struct {
	Server string
	dmnList *model.DmnList
	dmnMap model.DmnMap
}

func NewDmnApi(server string) (DmnApi) {
	return &dmnApi{server + epDmnPrefix, nil, nil}
}

func (this *dmnApi) DmnList() (*model.DmnList, error) {

	if this.dmnList != nil {
		return this.dmnList, nil
	}

	url := this.Server + epDmnList.String()

	if dl, err := model.NewDmnList(url); err != nil {
		return nil, err
	} else {
		this.dmnList = dl
		return this.dmnList, nil
	}
}

func (this *dmnApi) DmnMap() (model.DmnMap, error) {

	if this.dmnMap != nil {
		return this.dmnMap, nil
	}

	if dl, err := this.DmnList(); err != nil {
		return nil, err
	} else if dm, err := dl.Map(); err != nil {
		return nil, err
	} else {
		this.dmnMap = dm
		return this.dmnMap, nil
	}
}

func (this *dmnApi) DmnInfoById(id string) (*model.DmnInfo, error) {
	url := this.Server + epDmnInfoById.With(id)
	return model.NewDmnInfo(url)
}

func (this *dmnApi) DmnInfoByKey(key string) (*model.DmnInfo, error) {
	url := this.Server + epDmnInfoByKey.With(key)
	return model.NewDmnInfo(url)
}

func (this *dmnApi) DmnInfoByKeyVer(key string, ver int) (*model.DmnInfo, error) {

	if dm, err := this.DmnMap(); err != nil {
		return nil, err
	} else if di, err := dm.DmnInfo(key, ver); err != nil {
		return nil, err
	} else {
		return di, nil
	}
}

func (this *dmnApi) DmnXmlById(id string) (*model.DmnXml, error) {
	url := this.Server + epDmnXmlById.With(id)
	return model.NewDmnXml(url)
}

func (this *dmnApi) DmnXmlByKey(key string) (*model.DmnXml, error) {
	url := this.Server + epDmnXmlByKey.With(key)
	return model.NewDmnXml(url)
}

func (this *dmnApi) DmnXmlByKeyVer(key string, ver int) (*model.DmnXml, error) {

	if di, err := this.DmnInfoByKeyVer(key, ver); err != nil {
		return nil, err
	} else if dx, err := this.DmnXmlById(di.Id); err != nil {
		return nil, err
	} else {
		return dx, nil
	}
}

func (this *dmnApi) DmnById(id string) (*model.Dmn, error) {

	if dx, err := this.DmnXmlById(id); err != nil {
		return nil, err
	} else {
		return model.NewDmn(dx.DmnXml)
	}
}

func (this *dmnApi) DmnByKey(key string) (*model.Dmn, error) {

	if dx, err := this.DmnXmlByKey(key); err != nil {
		return nil, err
	} else {
		return model.NewDmn(dx.DmnXml)
	}
}

func (this *dmnApi) DmnByKeyVer(key string, ver int) (*model.Dmn, error) {

	if dx, err := this.DmnXmlByKeyVer(key, ver); err != nil {
		return nil, err
	} else {
		return model.NewDmn(dx.DmnXml)
	}
}
