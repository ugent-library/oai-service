package oaipmh

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	xmlnsXsi          = "http://www.w3.org/2001/XMLSchema-instance"
	xsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd"
)

var (
	OAIDC = MetadataFormat{
		MetadataPrefix:    "oai_dc",
		Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
	}

	verbs                   = map[string]struct{}{"Identify": {}, "ListMetadataFormats": {}, "ListSets": {}, "ListIdentifiers": {}, "ListRecords": {}, "GetRecord": {}}
	identifyArgs            = map[string]struct{}{"verb": {}}
	listMetadataFormatsArgs = map[string]struct{}{"verb": {}, "identifier": {}}
	listSetsArgs            = map[string]struct{}{"verb": {}, "resumptionToken": {}}
	listRecordsArgs         = map[string]struct{}{"verb": {}, "resumptionToken": {}, "metadataPrefix": {}, "set": {}, "from": {}, "until": {}}
	getRecordArgs           = map[string]struct{}{"verb": {}, "metadataPrefix": {}, "identifier": {}}

	// TODO not all errors should be public
	ErrVerbMissing              = &Error{Code: "badVerb", Value: "verb is missing"}
	ErrVerbRepeated             = &Error{Code: "badVerb", Value: "verb can't be repeated"}
	ErrVerbInvalid              = &Error{Code: "badVerb", Value: "verb is invalid"}
	ErrNoSetHierarchy           = &Error{Code: "noSetHierarchy", Value: "sets are not supported"}
	ErrIDDoesNotExist           = &Error{Code: "idDoesNotExist", Value: "identifier is unknown or illegal"}
	ErrNoRecordsMatch           = &Error{Code: "noRecordsMatch", Value: "no records match"}
	ErrNoMetadataFormats        = &Error{Code: "noMetadataFormats", Value: "there are no metadata formats available for the specified item"}
	ErrCannotDisseminateFormat  = &Error{Code: "cannotDisseminateFormat", Value: "the metadata format identified by the value given for the metadataPrefix argument is not supported by the item or by the repository"}
	ErrBadResumptionToken       = &Error{Code: "badResumptionToken", Value: "the value of the resumptionToken argument is invalid or expired"}
	ErrResumptiontokenExclusive = &Error{Code: "badArgument", Value: "resumptionToken cannot be combined with other attributes"}
	ErrMetadataPrefixMissing    = &Error{Code: "badArgument", Value: "metadataPrefix is missing"}
	ErrIdentifierMissing        = &Error{Code: "badArgument", Value: "identifier is missing"}
	ErrFromInvalid              = &Error{Code: "badArgument", Value: "from is not a valid datestamp"}
	ErrUntilInvalid             = &Error{Code: "badArgument", Value: "until is not a valid datestamp"}
	ErrSetDoesNotExist          = &Error{Code: "badArgument", Value: "set is unknown"}
)

type Request struct {
	XMLName         xml.Name `xml:"request"`
	URL             string   `xml:",chardata"`
	Verb            string   `xml:"verb,attr,omitempty"`
	MetadataPrefix  string   `xml:"metadataPrefix,attr,omitempty"`
	Identifier      string   `xml:"identifier,attr,omitempty"`
	Set             string   `xml:"set,attr,omitempty"`
	From            string   `xml:"from,attr,omitempty"`
	Until           string   `xml:"until,attr,omitempty"`
	ResumptionToken string   `xml:"resumptionToken,attr,omitempty"`
}

type response struct {
	provider          *Provider
	XMLName           xml.Name `xml:"http://www.openarchives.org/OAI/2.0/ OAI-PMH"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	ResponseDate      string   `xml:"responseDate"`
	Request           Request
	Errors            []*Error
	Body              any
}

type Error struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Value   string   `xml:",chardata"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Value)
}

