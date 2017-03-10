
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table `accounts` (
    `id` varchar(255) primary key,
    
    `name` text not null,
    `description` text not null,
    `status` text not null
);

create table `payment_historys` (
    `id` varchar(255) primary key,
    
    `name` text not null,
    `description` text not null,
    `invoice_id` text null,
    `payment_method` text null,
    `type` text null,
    `account_id` varchar(255) not null,
    `date` text null,
    `payment_info` text null,
    
    constraint `payment_historys_account_id_accounts_id` foreign key(`account_id`) REFERENCES `accounts`(id)
);

create table `payment_methods` (
    `id` varchar(255) primary key,
    
    `name` text not null,
    `description` text not null,
    `token` text null,
    `account_id` varchar(255) not null,
    
    constraint `payment_methods_account_id_accounts_id` foreign key(`account_id`) REFERENCES `accounts`(id)
);

create table `bills` (
    `id` varchar(255) primary key,
    
    `name` text not null,
    `description` text not null,
    `month` text null,
    `amount` integer not null,
    `detail` text null,
    `account_id` varchar(255) not null,
    
    constraint `bills_account_id_accounts_id` foreign key(`account_id`) REFERENCES `accounts`(id)
 );

create table `users` (
    `id` varchar(255) primary key,
    `name` text not null,
    `description` text not null,
    `email` varchar(255) not null,
    `password` varchar(255) not null,
    `account_id` varchar(255) not null,
    
    constraint `users_account_id_accounts_id` foreign key(`account_id`) REFERENCES `accounts`(id)
 );

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table accounts;

drop table payment_historys;

drop table payment_methods;

drop table bills;

drop table users;
