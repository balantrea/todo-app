CREATE TABLE users
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE users_lists
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id int references users(id) on delete cascade not null,
    list_id int references todo_lists(id) on delete cascade not null
);

CREATE TABLE todo_items
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title varchar(255) not null,
    description varchar(255),
    done boolean not null default false
);

CREATE TABLE list_item
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    list_id int references todo_lists(id) on delete cascade not null,
    item_id int references todo_items(id) on delete cascade not null
);