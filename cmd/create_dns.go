// SPDX-License-Identifier: MPL-2.0
//
// Copyright (c) 2023 Mark Cornick
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type createDNSInput struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Priority int    `json:"priority"`
	TTL      int    `json:"ttl"`
}

type createDNSOutput struct {
	Request  resultRequest `json:"request"`
	Response struct {
		Message  string `json:"message"`
		DataSent struct {
			Type     string      `json:"type"`
			Priority *int        `json:"priority"`
			TTL      json.Number `json:"ttl"`
			Name     string      `json:"name"`
			Content  string      `json:"content"`
		} `json:"data_sent"`
		ResponseReceived struct {
			Data struct {
				ID        json.Number `json:"id"`
				Name      string      `json:"name"`
				Content   string      `json:"content"`
				TTL       json.Number `json:"ttl"`
				Priority  *int        `json:"priority"`
				Type      string      `json:"type"`
				CreatedAt time.Time   `json:"created_at"`
				UpdatedAt time.Time   `json:"updated_at"`
			} `json:"data"`
		} `json:"response_received"`
	} `json:"response"`
}

var (
	createDNSPriority int
	createDNSTTL      int
	createDNSCmd      = &cobra.Command{
		Use:   "dns [name] [type] [data]",
		Short: "Create a DNS record",
		Long:  "Creates a DNS record.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := createDNS(args[0], args[1], args[2], createDNSPriority, createDNSTTL)
			if err != nil {
				return err
			}
			if result.Request.Success {
				fmt.Println(result.Response.Message)
			} else {
				return fmt.Errorf("%s", result.Response.Message)
			}
			return nil
		},
	}
)

func init() {
	createDNSCmd.Flags().IntVarP(
		&createDNSPriority,
		"priority",
		"p",
		0,
		"priority of the DNS record",
	)
	createDNSCmd.Flags().IntVarP(
		&createDNSTTL,
		"ttl",
		"T",
		3600,
		"time to live of the DNS record",
	)
	createCmd.AddCommand(createDNSCmd)
}

func createDNS(name string, recordType string, content string, priority int, ttl int) (createDNSOutput, error) {
	var result createDNSOutput
	dns := createDNSInput{strings.ToUpper(recordType), name, content, priority, ttl}
	body, err := callAPIWithParams(
		http.MethodPost,
		"/address/"+viper.GetString("address")+"/dns",
		dns,
		true,
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}
