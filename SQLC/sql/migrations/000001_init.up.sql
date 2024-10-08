CREATE TABLE categories (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(100) NOT NULL,
  description varchar(2000)
);

CREATE TABLE courses (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(100) NOT NULL,
  category_id varchar(36) NOT NULL REFERENCES categories(id),
  description varchar(2000),
  price decimal(10, 2) NOT NULL,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);
