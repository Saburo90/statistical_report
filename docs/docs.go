// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2018-11-29 18:02:12.4137186 +0800 CST m=+0.071288701

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "title": "统计报表API",
        "contact": {},
        "license": {},
        "version": "1.0"
    },
    "host": "39.106.39.7:8092",
    "basePath": "/statistical",
    "paths": {
        "/user/getOverview": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户统计"
                ],
                "summary": "获取用户统计数据",
                "parameters": [
                    {
                        "description": "调用者,调用时间,调用者IP,签名",
                        "name": "operationSign",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/protocol.OverviewReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":0,\"msg\":\"OK\",\"data\":{}}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/protocol.OverviewResp"
                        }
                    },
                    "400": {
                        "description": "{\"code\":exceptionCode,\"msg\":exceptionMsg,\"data\":{}}"
                    }
                }
            }
        }
    },
    "definitions": {
        "protocol.OverviewReq": {
            "type": "object",
            "properties": {
                "clientIP": {
                    "type": "string"
                },
                "operateTime": {
                    "type": "integer"
                },
                "operator": {
                    "type": "string"
                },
                "sign": {
                    "type": "string"
                }
            }
        },
        "protocol.OverviewResp": {
            "type": "object",
            "properties": {
                "total_roam_users": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
