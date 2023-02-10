package main

// Example implementation for the following mysql data structures:
/*
DROP TABLE IF EXISTS `rent_events`;
DROP TABLE IF EXISTS `vehicles`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `uuid` varchar(191) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `phone` longtext,
  `credit_card` longtext,
  PRIMARY KEY (`uuid`)
);

CREATE TABLE `vehicles` (
  `uuid` varchar(191) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `license_plate` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT '0',
  `year` smallint DEFAULT NULL,
  `price_per_hour` double DEFAULT NULL,
  `lot` bigint DEFAULT NULL,
  PRIMARY KEY (`uuid`)
);

CREATE TABLE `rent_events` (
  `uuid` varchar(191) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` varchar(128) DEFAULT NULL,
  `vehicle_id` varchar(128) DEFAULT NULL,
  `starting_time` datetime(3) DEFAULT NULL,
  `hours` bigint DEFAULT NULL,
  `checkin_time` datetime(3) DEFAULT NULL,
  `dropoff_time` datetime(3) DEFAULT NULL,
  `cancel_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`uuid`),
  KEY `fk_rent_events_vehicle` (`vehicle_id`),
  KEY `fk_rent_events_user` (`user_id`),
  CONSTRAINT `fk_rent_events_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`uuid`) ON DELETE CASCADE,
  CONSTRAINT `fk_rent_events_vehicle` FOREIGN KEY (`vehicle_id`) REFERENCES `vehicles` (`uuid`) ON DELETE CASCADE
);
*/

import (
	"fmt"
	"strconv"
	"time"

	"github.com/franciscoescher/gosimplerest"
	"github.com/gofrs/uuid"

	"gopkg.in/guregu/null.v3"
)

var UserResource = gosimplerest.Resource{
	Table:      "users",
	PrimaryKey: "uuid",
	Fields: map[string]gosimplerest.Field{
		"uuid":        {Validator: validateUUID},
		"first_name":  {Validator: validateLenght(1)},
		"last_name":   {Validator: validateLenght(1)},
		"phone":       {},
		"credit_card": {},
		"created_at":  {Validator: validateTime},
		"deleted_at":  {Validator: validateTime},
		"updated_at":  {Validator: validateTime},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}

var RentEventResource = gosimplerest.Resource{
	Table:      "rent_events",
	PrimaryKey: "uuid",
	Fields: map[string]gosimplerest.Field{
		"uuid":          {Validator: validateUUID},
		"user_id":       {Validator: validateUUID},
		"vehicle_id":    {Validator: validateUUID},
		"starting_time": {Validator: validateTime},
		"hours":         {Validator: validateIntPositive},
		"checkin_time":  {Validator: validateTime},
		"dropoff_time":  {Validator: validateTime},
		"cancel_time":   {Validator: validateTime},
		"created_at":    {Validator: validateTime},
		"deleted_at":    {Validator: validateTime},
		"updated_at":    {Validator: validateTime},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
	BelongsToFields: []gosimplerest.BelongsTo{
		{Table: "users", Field: "user_id"},
		{Table: "vehicles", Field: "vehicle_id"},
	},
}

var VehicleResource = gosimplerest.Resource{
	Table:      "vehicles",
	PrimaryKey: "uuid",
	Fields: map[string]gosimplerest.Field{
		"uuid":           {Validator: validateUUID},
		"license_plate":  {},
		"state":          {},
		"archived":       {},
		"year":           {Validator: validateIntPositive},
		"price_per_hour": {},
		"lot":            {Validator: validateIntPositive},
		"created_at":     {Validator: validateTime},
		"deleted_at":     {Validator: validateTime},
		"updated_at":     {Validator: validateTime},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}

func validateUUID(field string, val interface{}) error {
	if val == nil {
		return fmt.Errorf("%s is required", field)
	}
	_, err := uuid.FromString(val.(string))
	if err != nil {
		return fmt.Errorf(field+" is invalid: %s", val)
	}
	return nil
}

func validateIntPositive(field string, val interface{}) error {
	s, ok := val.(string)
	if !ok {
		return fmt.Errorf("%s is invalid: %s", field, val)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("%s is invalid: %s", field, val)
	}
	if i <= 0 {
		return fmt.Errorf("%s must be positive: %s", field, val)
	}
	return nil
}

func validateTime(field string, val interface{}) error {
	if val == nil {
		return fmt.Errorf("%s can't be null", field)
	}
	_, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		return fmt.Errorf("%s is invalid: %s", field, val)
	}
	return nil
}

func validateLenght(i int) gosimplerest.ValidatorFunc {
	return func(field string, val interface{}) error {
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s is invalid: %s", field, val)
		}
		if len(s) < i {
			return fmt.Errorf("%s must have %d characters: %s", field, i, val)
		}
		return nil
	}
}