type Identify struct {
	XMLName           xml.Name `xml:"Identify"`
	RepositoryName    string   `xml:"repositoryName"`
	BaseURL           string   `xml:"baseURL"`
	ProtocolVersion   string   `xml:"protocolVersion"`
	AdminEmail        []string `xml:"adminEmail"`
	Granularity       string   `xml:"granularity"`
	EarliestDatestamp string   `xml:"earliestDatestamp"`
	Compression       string   `xml:"compression,omitempty"`
	DeletedRecord     string   `xml:"deletedRecord"`
}

type ListMetadataFormats struct {
	XMLName         xml.Name          `xml:"ListMetadataFormats"`
	MetadataFormats []*MetadataFormat `xml:"metadataFormat"`
}

type ListSets struct {
	XMLName xml.Name `xml:"ListSets"`
	Sets    []*Set   `xml:"set"`
}

type GetRecord struct {
	XMLName xml.Name `xml:"GetRecord"`
	Record  *Record  `xml:"record"`
}

type ListIdentifiers struct {
	XMLName         xml.Name         `xml:"ListIdentifiers"`
	Headers         []*Header        `xml:"header"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
}

type ListRecords struct {
	XMLName         xml.Name         `xml:"ListRecords"`
	Records         []*Record        `xml:"record"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
}

type MetadataFormat struct {
	MetadataPrefix    string `xml:"metadataPrefix"`
	Schema            string `xml:"schema"`
	MetadataNamespace string `xml:"metadataNamespace"`
}

type Set struct {
	Spec        string   `xml:"setSpec"`
	Name        string   `xml:"setName"`
	Description *Payload `xml:"setDescription"`
}

type Header struct {
	Status     string   `xml:"status,attr,omitempty"`
	Identifier string   `xml:"identifier"`
	Datestamp  string   `xml:"datestamp"`
	SetSpec    []string `xml:"setSpec"`
}

type Payload struct {
	XML string `xml:",innerxml"`
}

type Record struct {
	Header   *Header  `xml:"header"`
	Metadata *Payload `xml:"metadata"`
}

type ResumptionToken struct {
	ExpirationDate   string `xml:"expirationDate,attr,omitempty"`
	CompleteListSize int    `xml:"completeListSize,attr,omitempty"`
	Cursor           int    `xml:"cursor,attr,omitempty"`
	Value            string `xml:",chardata"`
}

type Provider struct {
	ProviderConfig
	dateFormat string
	setMap     map[string]struct{}
}

// TODO use context in callbacks
type ProviderConfig struct {
	ErrorHandler        func(error)
	RepositoryName      string
	BaseURL             string
	AdminEmail          []string
	Granularity         string
	Compression         string
	DeletedRecord       string
	Sets                []*Set // TODO callback
	EarliestDatestamp   func() (time.Time, error)
	ListMetadataFormats func(*Request) ([]*MetadataFormat, error)
	GetRecord           func(*Request) (*Record, error)
	ListIdentifiers     func(*Request) ([]*Header, *ResumptionToken, error)
	ListRecords         func(*Request) ([]*Record, *ResumptionToken, error)
}

func NewProvider(conf ProviderConfig) (*Provider, error) {
	p := &Provider{
		ProviderConfig: conf,
		setMap:         make(map[string]struct{}),
	}

	if p.Granularity == "" {
		p.Granularity = "YYYY-MM-DDThh:mm:ssZ"
	}

	if p.Granularity == "YYYY-MM-DD" {
		p.dateFormat = "2006-01-02"
	} else if p.Granularity == "YYYY-MM-DDThh:mm:ssZ" {
		p.dateFormat = "2006-01-02T15:04:05Z"
	} else {
		return nil, errors.New("OAI-PMH granularity should be YYYY-MM-DD or YYYY-MM-DDThh:mm:ssZ")
	}

	if p.DeletedRecord == "" {
		p.DeletedRecord = "persistent"
	}

	for _, set := range p.Sets {
		p.setMap[set.Name] = struct{}{}
	}

	return p, nil
}

