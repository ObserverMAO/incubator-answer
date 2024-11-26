/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package entity

import "time"

// UserExternalLogin user external login
type UserExternalLogin struct {
	ID         int64     `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt  time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt  time.Time `xorm:"updated TIMESTAMP updated_at"`
	UserID     string    `xorm:"not null default 0 BIGINT(20) user_id"`
	Provider   string    `xorm:"not null default '' VARCHAR(100) provider"`
	ExternalID string    `xorm:"not null default '' VARCHAR(128) external_id"`
	MetaInfo   string    `xorm:"TEXT meta_info"`
}

// TableName  table name
func (UserExternalLogin) TableName() string {
	return "user_external_login"
}

type MixinUserInfoResponse struct {
	Data MixinUserInfo `json:"data"`
}

type MixinUserInfo struct {
	Type                          string      `json:"type"`
	UserID                        string      `json:"user_id"`
	IdentityNumber                string      `json:"identity_number"`
	Phone                         string      `json:"phone"`
	FullName                      string      `json:"full_name"`
	Biography                     string      `json:"biography"`
	AvatarURL                     string      `json:"avatar_url"`
	Relationship                  string      `json:"relationship"`
	MuteUntil                     time.Time   `json:"mute_until"`
	CreatedAt                     time.Time   `json:"created_at"`
	IsVerified                    bool        `json:"is_verified"`
	IsScam                        bool        `json:"is_scam"`
	IsDeactivated                 bool        `json:"is_deactivated"`
	CodeID                        string      `json:"code_id"`
	CodeURL                       string      `json:"code_url"`
	Features                      interface{} `json:"features"`
	HasSafe                       bool        `json:"has_safe"`
	Membership                    Membership  `json:"membership"`
	Email                         string      `json:"email"`
	SessionID                     string      `json:"session_id"`
	DeviceStatus                  string      `json:"device_status"`
	HasPin                        bool        `json:"has_pin"`
	SaltExportedAt                time.Time   `json:"salt_exported_at"`
	ReceiveMessageSource          string      `json:"receive_message_source"`
	AcceptConversationSource      string      `json:"accept_conversation_source"`
	AcceptSearchSource            string      `json:"accept_search_source"`
	FiatCurrency                  string      `json:"fiat_currency"`
	TransferNotificationThreshold int64       `json:"transfer_notification_threshold"`
	TransferConfirmationThreshold int64       `json:"transfer_confirmation_threshold"`
	PinTokenBase64                string      `json:"pin_token_base64"`
	PinToken                      string      `json:"pin_token"`
	SaltBase64                    string      `json:"salt_base64"`
	TipKeyBase64                  string      `json:"tip_key_base64"`
	TipCounter                    int64       `json:"tip_counter"`
	SpendPublicKey                string      `json:"spend_public_key"`
	HasEmergencyContact           bool        `json:"has_emergency_contact"`
}

type Membership struct {
	Plan      string    `json:"plan"`
	ExpiredAt time.Time `json:"expired_at"`
}
