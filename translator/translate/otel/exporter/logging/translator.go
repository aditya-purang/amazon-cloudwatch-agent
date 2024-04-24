// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package logging

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"

	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

type translator struct {
	name    string
	factory exporter.Factory
}

var _ common.Translator[component.Config] = (*translator)(nil)

func NewTranslator() common.Translator[component.Config] {
	t := &translator{factory: loggingexporter.NewFactory()}
	return t
}

func (t *translator) ID() component.ID {
	return component.NewIDWithName(t.factory.Type(), "")
}

func (t *translator) Translate(conf *confmap.Conf) (component.Config, error) {
	if conf == nil || !conf.IsSet(common.DebugLogging) {
		return nil, &common.MissingKeyError{ID: t.ID(), JsonKey: common.DebugLogging}
	}
	if enabled, _ := common.GetBool(conf, common.DebugLogging); !enabled {
		return nil, &common.DebugLoggingEnabledError{ID: t.ID(), JsonKey: common.DebugLogging}
	}

	cfg := t.factory.CreateDefaultConfig().(*loggingexporter.Config)

	c := confmap.NewFromStringMap(map[string]interface{}{
		"verbosity": configtelemetry.LevelDetailed,
	})

	if err := c.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal into loggingexporter config: %w", err)
	}

	return cfg, nil
}