// TODO description, earliestDatestamp
func (p *Provider) identify(r *response) error {
	r.Body = &Identify{
		RepositoryName:  p.RepositoryName,
		BaseURL:         p.BaseURL,
		ProtocolVersion: "2.0",
		AdminEmail:      p.AdminEmail,
		Granularity:     p.Granularity,
		Compression:     p.Compression,
		DeletedRecord:   p.DeletedRecord,
	}
	return nil
}

func (p *Provider) listMetadataFormats(r *response) error {
	formats, err := p.ListMetadataFormats(&r.Request)
	if err == ErrIDDoesNotExist {
		r.Errors = append(r.Errors, err.(*Error))
		return nil
	}
	if err != nil {
		return err
	}
	if len(formats) == 0 {
		r.Errors = append(r.Errors, ErrNoMetadataFormats)
		return nil
	}
	r.Body = &ListMetadataFormats{
		MetadataFormats: formats,
	}
	return nil
}

// TODO resumptionToken, badResumptionToken
func (p *Provider) listSets(r *response) error {
	if len(p.Sets) == 0 {
		r.Errors = append(r.Errors, ErrNoSetHierarchy)
		return nil
	}
	r.Body = &ListSets{
		Sets: p.Sets,
	}
	return nil
}

func (p *Provider) listIdentifiers(r *response) error {
	headers, token, err := p.ListIdentifiers(&r.Request)
	if err == ErrBadResumptionToken || err == ErrCannotDisseminateFormat {
		r.Errors = append(r.Errors, err.(*Error))
		return nil
	}
	if err != nil {
		return err
	}
	if len(headers) == 0 {
		r.Errors = append(r.Errors, ErrNoRecordsMatch)
		return nil
	}
	r.Body = &ListIdentifiers{
		Headers:         headers,
		ResumptionToken: token,
	}
	return nil
}

func (p *Provider) listRecords(r *response) error {
	recs, token, err := p.ListRecords(&r.Request)
	if err == ErrBadResumptionToken || err == ErrCannotDisseminateFormat {
		r.Errors = append(r.Errors, err.(*Error))
		return nil
	}
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		r.Errors = append(r.Errors, ErrNoRecordsMatch)
		return nil
	}
	r.Body = &ListRecords{
		Records:         recs,
		ResumptionToken: token,
	}
	return nil
}

func (p *Provider) getRecord(r *response) error {
	rec, err := p.GetRecord(&r.Request)
	if err == ErrIDDoesNotExist || err == ErrCannotDisseminateFormat {
		r.Errors = append(r.Errors, err.(*Error))
		return nil
	}
	if err != nil {
		return err
	}
	r.Body = &GetRecord{
		Record: rec,
	}
	return nil
}

