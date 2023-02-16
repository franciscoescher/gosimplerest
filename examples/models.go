package examples

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
  `hours` int DEFAULT NULL,
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
	"github.com/franciscoescher/gosimplerest/resource"

	"gopkg.in/guregu/null.v3"
)

var UserResource = resource.Resource{
	Table:      "users",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":        {Validator: "uuid4"},
		"first_name":  {},
		"last_name":   {},
		"phone":       {},
		"credit_card": {Unsearchable: true},
		"created_at":  {},
		"deleted_at":  {},
		"updated_at":  {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}

var RentEventResource = resource.Resource{
	Table:      "rent_events",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":          {Validator: "uuid4"},
		"user_id":       {Validator: "uuid4"},
		"vehicle_id":    {Validator: "uuid4"},
		"starting_time": {},
		"hours":         {},
		"checkin_time":  {},
		"dropoff_time":  {},
		"cancel_time":   {},
		"created_at":    {},
		"deleted_at":    {},
		"updated_at":    {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}

var VehicleResource = resource.Resource{
	Table:      "vehicles",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":           {Validator: "uuid4"},
		"license_plate":  {},
		"state":          {},
		"archived":       {},
		"year":           {},
		"price_per_hour": {},
		"lot":            {},
		"created_at":     {},
		"deleted_at":     {},
		"updated_at":     {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}
