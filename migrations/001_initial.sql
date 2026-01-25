-- Migration: Initial Schema
-- Generated automatically from GORM models

-- ====================================
-- CREATE TABLES
-- ====================================

CREATE TABLE IF NOT EXISTS stich."Channels" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT,
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  owner_user_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Customers" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  first_name TEXT,
  last_name TEXT,
  email TEXT,
  phone_number TEXT,
  whatsapp_number TEXT,
  address TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."DressTypes" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT,
  measurements TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."EmailNotifications" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  status TEXT,
  to_mail_address TEXT,
  subject TEXT,
  body TEXT,
  notification_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."EnquiryHistories" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  status TEXT,
  employee_comment TEXT,
  customer_comment TEXT,
  visiting_date TIMESTAMPTZ,
  call_back_date TIMESTAMPTZ,
  enquiry_date TIMESTAMPTZ,
  response_status TEXT,
  enquiry_id INTEGER,
  employee_id INTEGER,
  performed_at TIMESTAMPTZ NOT NULL,
  performed_by_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Enquiries" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  subject TEXT,
  status TEXT NOT NULL,
  notes TEXT,
  source TEXT,
  referred_by TEXT,
  referrer_phone_number TEXT,
  customer_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."MasterConfigs" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT,
  type TEXT,
  current_value TEXT,
  previous_value TEXT,
  default_value TEXT,
  use_default BOOLEAN,
  description TEXT,
  format TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Measurements" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  value JSONB,
  person_id INTEGER,
  dress_type_id INTEGER,
  taken_by_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."MeasurementHistories" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  action TEXT NOT NULL,
  old_values JSONB,
  measurement_id INTEGER,
  performed_at TIMESTAMPTZ NOT NULL,
  performed_by_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Notifications" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  status TEXT NOT NULL DEFAULT 'PENDING',
  source_entity TEXT,
  entity_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."OrderHistories" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  action TEXT NOT NULL,
  changed_fields TEXT,
  status TEXT,
  expected_delivery_date TIMESTAMPTZ,
  delivered_date TIMESTAMPTZ,
  order_item_id INTEGER,
  order_item_data JSONB,
  order_id INTEGER,
  performed_at TIMESTAMPTZ NOT NULL,
  performed_by_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Orders" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  status TEXT,
  notes TEXT,
  expected_delivery_date TIMESTAMPTZ,
  delivered_date TIMESTAMPTZ,
  customer_id INTEGER,
  order_taken_by_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."OrderItems" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  description TEXT,
  quantity INTEGER,
  price TEXT,
  total TEXT,
  expected_delivery_date TIMESTAMPTZ,
  delivered_date TIMESTAMPTZ,
  person_id INTEGER,
  measurement_id INTEGER,
  order_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Persons" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT,
  gender TEXT,
  age INTEGER,
  customer_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."UserChannelDetails" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  user_id INTEGER,
  user_channel_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."UserConfigs" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  config TEXT,
  user_id INTEGER,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."Users" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  first_name TEXT,
  last_name TEXT,
  extension TEXT NOT NULL,
  phone_number TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  role TEXT NOT NULL,
  is_login_disabled BOOLEAN,
  is_logged_in BOOLEAN,
  last_login_time TIMESTAMPTZ,
  login_failure_counter INTEGER,
  reset_password_string TEXT,
  experience TEXT,
  department TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stich."WhatsappNotifications" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  status TEXT,
  receipient_number TEXT,
  receipient_extension TEXT,
  subject TEXT,
  body TEXT,
  notification_id INTEGER,
  PRIMARY KEY (id)
);


-- ====================================
-- ADD FOREIGN KEY CONSTRAINTS
-- ====================================

ALTER TABLE stich."Channels" ADD CONSTRAINT fk_Channel_owner_user_id FOREIGN KEY (owner_user_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."EmailNotifications" ADD CONSTRAINT fk_EmailNotification_notification_id FOREIGN KEY (notification_id) REFERENCES stich."Notifications" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."EnquiryHistories" ADD CONSTRAINT fk_EnquiryHistory_employee_id FOREIGN KEY (employee_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."EnquiryHistories" ADD CONSTRAINT fk_EnquiryHistory_performed_by_id FOREIGN KEY (performed_by_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Enquiries" ADD CONSTRAINT fk_Enquiry_customer_id FOREIGN KEY (customer_id) REFERENCES stich."Customers" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Measurements" ADD CONSTRAINT fk_Measurement_person_id FOREIGN KEY (person_id) REFERENCES stich."Persons" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Measurements" ADD CONSTRAINT fk_Measurement_dress_type_id FOREIGN KEY (dress_type_id) REFERENCES stich."DressTypes" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Measurements" ADD CONSTRAINT fk_Measurement_taken_by_id FOREIGN KEY (taken_by_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."MeasurementHistories" ADD CONSTRAINT fk_MeasurementHistory_measurement_id FOREIGN KEY (measurement_id) REFERENCES stich."Measurements" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."MeasurementHistories" ADD CONSTRAINT fk_MeasurementHistory_performed_by_id FOREIGN KEY (performed_by_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."OrderHistories" ADD CONSTRAINT fk_OrderHistory_order_id FOREIGN KEY (order_id) REFERENCES stich."Orders" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."OrderHistories" ADD CONSTRAINT fk_OrderHistory_performed_by_id FOREIGN KEY (performed_by_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Orders" ADD CONSTRAINT fk_Order_customer_id FOREIGN KEY (customer_id) REFERENCES stich."Customers" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Orders" ADD CONSTRAINT fk_Order_order_taken_by_id FOREIGN KEY (order_taken_by_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."OrderItems" ADD CONSTRAINT fk_OrderItem_person_id FOREIGN KEY (person_id) REFERENCES stich."Persons" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."OrderItems" ADD CONSTRAINT fk_OrderItem_measurement_id FOREIGN KEY (measurement_id) REFERENCES stich."Measurements" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."OrderItems" ADD CONSTRAINT fk_OrderItem_order_id FOREIGN KEY (order_id) REFERENCES stich."Orders" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."Persons" ADD CONSTRAINT fk_Person_customer_id FOREIGN KEY (customer_id) REFERENCES stich."Customers" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."UserChannelDetails" ADD CONSTRAINT fk_UserChannelDetail_user_id FOREIGN KEY (user_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."UserChannelDetails" ADD CONSTRAINT fk_UserChannelDetail_user_channel_id FOREIGN KEY (user_channel_id) REFERENCES stich."Channels" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."UserConfigs" ADD CONSTRAINT fk_UserConfig_user_id FOREIGN KEY (user_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."WhatsappNotifications" ADD CONSTRAINT fk_WhatsappNotification_notification_id FOREIGN KEY (notification_id) REFERENCES stich."Notifications" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;


-- Migration completed successfully
-- Tables created: 18
-- Indexes created: 0
-- Foreign keys added: 22