func (p *Provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(p.BaseURL)
	u.RawQuery = r.URL.RawQuery

	res := &response{
		provider:          p,
		XmlnsXsi:          xmlnsXsi,
		XsiSchemaLocation: xsiSchemaLocation,
		ResponseDate:      time.Now().UTC().Format(time.RFC3339),
		Request: Request{
			URL: u.String(),
		},
	}

	args := r.URL.Query()

	res.setVerb(args)

	var handler func(*response) error

	switch res.Request.Verb {
	case "Identify":
		handler = p.identify
		res.validateArgs(args, identifyArgs)
	case "ListMetadataFormats":
		handler = p.listMetadataFormats
		res.validateArgs(args, listMetadataFormatsArgs)
		res.setIdentifier(args)
	case "ListSets":
		handler = p.listSets
		res.validateArgs(args, listSetsArgs)
		res.setResumptionToken(args)
	case "ListIdentifiers":
		handler = p.listIdentifiers
		res.validateArgs(args, listRecordsArgs)
		res.setResumptionToken(args)
		res.setRequiredMetadataPrefix(args)
		res.setSet(args)
		res.setFromUntil(args)
	case "ListRecords":
		handler = p.listRecords
		res.validateArgs(args, listRecordsArgs)
		res.setResumptionToken(args)
		res.setRequiredMetadataPrefix(args)
		res.setSet(args)
		res.setFromUntil(args)
	case "GetRecord":
		handler = p.getRecord
		res.validateArgs(args, getRecordArgs)
		res.setRequiredMetadataPrefix(args)
		res.setRequiredIdentifier(args)
	}

	status := 200

	if len(res.Errors) == 0 {
		if err := handler(res); err != nil {
			p.handleError(w, err)
			return
		}
	} else {
		status = 400
	}

	out, err := xml.MarshalIndent(res, "", "  ")
	if err != nil {
		p.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	w.Write(out)
}

func (p *Provider) handleError(w http.ResponseWriter, err error) {
	if p.ErrorHandler != nil {
		p.ErrorHandler(err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (r *response) validateArgs(q url.Values, attrs map[string]struct{}) {
	for attr := range q {
		if _, ok := attrs[attr]; !ok {
			r.Errors = append(r.Errors, &Error{Code: "badArgument", Value: fmt.Sprintf("argument '%s' is illegal", attr)})
		}
	}
}

func (r *response) setVerb(q url.Values) {
	vals, ok := q["verb"]

	if !ok {
		r.Errors = append(r.Errors, ErrVerbMissing)
		return
	}
	if len(vals) > 1 {
		r.Errors = append(r.Errors, ErrVerbRepeated)
		return
	}
	if _, ok := verbs[vals[0]]; !ok {
		r.Errors = append(r.Errors, ErrVerbInvalid)
		return
	}

	r.Request.Verb = vals[0]
}

func (r *response) setResumptionToken(q url.Values) {
	r.Request.ResumptionToken = r.getArg(q, "resumptionToken")
}

func (r *response) setRequiredMetadataPrefix(q url.Values) {
	val := r.getArg(q, "metadataPrefix")

	if val != "" && r.Request.ResumptionToken != "" {
		r.Errors = append(r.Errors, ErrResumptiontokenExclusive)
		return
	}

	if val == "" && r.Request.ResumptionToken == "" {
		r.Errors = append(r.Errors, ErrMetadataPrefixMissing)
		return
	}

	r.Request.MetadataPrefix = val
}

func (r *response) setIdentifier(q url.Values) {
	r.Request.Identifier = r.getArg(q, "identifier")
}

func (r *response) setRequiredIdentifier(q url.Values) {
	r.Request.Identifier = r.getArg(q, "identifier")
	if r.Request.Identifier == "" {
		r.Errors = append(r.Errors, ErrIdentifierMissing)
	}
}

func (r *response) setSet(q url.Values) {
	val := r.getArg(q, "set")

	if val != "" && r.Request.ResumptionToken != "" {
		r.Errors = append(r.Errors, ErrResumptiontokenExclusive)
		return
	}

	if val != "" && len(r.provider.Sets) == 0 {
		r.Errors = append(r.Errors, ErrNoSetHierarchy)
		return
	}

	if _, ok := r.provider.setMap[val]; !ok {
		r.Errors = append(r.Errors, ErrSetDoesNotExist)
		return
	}

	r.Request.Set = val
}

func (r *response) setFromUntil(q url.Values) {
	f := r.getArg(q, "from")
	u := r.getArg(q, "until")
	if f != "" {
		if _, err := time.Parse(r.provider.dateFormat, f); err == nil {
			r.Request.From = f
		} else {
			r.Errors = append(r.Errors, ErrFromInvalid)
		}
	}
	if u != "" {
		if _, err := time.Parse(r.provider.dateFormat, u); err == nil {
			r.Request.From = f
		} else {
			r.Errors = append(r.Errors, ErrUntilInvalid)
		}
	}
}

func (r *response) getArg(q url.Values, attr string) string {
	vals, ok := q[attr]
	if !ok {
		return ""
	}

	if len(vals) > 1 {
		err := &Error{Code: "badArgument", Value: fmt.Sprintf("%s can't be repeated", attr)}
		r.Errors = append(r.Errors, err)
		return ""
	}

	if vals[0] == "" {
		err := &Error{Code: "badArgument", Value: fmt.Sprintf("%s is missing", attr)}
		r.Errors = append(r.Errors, err)
		return ""
	}

	return vals[0]
}
