package oaipmh

import (
	"context"
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

type handleFunc func(context.Context, *Provider, *response, url.Values) error

var (
	OAIDC = &MetadataFormat{
		MetadataPrefix:    "oai_dc",
		Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
	}

	verbHandlers = map[string][]handleFunc{
		"Identify": {
			allowArgs(),
			identify,
		},
		"ListMetadataFormats": {
			allowArgs("identifier"),
			setIdentifier,
			listMetadataFormats,
		},
		"ListSets": {
			allowArgs("resumptionToken"),
			setResumptionToken,
			listSets,
		},
		"ListIdentifiers": {
			allowArgs("resumptionToken", "metadataPrefix", "set", "from", "until"),
			setResumptionToken,
			setRequiredMetadataPrefix,
			setSet,
			setFromUntil,
			listIdentifiers,
		},
		"ListRecords": {
			allowArgs("resumptionToken", "metadataPrefix", "set", "from", "until"),
			setResumptionToken,
			setRequiredMetadataPrefix,
			setSet,
			setFromUntil,
			listRecords,
		},
		"GetRecord": {
			allowArgs("metadataPrefix", "identifier"),
			setRequiredMetadataPrefix,
			setRequiredIdentifier,
			getRecord,
		},
	}

	// TODO make all errors private and remove Error method
	ErrCannotDisseminateFormat = &Error{Code: "cannotDisseminateFormat", Value: "the metadata format identified by the value given for the metadataPrefix argument is not supported by the item or by the repository"}
	ErrBadResumptionToken      = &Error{Code: "badResumptionToken", Value: "the value of the resumptionToken argument is invalid or expired"}

	errVerbMissing              = &Error{Code: "badVerb", Value: "verb is missing"}
	errVerbRepeated             = &Error{Code: "badVerb", Value: "verb can't be repeated"}
	errVerbInvalid              = &Error{Code: "badVerb", Value: "verb is invalid"}
	errNoSetHierarchy           = &Error{Code: "noSetHierarchy", Value: "sets are not supported"}
	errIDDoesNotExist           = &Error{Code: "idDoesNotExist", Value: "identifier is unknown or illegal"}
	errNoRecordsMatch           = &Error{Code: "noRecordsMatch", Value: "no records match"}
	errNoMetadataFormats        = &Error{Code: "noMetadataFormats", Value: "there are no metadata formats available for the specified item"}
	errResumptiontokenExclusive = &Error{Code: "badArgument", Value: "resumptionToken cannot be combined with other attributes"}
	errMetadataPrefixMissing    = &Error{Code: "badArgument", Value: "metadataPrefix is missing"}
	errIdentifierMissing        = &Error{Code: "badArgument", Value: "identifier is missing"}
	errFromInvalid              = &Error{Code: "badArgument", Value: "from is not a valid datestamp"}
	errUntilInvalid             = &Error{Code: "badArgument", Value: "until is not a valid datestamp"}
	errSetDoesNotExist          = &Error{Code: "badArgument", Value: "set is unknown"}
)

type request struct {
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
	XMLName           xml.Name `xml:"http://www.openarchives.org/OAI/2.0/ OAI-PMH"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	ResponseDate      string   `xml:"responseDate"`
	Request           request
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
	AdminEmails       []string `xml:"adminEmail"`
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
	XMLName         xml.Name         `xml:"ListSets"`
	Sets            []*Set           `xml:"set"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
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
	SetSpec        string   `xml:"setSpec"`
	SetName        string   `xml:"setName"`
	SetDescription *Payload `xml:"setDescription"`
}

type Header struct {
	Status     string   `xml:"status,attr,omitempty"`
	Identifier string   `xml:"identifier"`
	Datestamp  string   `xml:"datestamp"`
	SetSpecs   []string `xml:"setSpec"`
}

type Payload struct {
	XML string `xml:",innerxml"`
}

type Record struct {
	Header   *Header  `xml:"header"`
	Metadata *Payload `xml:"metadata"`
}

type ResumptionToken struct {
	CompleteListSize int    `xml:"completeListSize,attr,omitempty"`
	Value            string `xml:",chardata"`
	ExpirationDate   string `xml:"expirationDate,attr,omitempty"`
	Cursor           *int   `xml:"cursor,attr,omitempty"`
}

type Provider struct {
	ProviderConfig
	dateFormat string
}

// TODO use context in callbacks
type ProviderConfig struct {
	RepositoryName string
	BaseURL        string
	AdminEmails    []string
	Granularity    string
	Compression    string
	DeletedRecord  string
	StyleSheet     string
	ErrorHandler   func(error)
	Backend        ProviderBackend
}

type ProviderBackend interface {
	GetEarliestDatestamp(context.Context) (time.Time, error)
	HasMetadataFormat(context.Context, string) (bool, error)
	HasSets(context.Context) (bool, error)
	HasSet(context.Context, string) (bool, error)
	GetMetadataFormats(context.Context) ([]*MetadataFormat, error)
	GetRecordMetadataFormats(context.Context, string) ([]*MetadataFormat, error)
	GetSets(context.Context) ([]*Set, *ResumptionToken, error)
	GetMoreSets(context.Context, string) ([]*Set, *ResumptionToken, error)
	// TODO pass from, until as time objects
	GetIdentifiers(context.Context, string, string, string, string) ([]*Header, *ResumptionToken, error)
	GetMoreIdentifiers(context.Context, string) ([]*Header, *ResumptionToken, error)
	HasRecord(context.Context, string) (bool, error)
	GetRecords(context.Context, string, string, string, string) ([]*Record, *ResumptionToken, error)
	GetMoreRecords(context.Context, string) ([]*Record, *ResumptionToken, error)
	GetRecord(context.Context, string, string) (*Record, error)
}

func NewProvider(conf ProviderConfig) (*Provider, error) {
	p := &Provider{
		ProviderConfig: conf,
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

	return p, nil
}

func (p *Provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &response{
		XmlnsXsi:          xmlnsXsi,
		XsiSchemaLocation: xsiSchemaLocation,
		ResponseDate:      time.Now().UTC().Format(time.RFC3339),
		Request: request{
			URL: p.BaseURL,
		},
	}

	ctx := r.Context()

	args := r.URL.Query()

	verbs, ok := args["verb"]
	if !ok {
		res.Errors = append(res.Errors, errVerbMissing)
		p.writeResponse(w, res)
		return
	}
	if len(verbs) > 1 {
		res.Errors = append(res.Errors, errVerbRepeated)
		p.writeResponse(w, res)
		return
	}

	res.Request.Verb = verbs[0]

	handlers, ok := verbHandlers[res.Request.Verb]

	if !ok {
		res.Errors = append(res.Errors, errVerbInvalid)
		return
	}

	for i, h := range handlers {
		// only call last handler if there are no errors
		if i == len(handlers)-1 && len(res.Errors) > 0 {
			break
		}

		if err := h(ctx, p, res, args); err != nil {
			p.handleError(w, err)
			return
		}
	}

	p.writeResponse(w, res)
}

func (p *Provider) writeResponse(w http.ResponseWriter, res *response) {
	status := 200
	if len(res.Errors) > 0 {
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
	if p.StyleSheet != "" {
		w.Write([]byte(`<?xml-stylesheet type="text/xsl" href="` + p.StyleSheet + `"?>` + "\n"))
	}
	w.Write(out)
}

func (p *Provider) handleError(w http.ResponseWriter, err error) {
	if p.ErrorHandler != nil {
		p.ErrorHandler(err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// TODO description
func identify(ctx context.Context, p *Provider, res *response, q url.Values) error {
	t, err := p.Backend.GetEarliestDatestamp(ctx)
	if err != nil {
		return err
	}

	res.Body = &Identify{
		RepositoryName:    p.RepositoryName,
		BaseURL:           p.BaseURL,
		ProtocolVersion:   "2.0",
		AdminEmails:       p.AdminEmails,
		Granularity:       p.Granularity,
		Compression:       p.Compression,
		DeletedRecord:     p.DeletedRecord,
		EarliestDatestamp: t.Format(p.dateFormat),
	}

	return nil
}

func listMetadataFormats(ctx context.Context, p *Provider, res *response, q url.Values) error {
	var formats []*MetadataFormat
	var err error
	if identifier := res.Request.Identifier; identifier != "" {
		formats, err = p.Backend.GetRecordMetadataFormats(ctx, identifier)
	} else {
		formats, err = p.Backend.GetMetadataFormats(ctx)
	}

	if err != nil {
		return err
	}

	if len(formats) == 0 {
		res.Errors = append(res.Errors, errNoMetadataFormats)
		return nil
	}

	res.Body = &ListMetadataFormats{
		MetadataFormats: formats,
	}

	return nil
}

func listSets(ctx context.Context, p *Provider, res *response, args url.Values) error {
	var sets []*Set
	var token *ResumptionToken
	var err error
	if rt := res.Request.ResumptionToken; rt != "" {
		sets, token, err = p.Backend.GetMoreSets(ctx, rt)
		if err == ErrBadResumptionToken {
			res.Errors = append(res.Errors, err.(*Error))
			return nil
		}
	} else {
		sets, token, err = p.Backend.GetSets(ctx)
	}

	if err != nil {
		return err
	}

	if len(sets) == 0 {
		res.Errors = append(res.Errors, errNoSetHierarchy)
		return nil
	}

	res.Body = &ListSets{
		Sets:            sets,
		ResumptionToken: token,
	}

	return nil
}

func listIdentifiers(ctx context.Context, p *Provider, res *response, args url.Values) error {
	var headers []*Header
	var token *ResumptionToken
	var err error
	if rt := res.Request.ResumptionToken; rt != "" {
		headers, token, err = p.Backend.GetMoreIdentifiers(ctx, rt)
		if err == ErrBadResumptionToken {
			res.Errors = append(res.Errors, err.(*Error))
			return nil
		}
	} else {
		headers, token, err = p.Backend.GetIdentifiers(ctx,
			res.Request.MetadataPrefix,
			res.Request.Set,
			res.Request.From,
			res.Request.Until,
		)
	}

	if err != nil {
		return err
	}

	if len(headers) == 0 {
		res.Errors = append(res.Errors, errNoRecordsMatch)
		return nil
	}

	res.Body = &ListIdentifiers{
		Headers:         headers,
		ResumptionToken: token,
	}

	return nil
}

func listRecords(ctx context.Context, p *Provider, res *response, args url.Values) error {
	var records []*Record
	var token *ResumptionToken
	var err error
	if rt := res.Request.ResumptionToken; rt != "" {
		records, token, err = p.Backend.GetMoreRecords(ctx, rt)
		if err == ErrBadResumptionToken {
			res.Errors = append(res.Errors, err.(*Error))
			return nil
		}
	} else {
		records, token, err = p.Backend.GetRecords(ctx,
			res.Request.MetadataPrefix,
			res.Request.Set,
			res.Request.From,
			res.Request.Until,
		)
	}

	if err != nil {
		return err
	}

	if len(records) == 0 {
		res.Errors = append(res.Errors, errNoRecordsMatch)
		return nil
	}

	res.Body = &ListRecords{
		Records:         records,
		ResumptionToken: token,
	}

	return nil
}

func getRecord(ctx context.Context, p *Provider, res *response, args url.Values) error {
	rec, err := p.Backend.GetRecord(ctx, res.Request.Identifier, res.Request.MetadataPrefix)
	if err == ErrCannotDisseminateFormat {
		res.Errors = append(res.Errors, err.(*Error))
		return nil
	}
	if err != nil {
		return err
	}
	res.Body = &GetRecord{
		Record: rec,
	}
	return nil
}

func allowArgs(keys ...string) handleFunc {
	keys = append(keys, "verb")
	return func(ctx context.Context, p *Provider, res *response, args url.Values) error {
		for key := range args {
			ok := false
			for _, k := range keys {
				if k == key {
					ok = true
					break
				}
			}
			if !ok {
				res.Errors = append(res.Errors, &Error{Code: "badArgument", Value: fmt.Sprintf("argument %s is illegal", key)})
			}
		}
		return nil
	}
}

func setResumptionToken(ctx context.Context, p *Provider, res *response, args url.Values) error {
	res.Request.ResumptionToken = getArg(res, args, "resumptionToken")
	// only verb and resumptionToken can be set
	if res.Request.ResumptionToken != "" && len(args) > 2 {
		res.Errors = append(res.Errors, errResumptiontokenExclusive)
	}
	return nil
}

func setIdentifier(ctx context.Context, p *Provider, res *response, args url.Values) error {
	val := getArg(res, args, "identifier")

	if val != "" {
		exists, err := p.Backend.HasRecord(ctx, val)
		if err != nil {
			return err
		}
		if !exists {
			res.Errors = append(res.Errors, errIDDoesNotExist)
			return nil
		}

		res.Request.Identifier = val
	}

	return nil
}

func setRequiredMetadataPrefix(ctx context.Context, p *Provider, res *response, args url.Values) error {
	if res.Request.ResumptionToken != "" {
		return nil
	}

	val := getArg(res, args, "metadataPrefix")

	if val == "" {
		res.Errors = append(res.Errors, errMetadataPrefixMissing)
		return nil
	}

	exists, err := p.Backend.HasMetadataFormat(ctx, val)
	if err != nil {
		return err
	}
	if !exists {
		res.Errors = append(res.Errors, ErrCannotDisseminateFormat)
		return nil
	}

	res.Request.MetadataPrefix = val

	return nil
}

func setRequiredIdentifier(ctx context.Context, p *Provider, res *response, args url.Values) error {
	val := getArg(res, args, "identifier")

	if val == "" {
		res.Errors = append(res.Errors, errIdentifierMissing)
	}

	exists, err := p.Backend.HasRecord(ctx, val)
	if err != nil {
		return err
	}
	if !exists {
		res.Errors = append(res.Errors, errIDDoesNotExist)
		return nil
	}

	res.Request.Identifier = val

	return nil
}

func setSet(ctx context.Context, p *Provider, res *response, args url.Values) error {
	if res.Request.ResumptionToken != "" {
		return nil
	}

	val := getArg(res, args, "set")

	if val != "" {
		hasSets, err := p.Backend.HasSets(ctx)
		if err != nil {
			return err
		}
		if !hasSets {
			res.Errors = append(res.Errors, errNoSetHierarchy)
			return nil
		}

		setExists, err := p.Backend.HasSet(ctx, val)
		if err != nil {
			return err
		}
		if !setExists {
			res.Errors = append(res.Errors, errSetDoesNotExist)
			return nil
		}

		res.Request.Set = val
	}

	return nil
}

// TODO parse dates and check if from <= until
func setFromUntil(ctx context.Context, p *Provider, res *response, args url.Values) error {
	if res.Request.ResumptionToken != "" {
		return nil
	}

	f := getArg(res, args, "from")
	u := getArg(res, args, "until")

	if f != "" {
		if p.Granularity == "YYYY-MM-DDThh:mm:ssZ" && len(f) == 10 {
			f += "T00:00:00Z"
		}
		if _, err := time.Parse(p.dateFormat, f); err == nil {
			res.Request.From = f
		} else {
			res.Errors = append(res.Errors, errFromInvalid)
		}
	}
	if u != "" {
		if p.Granularity == "YYYY-MM-DDThh:mm:ssZ" && len(u) == 10 {
			u += "T00:00:00Z"
		}
		if _, err := time.Parse(p.dateFormat, u); err == nil {
			res.Request.Until = u
		} else {
			res.Errors = append(res.Errors, errUntilInvalid)
		}
	}

	return nil
}

func getArg(res *response, q url.Values, attr string) string {
	vals, ok := q[attr]
	if !ok {
		return ""
	}

	if len(vals) > 1 {
		err := &Error{Code: "badArgument", Value: fmt.Sprintf("%s can't be repeated", attr)}
		res.Errors = append(res.Errors, err)
		return ""
	}

	if vals[0] == "" {
		err := &Error{Code: "badArgument", Value: fmt.Sprintf("%s is missing", attr)}
		res.Errors = append(res.Errors, err)
		return ""
	}

	return vals[0]
}
