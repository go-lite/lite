package swagger

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type API struct {
	Routes     []Route
	Components Components

	openapi3.T
}

type Components struct {
	Schemas         map[string]*openapi3.SchemaRef
	Responses       map[string]*openapi3.ResponseRef
	Parameters      map[string]*openapi3.ParameterRef
	SecuritySchemes map[string]*openapi3.SecuritySchemeRef
}

//func NewAPI() *API {
//	return &API{
//		Routes: []Route{},
//		T: openapi3.T{
//			OpenAPI: "3.0.0",
//			Info: &openapi3.Info{
//				Title:   "API Documentation",
//				Version: "1.0.0",
//			},
//			Paths:      &openapi3.Paths{},
//			Components: &openapi3.Components{},
//		},
//		Components: Components{
//			Schemas:         make(map[string]*openapi3.SchemaRef),
//			Responses:       make(map[string]*openapi3.ResponseRef),
//			Parameters:      make(map[string]*openapi3.ParameterRef),
//			SecuritySchemes: make(map[string]*openapi3.SecuritySchemeRef),
//		},
//	}
//}
//
//func (api *API) AddServer(desc, url string) {
//	api.Servers = append(api.Servers, &openapi3.Server{
//		URL:         url,
//		Description: desc,
//	})
//}

//func (api *API) RegisterRoute(method, path, name string, tags []string, reqRes RequestResponse) {
//	route := Route{
//		Method:           method,
//		Path:             path,
//		Name:             name,
//		Tags:             tags,
//		RequestsResponse: reqRes,
//	}
//	api.Routes = append(api.Routes, route)
//}
//
//func (api *API) GenerateSwagger(filename string) error {
//
//	for name, schema := range api.Components.Schemas {
//		api.Components.Schemas[name] = schema
//	}
//
//	for name, response := range api.Components.Responses {
//		api.Components.Responses[name] = response
//	}
//
//	for name, parameter := range api.Components.Parameters {
//		api.Components.Parameters[name] = parameter
//	}
//
//	for name, securityScheme := range api.Components.SecuritySchemes {
//		api.Components.SecuritySchemes[name] = securityScheme
//	}
//
//	for _, route := range api.Routes {
//		operation := &openapi3.Operation{
//			Summary:     route.Name,
//			Description: route.Name,
//			Tags:        route.Tags,
//		}
//
//		parameters := openapi3.Parameters{}
//		for _, header := range route.RequestsResponse.Request.Headers {
//			parameters = append(parameters, &openapi3.ParameterRef{
//				Value: &openapi3.Parameter{
//					In:       "header",
//					Name:     header.Key,
//					Required: true,
//					Schema: &openapi3.SchemaRef{
//						Value: &openapi3.Schema{
//							Type: string(header.Value.Type),
//						},
//					},
//				},
//			})
//		}
//
//		for _, param := range route.RequestsResponse.Request.PathParameters {
//			parameters = append(parameters, &openapi3.ParameterRef{
//				Value: &openapi3.Parameter{
//					In:       "path",
//					Name:     param.Key,
//					Required: true,
//					Schema: &openapi3.SchemaRef{
//						Value: &openapi3.Schema{
//							Type: string(param.Value.Type),
//						},
//					},
//				},
//			})
//		}
//
//		for _, query := range route.RequestsResponse.Request.QueryParameters {
//			parameters = append(parameters, &openapi3.ParameterRef{
//				Value: &openapi3.Parameter{
//					In:       "query",
//					Name:     query.Key,
//					Required: true,
//					Schema: &openapi3.SchemaRef{
//						Value: &openapi3.Schema{
//							Type: string(query.Value.Type),
//						},
//					},
//				},
//			})
//		}
//
//		operation.Parameters = parameters
//
//		if len(route.RequestsResponse.Request.Body) > 0 {
//			requestBody := &openapi3.RequestBody{}
//			for contentType, value := range route.RequestsResponse.Request.Body {
//				sc, enc, err := valueToSchema(value, contentType, 1)
//				if err != nil {
//					return err
//				}
//
//				requestBody.Content = openapi3.Content{
//					contentType: &openapi3.MediaType{
//						Schema:   sc,
//						Encoding: enc,
//					},
//				}
//			}
//			operation.RequestBody = &openapi3.RequestBodyRef{
//				Value: requestBody,
//			}
//		}
//
//		if len(route.RequestsResponse.Response.Body) > 0 {
//			responses := openapi3.Responses{}
//			for status, content := range route.RequestsResponse.Response.Body {
//				for contentType, value := range content {
//					description := fmt.Sprintf("Response %d", status)
//
//					sc, enc, err := valueToSchema(value, contentType, 1)
//					if err != nil {
//						return err
//					}
//
//					responses[fmt.Sprintf("%d", status)] = &openapi3.ResponseRef{
//						Value: &openapi3.Response{
//							Description: &description,
//							Content: openapi3.Content{
//								contentType: &openapi3.MediaType{
//									Schema:   sc,
//									Encoding: enc,
//								},
//							},
//						},
//					}
//				}
//			}
//			operation.Responses = responses
//		}
//
//		if api.Paths[route.Path] == nil {
//			api.Paths[route.Path] = &openapi3.PathItem{}
//		}
//
//		switch route.Method {
//		case "get":
//			api.Paths[route.Path].Get = operation
//		case "post":
//			api.Paths[route.Path].Post = operation
//		case "put":
//			api.Paths[route.Path].Put = operation
//		case "delete":
//			api.Paths[route.Path].Delete = operation
//		case "patch":
//			api.Paths[route.Path].Patch = operation
//		case "options":
//			api.Paths[route.Path].Options = operation
//		case "head":
//			api.Paths[route.Path].Head = operation
//		case "trace":
//			api.Paths[route.Path].Trace = operation
//		default:
//			return errors.New("unsupported method")
//		}
//	}
//
//	data, err := json.MarshalIndent(api.T, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	var yamlData map[string]interface{}
//	if err := yaml.Unmarshal(data, &yamlData); err != nil {
//		return err
//	}
//
//	yamlBytes, err := yaml.Marshal(&yamlData)
//	if err != nil {
//		return err
//	}
//
//	return os.WriteFile(filename, yamlBytes, 0644)
//}
//
//func valueToSchema(
//	value Value,
//	ct string,
//	depth int,
//) (sc *openapi3.SchemaRef, enc map[string]*openapi3.Encoding, err error) {
//	s := new(openapi3.Schema)
//	s.Nullable = value.Nullable
//
//	switch value.Type {
//	case ValueTypeString:
//		s.Type = openapi3.TypeString
//
//		switch value.Format {
//		case ValueFormatDateTime:
//			s.Format = "date-time"
//		case ValueFormatBinary:
//			s.Format = "binary"
//		default:
//		}
//
//	case ValueTypeUint:
//		s.Min = new(float64)
//		fallthrough
//
//	case ValueTypeInt:
//		s.Type = openapi3.TypeInteger
//
//	case ValueTypeFloat:
//		s.Type = openapi3.TypeNumber
//		s.Format = "double"
//
//	case ValueTypeBool:
//		s.Type = openapi3.TypeBoolean
//
//	case ValueTypeObject:
//		s.Type = openapi3.TypeObject
//		for key, property := range value.Properties {
//			if s.Properties == nil {
//				s.Properties = make(map[string]*openapi3.SchemaRef)
//			}
//
//			var tmpEnc map[string]*openapi3.Encoding
//
//			sc, tmpEnc, err = valueToSchema(property, ct, depth+1)
//			if err != nil {
//				return nil, nil, err
//			}
//
//			s.Properties[key] = sc
//
//			if vEnc, ok := tmpEnc["def"]; ok {
//				if enc == nil {
//					enc = make(map[string]*openapi3.Encoding)
//				}
//
//				enc[key] = vEnc
//			}
//
//			s.Properties[key] = sc
//
//			if !property.Nullable {
//				s.Required = append(s.Required, key)
//			}
//		}
//
//		slices.SortFunc(s.Required, strings.Compare)
//		slices.Sort(s.Required)
//
//	case ValueTypeArray:
//		s.Type = openapi3.TypeArray
//
//		sc, enc, err = valueToSchema(*value.Items, ct, depth) // NOTE: don't increment depth here
//		if err != nil {
//			return nil, nil, err
//		}
//
//		if _, ok := enc["def"]; ok {
//			return nil, nil, errors.New("encoding not supported for array")
//		}
//
//		s.Items = sc
//
//	case ValueTypeMap:
//		s.Type = openapi3.TypeObject
//
//		if value.Keys.Type != ValueTypeString {
//			return nil, nil, errors.New("unsupported map key type")
//		}
//
//		sc, _, err = valueToSchema(*value.Values, ct, depth+1)
//		if err != nil {
//			return nil, nil, err
//		}
//
//		s.AdditionalProperties.Schema = sc
//
//	case ValueTypeAny:
//		s.Description = "any type"
//	}
//
//	if len(value.Enum) > 0 {
//		s.Enum = make([]interface{}, len(value.Enum))
//		for i, v := range value.Enum {
//			s.Enum[i] = v
//		}
//
//		if s.Nullable {
//			s.Enum = append(s.Enum, "null")
//			s.Nullable = false
//		}
//
//		sort.Slice(s.Enum, func(i, j int) bool {
//			vI, okI := s.Enum[i].(string)
//			vJ, okJ := s.Enum[j].(string)
//
//			return okI && okJ && (strings.Compare(vI, vJ) < 0)
//		})
//	}
//
//	if value.XMLName != "" {
//		s.XML = &openapi3.XML{
//			Name: value.XMLName,
//		}
//	}
//
//	if depth == 1 && (ct == "multipart/form-data" || ct == "multipart/mixed") && len(value.ContentType) > 0 {
//		enc = make(map[string]*openapi3.Encoding)
//		enc["def"] = &openapi3.Encoding{
//			ContentType: strings.Join(value.ContentType, ","),
//		}
//	}
//
//	sc = &openapi3.SchemaRef{
//		Value: s,
//	}
//
//	return sc, enc, nil
//}
