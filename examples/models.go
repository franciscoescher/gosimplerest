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

type Users struct {
	UUID       string      `json:"uuid" primary_key:"true"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Phone      string      `json:"phone"`
	CreditCard string      `json:"credit_card"`
	CreatedAt  string      `json:"created_at" created_at:"true"`
	DeletedAt  null.String `json:"deleted_at" soft_delete:"true"`
	UpdatedAt  null.String `json:"updated_at" updated_at:"true"`
}

var UserResource = resource.Resource{
	Data: Users{},
}

type RentEvents struct {
	UUID         string      `json:"uuid" primary_key:"true"`
	UserID       string      `json:"user_id" belongs_to:"users"`
	VehicleID    string      `json:"vehicle_id" belongs_to:"vehicles"`
	StartingTIme string      `json:"starting_time"`
	Hours        int         `json:"hours"`
	CheckinTime  string      `json:"checkin_time"`
	DropoffTime  string      `json:"dropoff_time"`
	CancelTime   string      `json:"cancel_time"`
	CreatedAt    string      `json:"created_at" created_at:"true"`
	DeletedAt    null.String `json:"deleted_at" soft_delete:"true"`
	UpdatedAt    null.String `json:"updated_at" updated_at:"true"`
}

var RentEventResource = resource.Resource{
	Data: RentEvents{},
}

type Vehicles struct {
	UUID         string      `json:"uuid" primary_key:"true"`
	LicensePlate string      `json:"license_plate"`
	State        string      `json:"state"`
	Archived     bool        `json:"archived"`
	Year         int         `json:"year"`
	PricePerHour float64     `json:"price_per_hour"`
	Lot          int         `json:"lot"`
	CreatedAt    string      `json:"created_at" created_at:"true"`
	DeletedAt    null.String `json:"deleted_at" soft_delete:"true"`
	UpdatedAt    null.String `json:"updated_at" updated_at:"true"`
}

var VehicleResource = resource.Resource{
	Data: Vehicles{},
}
