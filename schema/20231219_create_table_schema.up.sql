-- Create users table
Create TABLE users (
    user_id int PRIMARY KEY AUTO_INCREMENT NOT NULL,
    username varchar(255) not null unique,
    email varchar(255) not null unique,
    password varchar(512) not null,
    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create items table
Create TABLE items (
   item_id int PRIMARY KEY AUTO_INCREMENT NOT NULL,
   user_id int NOT NULL,
   title varchar(255) NOT NULL ,
   description text,
   status smallint NOT NULL DEFAULT 0,
   created_at datetime DEFAULT CURRENT_TIMESTAMP,
   updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   FOREIGN KEY (user_id) REFERENCES users (user_id)
);