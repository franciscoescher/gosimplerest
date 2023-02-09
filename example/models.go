package main

// Example implementation for the following mysql data structures:
/*
DROP TABLE IF EXISTS `rent_events`;
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

DROP TABLE IF EXISTS `vehicles`;
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
*/

import (
	"franciscoescher/gosimplerest"

	"gopkg.in/guregu/null.v3"
)

var UserResource = gosimplerest.Resource{
	Table:      "users",
	PrimaryKey: "uuid",
	Fields: []gosimplerest.Field{
		{Name: "uuid"},
		{Name: "first_name"},
		{Name: "last_name"},
		{Name: "phone"},
		{Name: "credit_card"},
		{Name: "created_at"},
		{Name: "deleted_at"},
		{Name: "updated_at"},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}

var RentEventResource = gosimplerest.Resource{
	Table:      "rent_events",
	PrimaryKey: "uuid",
	Fields: []gosimplerest.Field{
		{Name: "uuid"},
		{Name: "user_id"},
		{Name: "vehicle_id"},
		{Name: "starting_time"},
		{Name: "hours"},
		{Name: "checkin_time"},
		{Name: "dropoff_time"},
		{Name: "cancel_time"},
		{Name: "created_at"},
		{Name: "deleted_at"},
		{Name: "updated_at"},
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
	Fields: []gosimplerest.Field{
		{Name: "uuid"},
		{Name: "license_plate"},
		{Name: "state"},
		{Name: "archived"},
		{Name: "year"},
		{Name: "price_per_hour"},
		{Name: "lot"},
		{Name: "created_at"},
		{Name: "deleted_at"},
		{Name: "updated_at"},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	CreatedAtField:  null.NewString("created_at", true),
	UpdatedAtField:  null.NewString("updated_at", true),
}
